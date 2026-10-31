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
	llm LLM // may be nil: the parser then uses rules only
}

func New(llm LLM) *Parser { return &Parser{llm: llm} }

// Parse always returns a result; LLM failures become a clarifying question.
func (p *Parser) Parse(ctx context.Context, req Request) Result {
	if res, ok := ParseRules(req); ok {
		return res
	}
	if p.llm == nil {
		return clarify(fallbackQuestion)
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
		return clarify(fallbackQuestion)
	}
	res.Source = SourceLLM
	return res
}
