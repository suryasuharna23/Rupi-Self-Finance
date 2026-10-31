// Package transaction memiliki data transaksi dan kategori (calon
// Transaction Service di Fase 2). Setiap query wajib difilter user_id (SEC-02).
package transaction

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/suryasuharna23/rupi/internal/parser"
)

var ErrNotFound = errors.New("transaksi tidak ditemukan")

// UndoWindow adalah batas waktu /undo (FR-15).
const UndoWindow = 24 * time.Hour

type Transaction struct {
	ID          string
	Type        parser.TransactionType
	Amount      int64
	Description string
	Category    string
	OccurredAt  time.Time
}

type Store struct{ db *pgxpool.Pool }

func NewStore(db *pgxpool.Pool) *Store { return &Store{db: db} }

// Create menyimpan semua transaksi dari satu pesan secara atomik.
func (s *Store) Create(ctx context.Context, userID, sourceMessageID string, source parser.Source, txs []parser.ParsedTransaction) ([]Transaction, error) {
	out := make([]Transaction, 0, len(txs))
	err := pgx.BeginFunc(ctx, s.db, func(tx pgx.Tx) error {
		for _, p := range txs {
			var id string
			err := tx.QueryRow(ctx, `
				INSERT INTO transaction_svc.transactions
					(user_id, type, amount, description, category_id, occurred_at, source_message_id, parser_source)
				VALUES ($1, $2, $3, $4, (
					SELECT id FROM transaction_svc.categories
					WHERE name = $5 AND type = $2 AND (user_id IS NULL OR user_id = $1)
					ORDER BY user_id NULLS LAST LIMIT 1
				), $6, $7, $8)
				RETURNING id`,
				userID, int(p.Type), p.Amount, p.Description, p.Category, p.OccurredAt, sourceMessageID, source.String(),
			).Scan(&id)
			if err != nil {
				return fmt.Errorf("insert transaksi: %w", err)
			}
			out = append(out, Transaction{ID: id, Type: p.Type, Amount: p.Amount, Description: p.Description, Category: p.Category, OccurredAt: p.OccurredAt})
		}
		return nil
	})
	return out, err
}

// UpdateAmount mengoreksi nominal dan mencatat jejak audit (FR-14).
func (s *Store) UpdateAmount(ctx context.Context, userID, txID string, amount int64) (Transaction, error) {
	var t Transaction
	err := pgx.BeginFunc(ctx, s.db, func(tx pgx.Tx) error {
		var old int64
		err := tx.QueryRow(ctx, `
			SELECT amount FROM transaction_svc.transactions
			WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL FOR UPDATE`, txID, userID).Scan(&old)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		err = tx.QueryRow(ctx, `
			UPDATE transaction_svc.transactions t SET amount = $3, updated_at = now()
			FROM transaction_svc.categories c
			WHERE t.id = $1 AND t.user_id = $2 AND c.id = t.category_id
			RETURNING t.id, t.type, t.amount, t.description, c.name, t.occurred_at`, txID, userID, amount,
		).Scan(&t.ID, &t.Type, &t.Amount, &t.Description, &t.Category, &t.OccurredAt)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO transaction_svc.transaction_edits (transaction_id, field, old_value, new_value)
			VALUES ($1, 'amount', $2, $3)`, txID, strconv.FormatInt(old, 10), strconv.FormatInt(amount, 10))
		return err
	})
	return t, err
}

// DeleteByMessage menghapus transaksi dari satu pesan (tombol Batalkan, FR-13).
func (s *Store) DeleteByMessage(ctx context.Context, userID, sourceMessageID string) (int64, error) {
	tag, err := s.db.Exec(ctx, `
		UPDATE transaction_svc.transactions SET deleted_at = now()
		WHERE user_id = $1 AND source_message_id = $2 AND deleted_at IS NULL`, userID, sourceMessageID)
	return tag.RowsAffected(), err
}

// DeleteLast menghapus transaksi terakhir yang dibuat dalam UndoWindow (FR-15).
func (s *Store) DeleteLast(ctx context.Context, userID string) (Transaction, error) {
	var t Transaction
	err := s.db.QueryRow(ctx, `
		UPDATE transaction_svc.transactions t SET deleted_at = now()
		FROM transaction_svc.categories c
		WHERE c.id = t.category_id AND t.id = (
			SELECT id FROM transaction_svc.transactions
			WHERE user_id = $1 AND deleted_at IS NULL AND created_at > now() - $2::interval
			ORDER BY created_at DESC LIMIT 1
		)
		RETURNING t.id, t.type, t.amount, t.description, c.name, t.occurred_at`,
		userID, UndoWindow.String(),
	).Scan(&t.ID, &t.Type, &t.Amount, &t.Description, &t.Category, &t.OccurredAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return t, ErrNotFound
	}
	return t, err
}

// FindRecent mengembalikan transaksi terbaru sebagai konteks koreksi.
func (s *Store) FindRecent(ctx context.Context, userID string, limit int) ([]parser.RecentTransaction, error) {
	rows, err := s.db.Query(ctx, `
		SELECT t.id, t.description, t.amount, c.name, t.occurred_at
		FROM transaction_svc.transactions t
		JOIN transaction_svc.categories c ON c.id = t.category_id
		WHERE t.user_id = $1 AND t.deleted_at IS NULL
		ORDER BY t.created_at DESC LIMIT $2`, userID, limit)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (parser.RecentTransaction, error) {
		var t parser.RecentTransaction
		err := r.Scan(&t.ID, &t.Description, &t.Amount, &t.Category, &t.OccurredAt)
		return t, err
	})
}

// DeleteAllForUser menghapus permanen semua transaksi pengguna (FR-03).
func (s *Store) DeleteAllForUser(ctx context.Context, userID string) error {
	return pgx.BeginFunc(ctx, s.db, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `DELETE FROM transaction_svc.transactions WHERE user_id = $1`, userID); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `DELETE FROM transaction_svc.categories WHERE user_id = $1`, userID)
		return err
	})
}
