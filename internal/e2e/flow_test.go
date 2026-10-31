//go:build e2e

// Package e2e menguji alur chat lengkap terhadap PostgreSQL sungguhan dengan
// server Telegram tiruan (RUPI-101). Jalankan:
//
//	RUPI_TEST_DATABASE_URL=postgres://rupi:rupi@localhost:5432/rupi?sslmode=disable \
//	  go test -tags e2e ./internal/e2e/ -v
package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/suryasuharna23/rupi/internal/conversation"
	"github.com/suryasuharna23/rupi/internal/parser"
	"github.com/suryasuharna23/rupi/internal/report"
	"github.com/suryasuharna23/rupi/internal/storage"
	"github.com/suryasuharna23/rupi/internal/telegram"
	"github.com/suryasuharna23/rupi/internal/transaction"
	"github.com/suryasuharna23/rupi/internal/user"
)

const webhookSecret = "rahasia-test"

// fakeTelegram merekam setiap sendMessage yang dikirim bot.
type fakeTelegram struct {
	mu   sync.Mutex
	sent []sentMessage
}

type sentMessage struct {
	ChatID      int64  `json:"chat_id"`
	Text        string `json:"text"`
	ReplyMarkup *struct {
		InlineKeyboard [][]telegram.InlineButton `json:"inline_keyboard"`
	} `json:"reply_markup"`
}

func (f *fakeTelegram) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/sendMessage") {
		var m sentMessage
		json.NewDecoder(r.Body).Decode(&m)
		f.mu.Lock()
		f.sent = append(f.sent, m)
		f.mu.Unlock()
	}
	w.Write([]byte(`{"ok":true,"result":true}`))
}

func (f *fakeTelegram) drain() []sentMessage {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := f.sent
	f.sent = nil
	return out
}

type harness struct {
	t       *testing.T
	db      *pgxpool.Pool
	tg      *fakeTelegram
	webhook *httptest.Server
	nextID  int64
}

func setup(t *testing.T) *harness {
	url := os.Getenv("RUPI_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("RUPI_TEST_DATABASE_URL tidak diisi")
	}
	// Dijalankan dua kali: migrasi harus aman diulang.
	for range 2 {
		if err := storage.Migrate(url); err != nil {
			t.Fatalf("migrate: %v", err)
		}
	}
	db, err := storage.Connect(context.Background(), url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(db.Close)

	tg := &fakeTelegram{}
	tgServer := httptest.NewServer(tg)
	t.Cleanup(tgServer.Close)

	svc := conversation.NewService(user.NewStore(db), transaction.NewStore(db), report.NewStore(db),
		conversation.NewIdempotencyStore(db), parser.New(nil))
	adapter := telegram.NewAdapter(telegram.NewClient(tgServer.URL, "test-token"), svc)
	webhook := httptest.NewServer(adapter.WebhookHandler(webhookSecret))
	t.Cleanup(webhook.Close)

	return &harness{t: t, db: db, tg: tg, webhook: webhook, nextID: time.Now().UnixNano() % 1e12}
}

// newUser mengembalikan ID Telegram unik agar test bisa diulang tanpa reset database.
func (h *harness) newUser() int64 {
	h.nextID++
	return h.nextID
}

func (h *harness) post(update telegram.Update, secret string) int {
	body, _ := json.Marshal(update)
	req, _ := http.NewRequest(http.MethodPost, h.webhook.URL, bytes.NewReader(body))
	req.Header.Set("X-Telegram-Bot-Api-Secret-Token", secret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		h.t.Fatal(err)
	}
	resp.Body.Close()
	return resp.StatusCode
}

func (h *harness) messageUpdate(userID int64, text string) telegram.Update {
	h.nextID++
	return telegram.Update{UpdateID: h.nextID, Message: &telegram.Message{
		MessageID: h.nextID, From: &telegram.User{ID: userID}, Chat: telegram.Chat{ID: userID},
		Date: time.Now().Unix(), Text: text,
	}}
}

// say mengirim pesan dan mengembalikan balasan bot.
func (h *harness) say(userID int64, text string) []sentMessage {
	h.t.Helper()
	if code := h.post(h.messageUpdate(userID, text), webhookSecret); code != http.StatusOK {
		h.t.Fatalf("webhook %q: status %d", text, code)
	}
	return h.tg.drain()
}

func (h *harness) click(userID int64, data string) []sentMessage {
	h.t.Helper()
	h.nextID++
	u := telegram.Update{UpdateID: h.nextID, CallbackQuery: &telegram.CallbackQuery{
		ID: "cb" + time.Now().Format("150405.000000000"), From: telegram.User{ID: userID},
		Message: &telegram.Message{Chat: telegram.Chat{ID: userID}}, Data: data,
	}}
	if code := h.post(u, webhookSecret); code != http.StatusOK {
		h.t.Fatalf("callback %q: status %d", data, code)
	}
	return h.tg.drain()
}

func (h *harness) count(query string, args ...any) int {
	h.t.Helper()
	var n int
	if err := h.db.QueryRow(context.Background(), query, args...).Scan(&n); err != nil {
		h.t.Fatalf("%s: %v", query, err)
	}
	return n
}

func (h *harness) activeTx(telegramID int64) int {
	return h.count(`
		SELECT count(*) FROM transaction_svc.transactions t
		JOIN user_svc.users u ON u.id = t.user_id
		WHERE u.telegram_id = $1 AND t.deleted_at IS NULL`, telegramID)
}

func expectText(t *testing.T, got []sentMessage, want ...string) {
	t.Helper()
	all := make([]string, len(got))
	for i, m := range got {
		all[i] = m.Text
	}
	joined := strings.Join(all, "\n---\n")
	for _, w := range want {
		if !strings.Contains(joined, w) {
			t.Errorf("balasan tidak memuat %q:\n%s", w, joined)
		}
	}
}

func TestChatFlow(t *testing.T) {
	h := setup(t)
	dina := h.newUser()

	t.Run("onboarding", func(t *testing.T) {
		got := h.say(dina, "/start")
		expectText(t, got, "Aku Rupi", "/hapusakun", "Mau dipanggil siapa?")
		expectText(t, h.say(dina, "Dina"), "Siap, Dina!")
		if n := h.count(`SELECT count(*) FROM user_svc.users WHERE telegram_id = $1 AND onboarded`, dina); n != 1 {
			t.Fatalf("pengguna onboarded = %d, want 1", n)
		}
	})

	t.Run("catat satu transaksi dengan tombol", func(t *testing.T) {
		got := h.say(dina, "beli bakso 15rb")
		expectText(t, got, "Tercatat", "Rp15.000", "Makanan & Minuman")
		if len(got) != 1 || got[0].ReplyMarkup == nil || len(got[0].ReplyMarkup.InlineKeyboard[0]) != 2 {
			t.Fatalf("balasan harus punya tombol Ubah / Batalkan: %+v", got)
		}
	})

	t.Run("dua transaksi dalam satu pesan", func(t *testing.T) {
		before := h.activeTx(dina)
		expectText(t, h.say(dina, "bakso 15rb sama es teh goceng"), "Rp15.000", "Rp5.000")
		if got := h.activeTx(dina) - before; got != 2 {
			t.Fatalf("transaksi baru = %d, want 2", got)
		}
	})

	t.Run("webhook ganda tidak menggandakan transaksi", func(t *testing.T) {
		before := h.activeTx(dina)
		u := h.messageUpdate(dina, "parkir 2rb")
		h.post(u, webhookSecret)
		h.post(u, webhookSecret)
		if got := len(h.tg.drain()); got != 1 {
			t.Errorf("balasan = %d, want 1", got)
		}
		if got := h.activeTx(dina) - before; got != 1 {
			t.Fatalf("transaksi baru = %d, want 1", got)
		}
	})

	t.Run("tanya balik nominal", func(t *testing.T) {
		expectText(t, h.say(dina, "makan siang"), "Berapa nominalnya?")
		expectText(t, h.say(dina, "15rb"), "Tercatat", "makan siang", "Rp15.000")
	})

	t.Run("koreksi nominal dengan jejak audit", func(t *testing.T) {
		expectText(t, h.say(dina, "yang bakso tadi 18rb"), "Diubah", "Rp18.000")
		n := h.count(`
			SELECT count(*) FROM transaction_svc.transaction_edits e
			JOIN transaction_svc.transactions t ON t.id = e.transaction_id
			JOIN user_svc.users u ON u.id = t.user_id
			WHERE u.telegram_id = $1 AND e.field = 'amount' AND e.old_value = '15000' AND e.new_value = '18000'`, dina)
		if n != 1 {
			t.Fatalf("baris transaction_edits = %d, want 1", n)
		}
	})

	// Makanan: 15.000 + 18.000 (dikoreksi) + 5.000 + 15.000 = 53.000; parkir 2.000.
	t.Run("tanya total", func(t *testing.T) {
		expectText(t, h.say(dina, "bulan ini habis berapa buat makan?"), "Makanan & Minuman", "Rp53.000")
		expectText(t, h.say(dina, "bulan ini habis berapa?"), "Rp55.000")
	})

	t.Run("ringkasan", func(t *testing.T) {
		expectText(t, h.say(dina, "/ringkas"), "Hari ini: keluar Rp55.000", "Bulan ini: keluar Rp55.000")
	})

	t.Run("tombol batalkan", func(t *testing.T) {
		got := h.say(dina, "fotokopi 5rb")
		before := h.activeTx(dina)
		data := got[0].ReplyMarkup.InlineKeyboard[0][1].CallbackData
		expectText(t, h.click(dina, data), "1 transaksi dibatalkan")
		if h.activeTx(dina) != before-1 {
			t.Fatal("transaksi tidak terhapus")
		}
		expectText(t, h.click(dina, data), "sudah dibatalkan")
	})

	t.Run("undo", func(t *testing.T) {
		before := h.activeTx(dina)
		expectText(t, h.say(dina, "batalin yang terakhir"), "Dibatalkan", "makan siang")
		expectText(t, h.say(dina, "/undo"), "Dibatalkan")
		if h.activeTx(dina) != before-2 {
			t.Fatal("undo tidak menghapus transaksi")
		}
	})

	t.Run("pengguna lain tidak melihat data Dina", func(t *testing.T) {
		raka := h.newUser()
		h.say(raka, "/start")
		h.say(raka, "Raka")
		expectText(t, h.say(raka, "bulan ini habis berapa?"), "Rp0")
		expectText(t, h.say(raka, "yang bakso tadi 20rb"), "Belum ketemu")
	})

	t.Run("secret token salah ditolak", func(t *testing.T) {
		if code := h.post(h.messageUpdate(dina, "beli kopi 10rb"), "salah"); code != http.StatusForbidden {
			t.Fatalf("status = %d, want 403", code)
		}
	})

	t.Run("hapus akun", func(t *testing.T) {
		expectText(t, h.say(dina, "/hapusakun"), "Yakin hapus akun?")
		expectText(t, h.click(dina, "hapusakun:tidak"), "tidak dihapus")
		expectText(t, h.click(dina, "hapusakun:ya"), "sudah dihapus")
		if n := h.count(`SELECT count(*) FROM user_svc.users WHERE telegram_id = $1 AND deleted_at IS NULL`, dina); n != 0 {
			t.Fatalf("pengguna aktif = %d, want 0", n)
		}
		// Daftar ulang setelah dihapus membuat akun baru yang kosong.
		expectText(t, h.say(dina, "/start"), "Mau dipanggil siapa?")
		if h.activeTx(dina) != 0 {
			t.Fatal("akun baru tidak boleh mewarisi transaksi lama")
		}
	})
}
