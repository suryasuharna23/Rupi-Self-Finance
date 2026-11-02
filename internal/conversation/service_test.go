package conversation

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/suryasuharna23/rupi/internal/parser"
	"github.com/suryasuharna23/rupi/internal/report"
)

const dinaTG = 42

var errDB = errors.New("database down")

type harness struct {
	t       *testing.T
	svc     *Service
	users   *fakeUsers
	txs     *fakeTxs
	reports *fakeReports
	idem    *fakeIdempotency
	parser  *fakeParser
	now     time.Time
	nextMsg int
}

// newHarness returns a Service with Dina already onboarded.
func newHarness(t *testing.T) *harness {
	h := &harness{
		t: t, users: newFakeUsers(), txs: &fakeTxs{}, reports: &fakeReports{},
		idem: newFakeIdempotency(), parser: &fakeParser{result: parser.Result{Kind: parser.KindSmallTalk}},
		now: time.Date(2026, time.October, 1, 12, 0, 0, 0, time.UTC),
	}
	h.users.add(dinaTG, "Dina", true)
	h.svc = NewService(h.users, h.txs, h.reports, h.idem, h.parser)
	return h
}

func (h *harness) dinaID() string { return h.users.byTG[dinaTG].ID }

// incoming builds a new message from Dina with a fresh message ID.
func (h *harness) incoming(text string) Incoming {
	h.nextMsg++
	h.now = h.now.Add(time.Second)
	return Incoming{Channel: "telegram", MessageID: fmt.Sprint(h.nextMsg), ChatID: dinaTG, ExternalID: dinaTG, Text: text, ReceivedAt: h.now}
}

func (h *harness) handle(in Incoming) []Reply {
	h.t.Helper()
	replies, err := h.svc.HandleMessage(context.Background(), in)
	if err != nil {
		h.t.Fatalf("HandleMessage(%q): %v", in.Text+in.CallbackData, err)
	}
	return replies
}

func (h *harness) send(text string) []Reply { h.t.Helper(); return h.handle(h.incoming(text)) }

func (h *harness) tap(data string) []Reply {
	h.t.Helper()
	in := h.incoming("")
	in.CallbackData = data
	return h.handle(in)
}

// replyText joins every reply's text for easy matching.
func replyText(replies []Reply) string {
	parts := make([]string, len(replies))
	for i, r := range replies {
		parts[i] = r.Text
	}
	return strings.Join(parts, "\n")
}

func assertContains(t *testing.T, replies []Reply, want string) {
	t.Helper()
	if got := replyText(replies); !strings.Contains(got, want) {
		t.Errorf("reply %q does not contain %q", got, want)
	}
}

func expenseResult(amount int64, desc, category string, confidence float64) parser.Result {
	return parser.Result{Kind: parser.KindTransactions, Source: parser.SourceRule, Transactions: []parser.ParsedTransaction{
		{Type: parser.TypeExpense, Amount: amount, Description: desc, Category: category, Confidence: confidence},
	}}
}

func TestDuplicateMessageIsProcessedOnce(t *testing.T) {
	h := newHarness(t)
	h.parser.result = expenseResult(15000, "bakso", "Makanan & Minuman", 0.95)
	in := h.incoming("beli bakso 15rb")

	if replies := h.handle(in); len(replies) != 1 {
		t.Fatalf("first delivery: %d replies, want 1", len(replies))
	}
	if replies := h.handle(in); replies != nil {
		t.Errorf("redelivery replied %q, want nothing", replyText(replies))
	}
	if n := len(h.txs.live(h.dinaID())); n != 1 {
		t.Errorf("%d transactions stored, want 1", n)
	}
}

func TestStoreErrorUnmarksMessage(t *testing.T) {
	h := newHarness(t)
	h.parser.result = expenseResult(15000, "bakso", "Makanan & Minuman", 0.95)
	h.txs.createErr = errDB
	in := h.incoming("beli bakso 15rb")

	if _, err := h.svc.HandleMessage(context.Background(), in); !errors.Is(err, errDB) {
		t.Fatalf("err = %v, want %v", err, errDB)
	}
	if len(h.idem.unmarked) != 1 || h.idem.unmarked[0] != in.MessageID {
		t.Fatalf("unmarked = %v, want [%s]", h.idem.unmarked, in.MessageID)
	}

	h.txs.createErr = nil
	assertContains(t, h.handle(in), "Tercatat")
	if n := len(h.txs.live(h.dinaID())); n != 1 {
		t.Errorf("%d transactions after retry, want 1", n)
	}
}

func TestErrors(t *testing.T) {
	cases := []struct {
		name   string
		break_ func(h *harness)
		text   string
	}{
		{"mark processed", func(h *harness) { h.idem.markErr = errDB }, "halo"},
		{"get user", func(h *harness) { h.users.getErr = errDB }, "halo"},
		{"find recent", func(h *harness) { h.txs.recentErr = errDB }, "halo"},
		{"summary", func(h *harness) { h.reports.err = errDB }, "/ringkas"},
		{"undo", func(h *harness) { h.txs.deleteLastErr = errDB }, "/undo"},
		{"sum", func(h *harness) {
			h.reports.err = errDB
			h.parser.result = parser.Result{Kind: parser.KindQuery, Query: &parser.ReportQuery{}}
		}, "habis berapa?"},
		{"correction", func(h *harness) {
			h.txs.updateErr = errDB
			h.parser.result = parser.Result{Kind: parser.KindCorrection, Correction: &parser.Correction{TransactionID: "tx-1", NewAmount: 1}}
		}, "yang bakso tadi 18rb"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := newHarness(t)
			c.break_(h)
			if _, err := h.svc.HandleMessage(context.Background(), h.incoming(c.text)); !errors.Is(err, errDB) {
				t.Errorf("err = %v, want %v", err, errDB)
			}
		})
	}

	t.Run("unmark failure keeps the original error", func(t *testing.T) {
		h := newHarness(t)
		h.txs.recentErr = errDB
		h.idem.unmarkErr = errors.New("unmark failed")
		if _, err := h.svc.HandleMessage(context.Background(), h.incoming("halo")); !errors.Is(err, errDB) {
			t.Errorf("err = %v, want %v", err, errDB)
		}
	})
}

func TestLowConfidenceAsksInsteadOfSaving(t *testing.T) {
	h := newHarness(t)
	h.parser.result = expenseResult(15, "bakso", "Makanan & Minuman", 0.5) // FR-12

	replies := h.send("bakso 15")
	assertContains(t, replies, `Maksudnya Rp15 untuk "bakso"?`)
	if n := len(h.txs.live(h.dinaID())); n != 0 {
		t.Errorf("%d transactions saved, want 0", n)
	}
}

func TestAmountReplyIsMergedWithPendingMessage(t *testing.T) {
	h := newHarness(t)
	h.parser.result = parser.Result{Kind: parser.KindClarification, Clarification: "Berapa nominalnya?"}
	assertContains(t, h.send("beli bakso"), "Berapa nominalnya?")

	h.parser.result = expenseResult(15000, "bakso", "Makanan & Minuman", 0.95)
	h.send("15rb")
	if got := h.parser.lastText(); got != "beli bakso 15rb" {
		t.Errorf("parser got %q, want the merged %q", got, "beli bakso 15rb")
	}

	// The pending message is used once only.
	h.send("20rb")
	if got := h.parser.lastText(); got != "20rb" {
		t.Errorf("parser got %q after the merge, want %q", got, "20rb")
	}
}

func TestPendingMessageIsDroppedWithoutAmount(t *testing.T) {
	h := newHarness(t)
	h.parser.result = parser.Result{Kind: parser.KindClarification, Clarification: "Berapa nominalnya?"}
	h.send("beli bakso")

	h.parser.result = parser.Result{Kind: parser.KindSmallTalk}
	h.send("eh nggak jadi")
	h.send("15rb")
	if got := h.parser.lastText(); got != "15rb" {
		t.Errorf("parser got %q, want %q (pending message should be dropped)", got, "15rb")
	}
}

func TestOtherClarificationIsNotRemembered(t *testing.T) {
	h := newHarness(t)
	h.parser.result = parser.Result{Kind: parser.KindClarification, Clarification: "Maksudnya pengeluaran atau pemasukan?"}
	assertContains(t, h.send("bakso"), "pengeluaran atau pemasukan")
	h.send("15rb")
	if got := h.parser.lastText(); got != "15rb" {
		t.Errorf("parser got %q, want %q", got, "15rb")
	}
}

func TestRateLimit(t *testing.T) {
	h := newHarness(t)
	start := h.now
	for i := 1; i <= 30; i++ {
		in := h.incoming("halo")
		in.ReceivedAt = start
		if got := replyText(h.handle(in)); strings.Contains(got, "Pelan-pelan") {
			t.Fatalf("message %d was rate limited", i)
		}
	}
	in := h.incoming("halo")
	in.ReceivedAt = start.Add(59 * time.Second)
	assertContains(t, h.handle(in), "Pelan-pelan") // SEC-07
	if len(h.parser.requests) != 30 {
		t.Errorf("parser called %d times, want 30", len(h.parser.requests))
	}

	in = h.incoming("halo")
	in.ReceivedAt = start.Add(time.Minute)
	if got := replyText(h.handle(in)); strings.Contains(got, "Pelan-pelan") {
		t.Error("a new minute must be allowed again")
	}
}

func TestOnboarding(t *testing.T) {
	h := newHarness(t)
	const newTG = 7
	msg := func(text string) Incoming {
		in := h.incoming(text)
		in.ExternalID, in.ChatID = newTG, newTG
		return in
	}

	replies := h.handle(msg("/start"))
	if len(replies) != 2 {
		t.Fatalf("greeting: %d replies, want 2", len(replies))
	}
	assertContains(t, replies, "Aku Rupi")
	assertContains(t, replies, "/hapusakun") // privacy notice
	assertContains(t, replies, "Mau dipanggil siapa?")

	for _, bad := range []string{"   ", strings.Repeat("a", 51)} {
		assertContains(t, h.handle(msg(bad)), "1–50 karakter")
	}
	if h.users.byTG[newTG].Onboarded {
		t.Fatal("an invalid name must not finish onboarding")
	}

	assertContains(t, h.handle(msg("  Budi ")), "Siap, Budi!")
	if u := h.users.byTG[newTG]; !u.Onboarded || u.DisplayName != "Budi" {
		t.Errorf("user after onboarding = %+v", u)
	}

	h.parser.result = parser.Result{Kind: parser.KindSmallTalk}
	assertContains(t, h.handle(msg("halo")), "Halo Budi!")
}

func TestOnboardingCommandWhileAwaitingNameRestarts(t *testing.T) {
	h := newHarness(t)
	in := h.incoming("halo")
	in.ExternalID = 7
	h.handle(in)

	in = h.incoming("/start")
	in.ExternalID = 7
	assertContains(t, h.handle(in), "Mau dipanggil siapa?")
	if h.users.byTG[7].Onboarded {
		t.Error("a command must not be taken as the name")
	}
}

func TestOnboardingError(t *testing.T) {
	h := newHarness(t)
	h.users.onboardErr = errDB
	in := h.incoming("/start")
	in.ExternalID = 7
	h.handle(in)

	in = h.incoming("Budi")
	in.ExternalID = 7
	if _, err := h.svc.HandleMessage(context.Background(), in); !errors.Is(err, errDB) {
		t.Errorf("err = %v, want %v", err, errDB)
	}
}

func TestCommands(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{"/start", "Halo lagi, Dina!"},
		{"/bantuan", "Cara pakai Rupi"},
		{"/help", "Cara pakai Rupi"},
		{"/undo", "Tidak ada transaksi dalam 24 jam terakhir"},
		{"batalin yang terakhir", "Tidak ada transaksi dalam 24 jam terakhir"},
	}
	for _, c := range cases {
		t.Run(c.text, func(t *testing.T) {
			h := newHarness(t)
			assertContains(t, h.send(c.text), c.want)
			if len(h.parser.requests) != 0 {
				t.Errorf("commands must not reach the parser")
			}
		})
	}
}

func TestSummaryCommand(t *testing.T) {
	h := newHarness(t)
	h.reports.periods = []report.Period{{Label: "Hari ini", Expense: 15000}, {Label: "Minggu ini", Expense: 65000, Income: 500000}}
	got := replyText(h.send("/ringkas"))
	want := "Ringkasan:\n• Hari ini: keluar Rp15.000, masuk Rp0\n• Minggu ini: keluar Rp65.000, masuk Rp500.000"
	if got != want {
		t.Errorf("summary = %q, want %q", got, want)
	}
	if h.reports.summaryNow.Location().String() != "Asia/Jakarta" {
		t.Errorf("summary computed in %s, want the user's time zone", h.reports.summaryNow.Location())
	}
}

func TestUnknownTimezoneFallsBackToWIB(t *testing.T) {
	h := newHarness(t)
	h.users.byTG[dinaTG].Timezone = "Mars/Olympus"
	h.send("/ringkas")
	if _, offset := h.reports.summaryNow.Zone(); offset != 7*3600 {
		t.Errorf("offset = %d, want WIB (+7)", offset)
	}
}

func TestUndo(t *testing.T) {
	h := newHarness(t)
	h.parser.result = expenseResult(15000, "bakso", "Makanan & Minuman", 0.95)
	h.send("beli bakso 15rb")

	assertContains(t, h.send("/undo"), "Dibatalkan: bakso Rp15.000.")
	if n := len(h.txs.live(h.dinaID())); n != 0 {
		t.Errorf("%d transactions left, want 0", n)
	}
}

func TestParserResults(t *testing.T) {
	income := parser.ParsedTransaction{Type: parser.TypeIncome, Amount: 1500000, Description: "", Category: "Uang Saku", Confidence: 0.9}
	start := time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		name   string
		result parser.Result
		want   string
	}{
		{"expense and income", parser.Result{Kind: parser.KindTransactions, Transactions: []parser.ParsedTransaction{
			{Type: parser.TypeExpense, Amount: 15000, Description: "bakso", Category: "Makanan & Minuman", Confidence: 0.9}, income,
		}}, "Tercatat:\n• Pengeluaran bakso — Rp15.000 (Makanan & Minuman)\n• Pemasukan - — Rp1.500.000 (Uang Saku)"},
		{"expense query", parser.Result{Kind: parser.KindQuery, Query: &parser.ReportQuery{
			Type: parser.TypeExpense, Categories: []string{"Makanan & Minuman"}, Range: parser.DateRange{Start: start, End: start.AddDate(0, 1, 0)},
		}}, "Pengeluaran Makanan & Minuman 1 Oct – 31 Oct 2026: Rp43.000"},
		{"income query", parser.Result{Kind: parser.KindQuery, Query: &parser.ReportQuery{
			Type: parser.TypeIncome, Range: parser.DateRange{Start: start, End: start.AddDate(0, 0, 1)},
		}}, "Pemasukan 1 Oct – 1 Oct 2026: Rp43.000"},
		{"correction of a missing transaction", parser.Result{Kind: parser.KindCorrection, Correction: &parser.Correction{TransactionID: "nope", NewAmount: 1}},
			"Transaksinya tidak ketemu"},
		{"small talk", parser.Result{Kind: parser.KindSmallTalk}, `Halo Dina! Ceritakan saja pengeluaranmu`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := newHarness(t)
			h.reports.sum = 43000
			h.parser.result = c.result
			replies := h.send("pesan")
			if got := replyText(replies); got != c.want && !strings.HasPrefix(got, c.want) {
				t.Errorf("reply = %q, want %q", got, c.want)
			}
		})
	}
}

func TestSavedTransactionsHaveButtons(t *testing.T) {
	h := newHarness(t)
	h.parser.result = expenseResult(15000, "bakso", "Makanan & Minuman", 0.95)
	in := h.incoming("beli bakso 15rb")
	replies := h.handle(in)
	want := [][]Button{{{Text: "Ubah", Data: "ubah:" + in.MessageID}, {Text: "Batalkan", Data: "batal:" + in.MessageID}}}
	if fmt.Sprint(replies[0].Buttons) != fmt.Sprint(want) {
		t.Errorf("buttons = %v, want %v", replies[0].Buttons, want)
	}
}

func TestParserGetsRecentTransactionsAsContext(t *testing.T) {
	h := newHarness(t)
	h.parser.result = expenseResult(15000, "bakso", "Makanan & Minuman", 0.95)
	h.send("beli bakso 15rb")

	h.parser.result = parser.Result{Kind: parser.KindCorrection, Correction: &parser.Correction{TransactionID: "tx-1", NewAmount: 18000}}
	assertContains(t, h.send("yang bakso tadi 18rb"), "Diubah: bakso jadi Rp18.000.")

	req := h.parser.requests[len(h.parser.requests)-1]
	if len(req.Context) != 1 || req.Context[0].ID != "tx-1" {
		t.Errorf("parser context = %+v, want the bakso transaction", req.Context)
	}
	if req.UserID != h.dinaID() || req.Location.String() != "Asia/Jakarta" {
		t.Errorf("request user/location = %s/%s", req.UserID, req.Location)
	}
}

func TestCallbacks(t *testing.T) {
	t.Run("cancel", func(t *testing.T) {
		h := newHarness(t)
		h.parser.result = parser.Result{Kind: parser.KindTransactions, Transactions: []parser.ParsedTransaction{
			{Type: parser.TypeExpense, Amount: 15000, Category: "Makanan & Minuman", Confidence: 0.9},
			{Type: parser.TypeExpense, Amount: 5000, Category: "Makanan & Minuman", Confidence: 0.9},
		}}
		in := h.incoming("bakso 15rb sama es teh goceng")
		h.handle(in)

		assertContains(t, h.tap("batal:"+in.MessageID), "2 transaksi dibatalkan.")
		assertContains(t, h.tap("batal:"+in.MessageID), "sudah dibatalkan sebelumnya")
	})

	t.Run("edit", func(t *testing.T) {
		assertContains(t, newHarness(t).tap("ubah:1"), "Ketik koreksinya")
	})

	t.Run("keep account", func(t *testing.T) {
		h := newHarness(t)
		assertContains(t, h.send("/hapusakun"), "Yakin hapus akun?")
		assertContains(t, h.tap("hapusakun:tidak"), "tidak dihapus")
		if len(h.users.deleted) != 0 {
			t.Error("account was deleted")
		}
	})

	t.Run("delete account", func(t *testing.T) {
		h := newHarness(t)
		id := h.dinaID()
		replies := h.send("/hapusakun")
		if len(replies[0].Buttons) != 1 || len(replies[0].Buttons[0]) != 2 {
			t.Fatalf("confirmation buttons = %v", replies[0].Buttons)
		}
		assertContains(t, h.tap(replies[0].Buttons[0][0].Data), "sudah dihapus")
		if len(h.txs.deletedAllFor) != 1 || len(h.users.deleted) != 1 || h.users.deleted[0] != id {
			t.Errorf("deleted transactions for %v, users %v; want %s", h.txs.deletedAllFor, h.users.deleted, id)
		}
		// Signing up again starts fresh onboarding.
		assertContains(t, h.send("halo"), "Mau dipanggil siapa?")
	})

	t.Run("unknown action", func(t *testing.T) {
		if replies := newHarness(t).tap("entah:1"); replies != nil {
			t.Errorf("replies = %v, want none", replies)
		}
	})

	errCases := []struct {
		name   string
		break_ func(h *harness)
		data   string
	}{
		{"cancel fails", func(h *harness) { h.txs.deleteMsgErr = errDB }, "batal:1"},
		{"deleting transactions fails", func(h *harness) { h.txs.deleteAllErr = errDB }, "hapusakun:ya"},
		{"deleting user fails", func(h *harness) { h.users.deleteErr = errDB }, "hapusakun:ya"},
	}
	for _, c := range errCases {
		t.Run(c.name, func(t *testing.T) {
			h := newHarness(t)
			c.break_(h)
			in := h.incoming("")
			in.CallbackData = c.data
			if _, err := h.svc.HandleMessage(context.Background(), in); !errors.Is(err, errDB) {
				t.Errorf("err = %v, want %v", err, errDB)
			}
		})
	}
}

// Compile-time checks that the fakes satisfy the Service interfaces.
var (
	_ Users        = (*fakeUsers)(nil)
	_ Transactions = (*fakeTxs)(nil)
	_ Reports      = (*fakeReports)(nil)
	_ Idempotency  = (*fakeIdempotency)(nil)
	_ Parser       = (*fakeParser)(nil)
)
