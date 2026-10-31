// Command rupi runs the Phase 1 Rupi bot: a single process containing the Channel
// Adapter, Conversation, Parser, Transaction, User, and Report modules.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata" // time zones stay available on minimal images and Windows

	"github.com/suryasuharna23/rupi/internal/config"
	"github.com/suryasuharna23/rupi/internal/conversation"
	"github.com/suryasuharna23/rupi/internal/parser"
	"github.com/suryasuharna23/rupi/internal/report"
	"github.com/suryasuharna23/rupi/internal/storage"
	"github.com/suryasuharna23/rupi/internal/telegram"
	"github.com/suryasuharna23/rupi/internal/transaction"
	"github.com/suryasuharna23/rupi/internal/user"
)

// shutdownTimeout follows the NFR-10 graceful shutdown limit.
const shutdownTimeout = 30 * time.Second

func main() {
	if err := run(); err != nil {
		slog.Error("rupi stopped", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	var level slog.Level
	if err := level.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		return err
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := storage.Migrate(cfg.DatabaseURL); err != nil {
		return err
	}
	db, err := storage.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	svc := conversation.NewService(
		user.NewStore(db),
		transaction.NewStore(db),
		report.NewStore(db),
		conversation.NewIdempotencyStore(db),
		parser.New(nil), // no LLM yet: messages the rules don't understand get a clarifying question
	)
	tg := telegram.NewClient(cfg.TelegramAPIURL, cfg.TelegramToken)
	adapter := telegram.NewAdapter(tg, svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(r.Context()); err != nil {
			http.Error(w, "db unavailable", http.StatusServiceUnavailable)
			return
		}
		w.Write([]byte("ok"))
	})
	if cfg.TelegramMode == "webhook" {
		mux.Handle("POST /telegram/webhook", adapter.WebhookHandler(cfg.WebhookSecret))
		if err := tg.SetWebhook(ctx, cfg.WebhookURL, cfg.WebhookSecret); err != nil {
			return err
		}
	} else {
		go func() {
			if err := adapter.Poll(ctx); err != nil {
				slog.Error("polling stopped", "err", err)
				stop()
			}
		}()
	}

	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: mux, ReadHeaderTimeout: 10 * time.Second}
	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()
	slog.Info("rupi running", "mode", cfg.TelegramMode, "addr", cfg.HTTPAddr)

	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
