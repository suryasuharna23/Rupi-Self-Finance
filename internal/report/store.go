// Package report runs transaction aggregations (the future Report Service in Phase 2).
package report

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/suryasuharna23/rupi/internal/parser"
)

type Store struct{ db *pgxpool.Pool }

func NewStore(db *pgxpool.Pool) *Store { return &Store{db: db} }

// Sum answers a ReportQuery with the SUM metric (FR-20). The user_id filter
// always comes from the authenticated user, never from message content (SEC-02).
func (s *Store) Sum(ctx context.Context, userID string, q parser.ReportQuery) (int64, error) {
	var total int64
	err := s.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(t.amount), 0)
		FROM transaction_svc.transactions t
		JOIN transaction_svc.categories c ON c.id = t.category_id
		WHERE t.user_id = $1 AND t.deleted_at IS NULL AND t.type = $2
		  AND t.occurred_at >= $3 AND t.occurred_at < $4
		  AND (COALESCE(cardinality($5::text[]), 0) = 0 OR c.name = ANY($5))`, // a nil slice is sent as NULL
		userID, int(q.Type), q.Range.Start, q.Range.End, q.Categories,
	).Scan(&total)
	return total, err
}

type Period struct {
	Label   string
	Expense int64
	Income  int64
}

// Summary computes totals for today, this week, and this month (FR-21).
func (s *Store) Summary(ctx context.Context, userID string, now time.Time) ([]Period, error) {
	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	week := day.AddDate(0, 0, -((int(now.Weekday()) + 6) % 7)) // Monday
	month := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	periods := []Period{{Label: "Hari ini"}, {Label: "Minggu ini"}, {Label: "Bulan ini"}}
	err := s.db.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(amount) FILTER (WHERE type = 1 AND occurred_at >= $2), 0),
			COALESCE(SUM(amount) FILTER (WHERE type = 2 AND occurred_at >= $2), 0),
			COALESCE(SUM(amount) FILTER (WHERE type = 1 AND occurred_at >= $3), 0),
			COALESCE(SUM(amount) FILTER (WHERE type = 2 AND occurred_at >= $3), 0),
			COALESCE(SUM(amount) FILTER (WHERE type = 1 AND occurred_at >= $4), 0),
			COALESCE(SUM(amount) FILTER (WHERE type = 2 AND occurred_at >= $4), 0)
		FROM transaction_svc.transactions
		WHERE user_id = $1 AND deleted_at IS NULL AND occurred_at >= LEAST($3, $4)`,
		userID, day, week, month,
	).Scan(&periods[0].Expense, &periods[0].Income, &periods[1].Expense, &periods[1].Income, &periods[2].Expense, &periods[2].Income)
	return periods, err
}
