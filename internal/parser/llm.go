package parser

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"
)

// LLM is the parser's second layer. Implementations receive only the message
// text and the category list (SEC-04), and only return structured data: the
// LLM never writes to the database or performs actions (SEC-06).
//
// The LLM provider is still an open question in the requirements document, so
// there is no concrete implementation yet.
type LLM interface {
	Parse(ctx context.Context, text string, now time.Time, categories []string) (Result, error)
}

var ErrInvalidLLMOutput = errors.New("invalid LLM output")

// Validate checks LLM output against sane limits (SEC-05).
func Validate(res Result, now time.Time) error {
	switch res.Kind {
	case KindTransactions:
		if len(res.Transactions) == 0 {
			return fmt.Errorf("%w: empty transaction list", ErrInvalidLLMOutput)
		}
		for _, tx := range res.Transactions {
			if tx.Amount < MinAmount || tx.Amount > MaxAmount {
				return fmt.Errorf("%w: amount %d out of range", ErrInvalidLLMOutput, tx.Amount)
			}
			if tx.Type != TypeExpense && tx.Type != TypeIncome {
				return fmt.Errorf("%w: transaction type %s", ErrInvalidLLMOutput, tx.Type)
			}
			if tx.OccurredAt.After(now.Add(time.Minute)) || tx.OccurredAt.Before(now.AddDate(-1, 0, 0)) {
				return fmt.Errorf("%w: implausible date %s", ErrInvalidLLMOutput, tx.OccurredAt)
			}
			valid := ExpenseCategories
			if tx.Type == TypeIncome {
				valid = IncomeCategories
			}
			if !slices.Contains(valid, tx.Category) {
				return fmt.Errorf("%w: category %q", ErrInvalidLLMOutput, tx.Category)
			}
		}
	case KindQuery:
		if res.Query == nil || !res.Query.Range.Start.Before(res.Query.Range.End) {
			return fmt.Errorf("%w: incomplete query", ErrInvalidLLMOutput)
		}
	case KindClarification, KindSmallTalk:
	default:
		// Corrections always go through the rule path because they need a transaction ID.
		return fmt.Errorf("%w: result kind %d", ErrInvalidLLMOutput, res.Kind)
	}
	return nil
}
