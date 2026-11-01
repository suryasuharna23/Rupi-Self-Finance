//go:build integration

// Package storagetest starts a real PostgreSQL in a container for store
// integration tests (RUPI-102). Use it from TestMain:
//
//	func TestMain(m *testing.M) { os.Exit(storagetest.Run(m)) }
//
// and call storagetest.DB(t) in each test. Run with:
//
//	go test -tags integration ./...
package storagetest

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/suryasuharna23/rupi/internal/storage"
)

const image = "postgres:16-alpine"

var (
	pool       *pgxpool.Pool
	telegramID atomic.Int64
)

// Run starts the container, applies migrations, runs the tests, and tears
// everything down. It returns the exit code for os.Exit.
func Run(m *testing.M) int {
	ctx := context.Background()
	ctr, err := postgres.Run(ctx, image,
		postgres.WithDatabase("rupi"),
		postgres.WithUsername("rupi"),
		postgres.WithPassword("rupi"),
		postgres.BasicWaitStrategies(),
	)
	defer testcontainers.TerminateContainer(ctr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start postgres container: %v\n", err)
		return 1
	}

	url, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "postgres connection string: %v\n", err)
		return 1
	}
	if err := storage.Migrate(url); err != nil {
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
		return 1
	}
	pool, err = storage.Connect(ctx, url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect: %v\n", err)
		return 1
	}
	defer pool.Close()

	return m.Run()
}

// DB returns the shared pool. Tests stay isolated by creating their own
// users rather than truncating tables, so they may run in parallel.
func DB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	if pool == nil {
		t.Fatal("storagetest.Run was not called from TestMain")
	}
	return pool
}

// TelegramID returns a Telegram ID no other test in this package uses
// (each test package gets its own container).
func TelegramID() int64 { return telegramID.Add(1) }

// NewUser inserts a fresh user and returns its ID.
func NewUser(t *testing.T) string {
	t.Helper()
	var id string
	err := DB(t).QueryRow(context.Background(),
		`INSERT INTO user_svc.users (telegram_id) VALUES ($1) RETURNING id`, TelegramID(),
	).Scan(&id)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
	return id
}
