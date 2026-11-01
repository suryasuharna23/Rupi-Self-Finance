//go:build integration

package report_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/suryasuharna23/rupi/internal/parser"
	"github.com/suryasuharna23/rupi/internal/report"
	"github.com/suryasuharna23/rupi/internal/storage/storagetest"
	"github.com/suryasuharna23/rupi/internal/transaction"
)

func TestMain(m *testing.M) { os.Exit(storagetest.Run(m)) }

var jakarta = time.FixedZone("WIB", 7*3600)

func at(day, hour int) time.Time {
	return time.Date(2026, time.October, day, hour, 0, 0, 0, jakarta)
}

func seed(t *testing.T, userID string, txs ...parser.ParsedTransaction) {
	t.Helper()
	store := transaction.NewStore(storagetest.DB(t))
	if _, err := store.Create(context.Background(), userID, "seed", parser.SourceRule, txs); err != nil {
		t.Fatal(err)
	}
}

func tx(typ parser.TransactionType, amount int64, category string, when time.Time) parser.ParsedTransaction {
	return parser.ParsedTransaction{Type: typ, Amount: amount, Category: category, OccurredAt: when}
}

// TestSummaryWeekSpanningMonthStart uses Thursday 1 October 2026: the week
// started on Monday 28 September, before the month did.
func TestSummaryWeekSpanningMonthStart(t *testing.T) {
	ctx := context.Background()
	user := storagetest.NewUser(t)
	seed(t, user,
		tx(parser.TypeExpense, 1000, "Lainnya", at(-3, 12)),           // Sun 27 Sep: last week
		tx(parser.TypeExpense, 20000, "Transportasi", at(-2, 0)),      // Mon 28 Sep 00:00: first instant of the week
		tx(parser.TypeIncome, 500000, "Uang Saku", at(0, 9)),          // Wed 30 Sep: this week, not this month
		tx(parser.TypeExpense, 30000, "Belanja", at(1, 0)),            // Thu 1 Oct 00:00: today
		tx(parser.TypeExpense, 15000, "Makanan & Minuman", at(1, 12)), // today
	)
	// Someone else's spending must never show up (SEC-02).
	seed(t, storagetest.NewUser(t), tx(parser.TypeExpense, 777, "Lainnya", at(1, 12)))

	got, err := report.NewStore(storagetest.DB(t)).Summary(ctx, user, at(1, 18))
	if err != nil {
		t.Fatal(err)
	}
	want := []report.Period{
		{Label: "Hari ini", Expense: 45000, Income: 0},
		{Label: "Minggu ini", Expense: 65000, Income: 500000},
		{Label: "Bulan ini", Expense: 45000, Income: 0},
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("period %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestSum(t *testing.T) {
	ctx := context.Background()
	store := report.NewStore(storagetest.DB(t))
	user := storagetest.NewUser(t)
	seed(t, user,
		tx(parser.TypeExpense, 15000, "Makanan & Minuman", at(5, 12)),
		tx(parser.TypeExpense, 8000, "Makanan & Minuman", at(6, 12)),
		tx(parser.TypeExpense, 20000, "Transportasi", at(6, 13)),
		tx(parser.TypeIncome, 500000, "Uang Saku", at(6, 9)),
		tx(parser.TypeExpense, 1000, "Lainnya", at(7, 0)), // outside the range (End is exclusive)
	)
	seed(t, storagetest.NewUser(t), tx(parser.TypeExpense, 777, "Makanan & Minuman", at(6, 12)))

	october := parser.DateRange{Start: at(1, 0), End: at(7, 0)}
	cases := []struct {
		name string
		q    parser.ReportQuery
		want int64
	}{
		{"all expenses", parser.ReportQuery{Type: parser.TypeExpense, Range: october}, 43000},
		{"one category", parser.ReportQuery{Type: parser.TypeExpense, Range: october, Categories: []string{"Makanan & Minuman"}}, 23000},
		{"income", parser.ReportQuery{Type: parser.TypeIncome, Range: october}, 500000},
	}
	for _, c := range cases {
		got, err := store.Sum(ctx, user, c.q)
		if err != nil {
			t.Fatal(err)
		}
		if got != c.want {
			t.Errorf("%s: Sum = %d, want %d", c.name, got, c.want)
		}
	}
}
