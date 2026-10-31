// Package user memiliki data profil pengguna (calon User Service di Fase 2).
package user

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID          string
	TelegramID  int64
	DisplayName string
	Timezone    string
	Onboarded   bool
}

type Store struct{ db *pgxpool.Pool }

func NewStore(db *pgxpool.Pool) *Store { return &Store{db: db} }

// GetOrCreate mendaftarkan pengguna otomatis berdasarkan ID Telegram (FR-01).
func (s *Store) GetOrCreate(ctx context.Context, telegramID int64) (User, bool, error) {
	var u User
	var created bool
	err := s.db.QueryRow(ctx, `
		WITH ins AS (
			INSERT INTO user_svc.users (telegram_id) VALUES ($1)
			ON CONFLICT (telegram_id) WHERE deleted_at IS NULL DO NOTHING
			RETURNING id, telegram_id, display_name, timezone, onboarded, true AS created
		)
		SELECT * FROM ins
		UNION ALL
		SELECT id, telegram_id, display_name, timezone, onboarded, false
		FROM user_svc.users WHERE telegram_id = $1 AND deleted_at IS NULL
		LIMIT 1`, telegramID,
	).Scan(&u.ID, &u.TelegramID, &u.DisplayName, &u.Timezone, &u.Onboarded, &created)
	if err != nil {
		return User{}, false, fmt.Errorf("get or create user: %w", err)
	}
	return u, created, nil
}

func (s *Store) CompleteOnboarding(ctx context.Context, userID, displayName string) error {
	_, err := s.db.Exec(ctx,
		`UPDATE user_svc.users SET display_name = $2, onboarded = true WHERE id = $1`,
		userID, displayName)
	return err
}

// Delete melakukan soft delete; data dipurge permanen setelah 30 hari (SEC-10).
func (s *Store) Delete(ctx context.Context, userID string) error {
	_, err := s.db.Exec(ctx,
		`UPDATE user_svc.users SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`, userID)
	return err
}
