// Package storage menyiapkan koneksi PostgreSQL dan migrasi skema.
package storage

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Connect(ctx context.Context, url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("buka pool postgres: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return pool, nil
}

// Migrate menjalankan semua migrasi yang belum diterapkan.
func Migrate(url string) error {
	src, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}
	// Driver golang-migrate untuk pgx v5 memakai skema URL "pgx5".
	dbURL := "pgx5://" + strings.TrimPrefix(strings.TrimPrefix(url, "postgres://"), "postgresql://")
	m, err := migrate.NewWithSourceInstance("iofs", src, dbURL)
	if err != nil {
		return fmt.Errorf("siapkan migrasi: %w", err)
	}
	defer m.Close()
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("jalankan migrasi: %w", err)
	}
	return nil
}
