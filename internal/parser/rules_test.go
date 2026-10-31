package parser

import (
	"testing"
	"time"
)

var jakarta = time.FixedZone("WIB", 7*3600)

// Jumat, 30 Oktober 2026 pukul 20.00 WIB.
var testNow = time.Date(2026, 10, 30, 20, 0, 0, 0, jakarta)

func parse(text string, recent ...RecentTransaction) Result {
	res, _ := ParseRules(Request{Text: text, Location: jakarta, ReceivedAt: testNow, Context: recent})
	return res
}

func TestParseRulesTransactions(t *testing.T) {
	type want struct {
		typ      TransactionType
		amount   int64
		category string
		date     time.Time
	}
	today, yesterday := testNow, time.Date(2026, 10, 29, 12, 0, 0, 0, jakarta)
	cases := []struct {
		in   string
		want []want
	}{
		{"beli bakso 15rb", []want{{TypeExpense, 15_000, CatFood, today}}},
		{"bakso 15rb sama es teh goceng", []want{{TypeExpense, 15_000, CatFood, today}, {TypeExpense, 5_000, CatFood, today}}},
		{"kemarin isi bensin 30rb", []want{{TypeExpense, 30_000, CatTransport, yesterday}}},
		{"dapet transferan ortu 1,5jt", []want{{TypeIncome, 1_500_000, CatAllowance, today}}},
		{"print tugas 3500", []want{{TypeExpense, 3_500, CatKuliah, today}}},
		{"makan sama temen 20rb", []want{{TypeExpense, 20_000, CatFood, today}}},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			res := parse(c.in)
			if res.Kind != KindTransactions || len(res.Transactions) != len(c.want) {
				t.Fatalf("got %+v", res)
			}
			for i, w := range c.want {
				tx := res.Transactions[i]
				if tx.Type != w.typ || tx.Amount != w.amount || tx.Category != w.category || !tx.OccurredAt.Equal(w.date) {
					t.Errorf("tx[%d] = %+v, want %+v", i, tx, w)
				}
				if tx.Confidence < ConfidenceThreshold {
					t.Errorf("tx[%d] confidence %.2f di bawah ambang", i, tx.Confidence)
				}
			}
		})
	}
}

func TestParseRulesOtherKinds(t *testing.T) {
	if res := parse("makan siang"); res.Kind != KindClarification || res.Clarification != "Berapa nominalnya?" {
		t.Errorf("makan siang: got %+v", res)
	}
	if res := parse("halo rupi"); res.Kind != KindSmallTalk {
		t.Errorf("halo rupi: got %+v", res)
	}
	if _, ok := ParseRules(Request{Text: "apa kabar dunia", ReceivedAt: testNow}); ok {
		t.Errorf("kalimat tak dikenal harus diteruskan ke LLM")
	}
}

func TestParseRulesQuery(t *testing.T) {
	res := parse("bulan ini habis berapa buat makan?")
	if res.Kind != KindQuery {
		t.Fatalf("got %+v", res)
	}
	q := res.Query
	wantStart := time.Date(2026, 10, 1, 0, 0, 0, 0, jakarta)
	if q.Type != TypeExpense || len(q.Categories) != 1 || q.Categories[0] != CatFood || !q.Range.Start.Equal(wantStart) {
		t.Errorf("query = %+v", q)
	}
}

func TestParseRulesCorrection(t *testing.T) {
	recent := []RecentTransaction{
		{ID: "tx-2", Description: "es teh", Amount: 5_000},
		{ID: "tx-1", Description: "bakso", Amount: 15_000},
	}
	res := parse("yang bakso tadi 18rb", recent...)
	if res.Kind != KindCorrection || res.Correction.TransactionID != "tx-1" || res.Correction.NewAmount != 18_000 {
		t.Errorf("got %+v", res)
	}
}

func TestResolveTime(t *testing.T) {
	cases := map[string]time.Time{
		"kemarin lusa":  time.Date(2026, 10, 28, 12, 0, 0, 0, jakarta),
		"senin kemarin": time.Date(2026, 10, 26, 12, 0, 0, 0, jakarta),
		"jumat lalu":    time.Date(2026, 10, 23, 12, 0, 0, 0, jakarta),
		"minggu lalu":   time.Date(2026, 10, 23, 12, 0, 0, 0, jakarta),
		"tgl 3":         time.Date(2026, 10, 3, 12, 0, 0, 0, jakarta),
		"tgl 31":        time.Date(2026, 9, 30, 12, 0, 0, 0, jakarta), // September hanya 30 hari
		"tadi pagi":     time.Date(2026, 10, 30, 8, 0, 0, 0, jakarta),
		"tanpa waktu":   testNow,
	}
	for in, want := range cases {
		if got, _ := ResolveTime(in, testNow); !got.Equal(want) {
			t.Errorf("ResolveTime(%q) = %s, want %s", in, got, want)
		}
	}
}
