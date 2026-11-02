package parser

import (
	"context"
	"log/slog"
	"slices"
	"time"
)

// LLMTimeout is how long an LLM call may take before falling back to a clarifying question (NFR-02).
const LLMTimeout = 5 * time.Second

const fallbackQuestion = `Maaf, Rupi belum paham. Coba tulis seperti "beli bakso 15rb" atau "bulan ini habis berapa?"`

type Parser struct {
	llm  LLM     // may be nil: the parser then uses rules only
	logs LogSink // may be nil: nothing is logged
}

func New(llm LLM) *Parser { return &Parser{llm: llm} }

// WithLogs makes every Parse call write one parser_logs entry (SEC-08).
func (p *Parser) WithLogs(logs LogSink) *Parser {
	p.logs = logs
	return p
}

// Parse always returns a result; LLM failures become a clarifying question.
func (p *Parser) Parse(ctx context.Context, req Request) Result {
	start := time.Now()
	res, source := p.parse(ctx, req)
	p.writeLog(ctx, req, res, source, time.Since(start))
	return res
}

func (p *Parser) parse(ctx context.Context, req Request) (Result, string) {
	if res, ok := ParseRules(req); ok {
		return res, LogSourceRule
	}
	if p.llm == nil {
		return clarify(fallbackQuestion), LogSourceNoLLM
	}

	ctx, cancel := context.WithTimeout(ctx, LLMTimeout)
	defer cancel()
	now := req.ReceivedAt.In(location(req))
	res, err := p.llm.Parse(ctx, req.Text, now, slices.Concat(ExpenseCategories, IncomeCategories))
	if err == nil {
		err = Validate(res, now)
	}
	if err != nil {
		slog.WarnContext(ctx, "llm parse failed", "err", err)
		return clarify(fallbackQuestion), LogSourceLLMFailed
	}
	res.Source = SourceLLM
	return res, LogSourceLLM
}
