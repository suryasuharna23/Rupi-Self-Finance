// Package conversation mengatur alur percakapan: onboarding, klarifikasi,
// koreksi, dan merangkai balasan (calon Conversation Service di Fase 2).
//
// Dependensi didefinisikan sebagai interface agar batasnya sama dengan
// kontrak gRPC di Fase 2.
package conversation

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/suryasuharna23/rupi/internal/parser"
	"github.com/suryasuharna23/rupi/internal/report"
	"github.com/suryasuharna23/rupi/internal/transaction"
	"github.com/suryasuharna23/rupi/internal/user"
)

type Incoming struct {
	Channel      string
	MessageID    string // unik per kanal, untuk idempotensi
	ChatID       int64
	ExternalID   int64 // ID pengguna di kanal (Telegram user ID)
	Text         string
	CallbackData string
	ReceivedAt   time.Time
}

type Button struct {
	Text string
	Data string
}

type Reply struct {
	Text    string
	Buttons [][]Button
}

type Users interface {
	GetOrCreate(ctx context.Context, telegramID int64) (user.User, bool, error)
	CompleteOnboarding(ctx context.Context, userID, displayName string) error
	Delete(ctx context.Context, userID string) error
}

type Transactions interface {
	Create(ctx context.Context, userID, sourceMessageID string, source parser.Source, txs []parser.ParsedTransaction) ([]transaction.Transaction, error)
	UpdateAmount(ctx context.Context, userID, txID string, amount int64) (transaction.Transaction, error)
	DeleteByMessage(ctx context.Context, userID, sourceMessageID string) (int64, error)
	DeleteLast(ctx context.Context, userID string) (transaction.Transaction, error)
	FindRecent(ctx context.Context, userID string, limit int) ([]parser.RecentTransaction, error)
	DeleteAllForUser(ctx context.Context, userID string) error
}

type Reports interface {
	Sum(ctx context.Context, userID string, q parser.ReportQuery) (int64, error)
	Summary(ctx context.Context, userID string, now time.Time) ([]report.Period, error)
}

type Idempotency interface {
	MarkProcessed(ctx context.Context, channel, messageID string) (bool, error)
	Unmark(ctx context.Context, channel, messageID string) error
}

type Parser interface {
	Parse(ctx context.Context, req parser.Request) parser.Result
}

type Service struct {
	users    Users
	txs      Transactions
	reports  Reports
	idem     Idempotency
	parser   Parser
	sessions *sessionStore
	limiter  *rateLimiter
}

func NewService(users Users, txs Transactions, reports Reports, idem Idempotency, p Parser) *Service {
	return &Service{
		users: users, txs: txs, reports: reports, idem: idem, parser: p,
		sessions: newSessionStore(30 * time.Minute),
		limiter:  newRateLimiter(30, time.Minute), // SEC-07
	}
}

const recentContextSize = 10

var undoRe = regexp.MustCompile(`(?i)^\s*(/undo|batal(in|kan)?\s+yang\s+terakhir|hapus\s+yang\s+terakhir)\s*$`)

const privacyNotice = "Rupi menyimpan transaksi yang kamu catat hanya untuk akunmu sendiri. " +
	"Teks pesan yang tidak dipahami aturan bisa dikirim ke layanan AI tanpa nama atau ID Telegram-mu. " +
	"Ketik /hapusakun kapan saja untuk menghapus semua datamu."

const helpText = "Cara pakai Rupi:\n" +
	"• Catat: \"beli bakso 15rb\", \"bakso 15rb sama es teh goceng\"\n" +
	"• Pemasukan: \"dapet transferan ortu 1,5jt\"\n" +
	"• Koreksi: \"yang bakso tadi 18rb\"\n" +
	"• Batalkan: /undo atau \"batalin yang terakhir\"\n" +
	"• Tanya: \"bulan ini habis berapa buat makan?\"\n" +
	"• Ringkasan: /ringkas\n" +
	"• Hapus akun: /hapusakun"

// HandleMessage memproses satu pesan masuk. Pesan yang sama yang dikirim
// ulang (retry webhook) tidak diproses dua kali (NFR-07).
func (s *Service) HandleMessage(ctx context.Context, in Incoming) (replies []Reply, err error) {
	first, err := s.idem.MarkProcessed(ctx, in.Channel, in.MessageID)
	if err != nil {
		return nil, err
	}
	if !first {
		return nil, nil
	}
	defer func() {
		if err != nil {
			if uerr := s.idem.Unmark(context.WithoutCancel(ctx), in.Channel, in.MessageID); uerr != nil {
				slog.Error("unmark processed message gagal", "err", uerr)
			}
		}
	}()

	if !s.limiter.Allow(in.ExternalID, in.ReceivedAt) {
		return text("Pelan-pelan ya, coba lagi sebentar lagi."), nil
	}

	u, created, err := s.users.GetOrCreate(ctx, in.ExternalID)
	if err != nil {
		return nil, err
	}
	loc, err := time.LoadLocation(u.Timezone)
	if err != nil {
		loc = time.FixedZone("WIB", 7*3600)
	}
	now := in.ReceivedAt.In(loc)

	if in.CallbackData != "" {
		return s.handleCallback(ctx, u, in.CallbackData)
	}

	msg := strings.TrimSpace(in.Text)
	sess := s.sessions.Get(u.ID, now)

	switch {
	case sess.awaitingName && !strings.HasPrefix(msg, "/"):
		return s.finishOnboarding(ctx, u, msg, now)
	case created || !u.Onboarded || msg == "/start":
		if u.Onboarded {
			return text(fmt.Sprintf("Halo lagi, %s! %s", u.DisplayName, helpText)), nil
		}
		s.sessions.Update(u.ID, now, func(st *session) { st.awaitingName = true })
		return []Reply{{Text: "Halo! Aku Rupi, asisten catat keuanganmu.\n\n" + privacyNotice}, {Text: "Mau dipanggil siapa?"}}, nil
	case msg == "/bantuan" || msg == "/help":
		return text(helpText), nil
	case msg == "/ringkas":
		return s.summary(ctx, u, now)
	case msg == "/hapusakun":
		return []Reply{{
			Text:    "Yakin hapus akun? Semua transaksimu akan dihapus permanen.",
			Buttons: [][]Button{{{Text: "Ya, hapus", Data: "hapusakun:ya"}, {Text: "Tidak", Data: "hapusakun:tidak"}}},
		}}, nil
	case undoRe.MatchString(msg):
		return s.undoLast(ctx, u)
	}

	// Jawaban untuk "Berapa nominalnya?" digabung dengan pesan sebelumnya.
	if sess.pendingText != "" && len(parser.FindAmounts(msg)) > 0 {
		msg = sess.pendingText + " " + msg
	}
	s.sessions.Update(u.ID, now, func(st *session) { st.pendingText = "" })

	recent, err := s.txs.FindRecent(ctx, u.ID, recentContextSize)
	if err != nil {
		return nil, err
	}
	res := s.parser.Parse(ctx, parser.Request{UserID: u.ID, Text: msg, Location: loc, ReceivedAt: now, Context: recent})

	switch res.Kind {
	case parser.KindTransactions:
		return s.saveTransactions(ctx, u, in.MessageID, res)
	case parser.KindQuery:
		return s.answerQuery(ctx, u, *res.Query)
	case parser.KindCorrection:
		return s.correct(ctx, u, *res.Correction)
	case parser.KindClarification:
		if res.Clarification == "Berapa nominalnya?" {
			s.sessions.Update(u.ID, now, func(st *session) { st.pendingText = msg })
		}
		return text(res.Clarification), nil
	default:
		return text(fmt.Sprintf("Halo %s! Ceritakan saja pengeluaranmu, misalnya \"beli bakso 15rb\".", u.DisplayName)), nil
	}
}

func (s *Service) finishOnboarding(ctx context.Context, u user.User, name string, now time.Time) ([]Reply, error) {
	name = strings.TrimSpace(name)
	if name == "" || len([]rune(name)) > 50 {
		return text("Nama panggilannya 1–50 karakter ya. Mau dipanggil siapa?"), nil
	}
	if err := s.users.CompleteOnboarding(ctx, u.ID, name); err != nil {
		return nil, err
	}
	s.sessions.Update(u.ID, now, func(st *session) { st.awaitingName = false })
	return text(fmt.Sprintf("Siap, %s! Zona waktumu WIB.\n\n%s", name, helpText)), nil
}

func (s *Service) saveTransactions(ctx context.Context, u user.User, messageID string, res parser.Result) ([]Reply, error) {
	for _, tx := range res.Transactions {
		if tx.Confidence < parser.ConfidenceThreshold { // FR-12
			return text(fmt.Sprintf("Maksudnya %s untuk %q? Tulis ulang dengan satuan ya, misalnya \"%s 15rb\".",
				Rupiah(tx.Amount), orDash(tx.Description), orDash(tx.Description))), nil
		}
	}
	saved, err := s.txs.Create(ctx, u.ID, messageID, res.Source, res.Transactions)
	if err != nil {
		return nil, err
	}
	var b strings.Builder
	b.WriteString("Tercatat:\n")
	for _, t := range saved {
		label := "Pengeluaran"
		if t.Type == parser.TypeIncome {
			label = "Pemasukan"
		}
		fmt.Fprintf(&b, "• %s %s — %s (%s)\n", label, orDash(t.Description), Rupiah(t.Amount), t.Category)
	}
	return []Reply{{
		Text:    strings.TrimSpace(b.String()),
		Buttons: [][]Button{{{Text: "Ubah", Data: "ubah:" + messageID}, {Text: "Batalkan", Data: "batal:" + messageID}}},
	}}, nil
}

func (s *Service) answerQuery(ctx context.Context, u user.User, q parser.ReportQuery) ([]Reply, error) {
	total, err := s.reports.Sum(ctx, u.ID, q)
	if err != nil {
		return nil, err
	}
	what := "Pengeluaran"
	if q.Type == parser.TypeIncome {
		what = "Pemasukan"
	}
	if len(q.Categories) > 0 {
		what += " " + strings.Join(q.Categories, ", ")
	}
	return text(fmt.Sprintf("%s %s – %s: %s",
		what, q.Range.Start.Format("2 Jan"), q.Range.End.AddDate(0, 0, -1).Format("2 Jan 2006"), Rupiah(total))), nil
}

func (s *Service) correct(ctx context.Context, u user.User, c parser.Correction) ([]Reply, error) {
	t, err := s.txs.UpdateAmount(ctx, u.ID, c.TransactionID, c.NewAmount)
	if errors.Is(err, transaction.ErrNotFound) {
		return text("Transaksinya tidak ketemu, mungkin sudah dihapus."), nil
	}
	if err != nil {
		return nil, err
	}
	return text(fmt.Sprintf("Diubah: %s jadi %s.", orDash(t.Description), Rupiah(t.Amount))), nil
}

func (s *Service) undoLast(ctx context.Context, u user.User) ([]Reply, error) {
	t, err := s.txs.DeleteLast(ctx, u.ID)
	if errors.Is(err, transaction.ErrNotFound) {
		return text("Tidak ada transaksi dalam 24 jam terakhir yang bisa dibatalkan."), nil
	}
	if err != nil {
		return nil, err
	}
	return text(fmt.Sprintf("Dibatalkan: %s %s.", orDash(t.Description), Rupiah(t.Amount))), nil
}

func (s *Service) summary(ctx context.Context, u user.User, now time.Time) ([]Reply, error) {
	periods, err := s.reports.Summary(ctx, u.ID, now)
	if err != nil {
		return nil, err
	}
	var b strings.Builder
	b.WriteString("Ringkasan:\n")
	for _, p := range periods {
		fmt.Fprintf(&b, "• %s: keluar %s, masuk %s\n", p.Label, Rupiah(p.Expense), Rupiah(p.Income))
	}
	return text(strings.TrimSpace(b.String())), nil
}

func (s *Service) handleCallback(ctx context.Context, u user.User, data string) ([]Reply, error) {
	action, arg, _ := strings.Cut(data, ":")
	switch action {
	case "batal":
		n, err := s.txs.DeleteByMessage(ctx, u.ID, arg)
		if err != nil {
			return nil, err
		}
		if n == 0 {
			return text("Transaksinya sudah dibatalkan sebelumnya."), nil
		}
		return text(fmt.Sprintf("%d transaksi dibatalkan.", n)), nil
	case "ubah":
		return text("Ketik koreksinya, misalnya \"yang bakso tadi 18rb\"."), nil
	case "hapusakun":
		if arg != "ya" {
			return text("Oke, akunmu tidak dihapus."), nil
		}
		if err := s.txs.DeleteAllForUser(ctx, u.ID); err != nil {
			return nil, err
		}
		if err := s.users.Delete(ctx, u.ID); err != nil {
			return nil, err
		}
		s.sessions.Delete(u.ID)
		return text("Akun dan semua transaksimu sudah dihapus. Sampai jumpa!"), nil
	}
	return nil, nil
}

func text(s string) []Reply { return []Reply{{Text: s}} }

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
