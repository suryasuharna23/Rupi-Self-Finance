package parser

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// LogRetention is how long parser_logs rows are kept.
const LogRetention = 30 * 24 * time.Hour

// Values of LogEntry.Source.
const (
	LogSourceRule      = "RULE"       // the rules understood the message
	LogSourceLLM       = "LLM"        // the LLM understood the message
	LogSourceLLMFailed = "LLM_FAILED" // the LLM failed or its output was invalid
	LogSourceNoLLM     = "NO_LLM"     // the rules failed and no LLM is configured
)

// LogEntry is one parser_logs row. It never holds the raw message text or
// the raw user ID (SEC-08).
type LogEntry struct {
	UserHash     string
	TextRedacted string
	Result       Result // transaction descriptions are redacted too
	Source       string
	Latency      time.Duration
	Tokens       int // filled in once the LLM reports usage (RUPI-106)
}

type LogSink interface {
	WriteLog(ctx context.Context, e LogEntry) error
}

// HashUserID is the user_id_hash stored in parser_logs, so the rows of a
// deleted account can still be found and purged.
func HashUserID(userID string) string {
	sum := sha256.Sum256([]byte(userID))
	return hex.EncodeToString(sum[:])
}

func (p *Parser) writeLog(ctx context.Context, req Request, res Result, source string, latency time.Duration) {
	if p.logs == nil {
		return
	}
	logged := res
	logged.Transactions = slices.Clone(res.Transactions)
	for i := range logged.Transactions {
		logged.Transactions[i].Description = Redact(logged.Transactions[i].Description)
	}
	err := p.logs.WriteLog(ctx, LogEntry{
		UserHash:     HashUserID(req.UserID),
		TextRedacted: Redact(req.Text),
		Result:       logged,
		Source:       source,
		Latency:      latency,
	})
	if err != nil {
		slog.WarnContext(ctx, "write parser log failed", "err", err)
	}
}

// LogStore keeps parser_logs in PostgreSQL (the future Parser Service schema).
type LogStore struct{ db *pgxpool.Pool }

func NewLogStore(db *pgxpool.Pool) *LogStore { return &LogStore{db: db} }

func (s *LogStore) WriteLog(ctx context.Context, e LogEntry) error {
	result, err := json.Marshal(e.Result)
	if err != nil {
		return fmt.Errorf("encode parser result: %w", err)
	}
	_, err = s.db.Exec(ctx, `
		INSERT INTO parser_svc.parser_logs (user_id_hash, text_redacted, result, source, latency_ms, tokens)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		e.UserHash, e.TextRedacted, result, e.Source, e.Latency.Milliseconds(), e.Tokens)
	return err
}

// DeleteOlderThan removes rows created before cutoff and returns how many.
func (s *LogStore) DeleteOlderThan(ctx context.Context, cutoff time.Time) (int64, error) {
	tag, err := s.db.Exec(ctx, `DELETE FROM parser_svc.parser_logs WHERE created_at < $1`, cutoff)
	return tag.RowsAffected(), err
}
