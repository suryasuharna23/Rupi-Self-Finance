package parser

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"
)

type fakeSink struct {
	entries []LogEntry
	err     error
}

func (f *fakeSink) WriteLog(_ context.Context, e LogEntry) error {
	f.entries = append(f.entries, e)
	return f.err
}

type fakeLLM struct {
	res Result
	err error
}

func (f fakeLLM) Parse(context.Context, string, time.Time, []string) (Result, error) {
	return f.res, f.err
}

func TestParseWritesOneLogPerCall(t *testing.T) {
	llmExpense := Result{Kind: KindTransactions, Transactions: []ParsedTransaction{
		{Type: TypeExpense, Amount: 50_000, Description: "transfer ke Budi", Category: CatOther, OccurredAt: testNow, Confidence: 0.9},
	}}
	cases := []struct {
		name   string
		llm    LLM
		text   string
		source string
	}{
		{"rules", nil, "beli bakso 15rb", LogSourceRule},
		{"no LLM configured", nil, "apa kabar dunia", LogSourceNoLLM},
		{"LLM", fakeLLM{res: llmExpense}, "apa kabar Budi", LogSourceLLM},
		{"LLM error", fakeLLM{err: errors.New("timeout")}, "apa kabar dunia", LogSourceLLMFailed},
		{"invalid LLM output", fakeLLM{res: Result{Kind: KindCorrection}}, "apa kabar dunia", LogSourceLLMFailed},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sink := &fakeSink{}
			p := New(c.llm).WithLogs(sink)
			res := p.Parse(context.Background(), Request{UserID: "user-1", Text: c.text, Location: jakarta, ReceivedAt: testNow})

			if len(sink.entries) != 1 {
				t.Fatalf("%d log entries, want 1", len(sink.entries))
			}
			e := sink.entries[0]
			if e.Source != c.source {
				t.Errorf("source = %s, want %s", e.Source, c.source)
			}
			if e.UserHash != HashUserID("user-1") || strings.Contains(e.UserHash, "user-1") {
				t.Errorf("user hash = %q", e.UserHash)
			}
			if e.TextRedacted != Redact(c.text) {
				t.Errorf("text = %q, want %q", e.TextRedacted, Redact(c.text))
			}
			if e.Result.Kind != res.Kind || len(e.Result.Transactions) != len(res.Transactions) {
				t.Errorf("logged result %+v differs from returned %+v", e.Result, res)
			}
		})
	}
}

func TestLoggedResultIsRedactedButReturnedResultIsNot(t *testing.T) {
	sink := &fakeSink{}
	llm := fakeLLM{res: Result{Kind: KindTransactions, Transactions: []ParsedTransaction{
		{Type: TypeExpense, Amount: 50_000, Description: "transfer ke Budi", Category: CatOther, OccurredAt: testNow, Confidence: 0.9},
	}}}
	res := New(llm).WithLogs(sink).Parse(context.Background(), Request{UserID: "u", Text: "apa kabar Budi", Location: jakarta, ReceivedAt: testNow})

	if got := sink.entries[0].Result.Transactions[0].Description; got != "transfer ke <name>" {
		t.Errorf("logged description = %q", got)
	}
	if got := res.Transactions[0].Description; got != "transfer ke Budi" {
		t.Errorf("returned description = %q, must stay unredacted", got)
	}
}

func TestLogWriteFailureDoesNotBreakParsing(t *testing.T) {
	p := New(nil).WithLogs(&fakeSink{err: errors.New("db down")})
	res := p.Parse(context.Background(), Request{UserID: "u", Text: "beli bakso 15rb", Location: jakarta, ReceivedAt: testNow})
	if res.Kind != KindTransactions {
		t.Errorf("kind = %d, want transactions", res.Kind)
	}
}

// TestLogsDoNotContainMessageText checks SEC-08 at every level, which is
// stricter than the INFO-level requirement.
func TestLogsDoNotContainMessageText(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	t.Cleanup(func() { slog.SetDefault(prev) })

	const secret = "apa kabar Budi rahasia"
	sinks := []LogSink{&fakeSink{}, &fakeSink{err: errors.New("db down")}}
	llms := []LLM{nil, fakeLLM{err: errors.New("timeout")}, fakeLLM{res: Result{Kind: KindCorrection}}}
	for _, sink := range sinks {
		for _, llm := range llms {
			New(llm).WithLogs(sink).Parse(context.Background(), Request{UserID: "u", Text: secret, Location: jakarta, ReceivedAt: testNow})
		}
	}
	if buf.Len() == 0 {
		t.Fatal("expected some log output to check")
	}
	for _, part := range []string{"Budi", "rahasia"} {
		if strings.Contains(buf.String(), part) {
			t.Errorf("logs contain %q:\n%s", part, buf.String())
		}
	}
}
