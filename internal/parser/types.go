// Package parser turns natural-language messages into structured results.
// The first layer uses rules (regex + dictionaries); if they fail, the second
// layer calls an LLM whose output is always validated (see llm.go).
package parser

import "time"

type TransactionType int

const (
	TypeUnspecified TransactionType = iota
	TypeExpense
	TypeIncome
	TypeTransfer
)

func (t TransactionType) String() string {
	switch t {
	case TypeExpense:
		return "EXPENSE"
	case TypeIncome:
		return "INCOME"
	case TypeTransfer:
		return "TRANSFER"
	default:
		return "UNSPECIFIED"
	}
}

type Source int

const (
	SourceRule Source = iota + 1
	SourceLLM
)

func (s Source) String() string {
	if s == SourceLLM {
		return "LLM"
	}
	return "RULE"
}

type Kind int

const (
	KindTransactions Kind = iota + 1
	KindQuery
	KindCorrection
	KindClarification
	KindSmallTalk
)

// ConfidenceThreshold: below this value the bot asks back (FR-12).
const ConfidenceThreshold = 0.7

type ParsedTransaction struct {
	Type        TransactionType
	Amount      int64 // rupiah, whole number
	Description string
	Category    string
	OccurredAt  time.Time
	Confidence  float64
}

// RecentTransaction is context for corrections like "yang bakso tadi".
type RecentTransaction struct {
	ID          string
	Description string
	Amount      int64
	Category    string
	OccurredAt  time.Time
}

type Metric int

const (
	MetricSum Metric = iota + 1
	MetricCount
	MetricTopN
	MetricBreakdown
)

type DateRange struct {
	Start time.Time // inclusive
	End   time.Time // exclusive
}

type ReportQuery struct {
	Metric     Metric
	Type       TransactionType
	Categories []string
	Range      DateRange
	Limit      int
}

type Correction struct {
	TransactionID string
	NewAmount     int64
}

type Request struct {
	UserID     string
	Text       string
	Location   *time.Location
	ReceivedAt time.Time
	Context    []RecentTransaction
}

type Result struct {
	Kind          Kind
	Transactions  []ParsedTransaction
	Query         *ReportQuery
	Correction    *Correction
	Clarification string
	Source        Source
}
