// Package parser mengubah pesan bahasa natural menjadi hasil terstruktur.
// Lapis pertama memakai aturan (regex + kamus); bila gagal, lapis kedua
// memanggil LLM yang keluarannya selalu divalidasi (lihat llm.go).
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

// ConfidenceThreshold: di bawah nilai ini bot bertanya balik (FR-12).
const ConfidenceThreshold = 0.7

type ParsedTransaction struct {
	Type        TransactionType
	Amount      int64 // rupiah, bilangan bulat
	Description string
	Category    string
	OccurredAt  time.Time
	Confidence  float64
}

// RecentTransaction adalah konteks untuk koreksi seperti "yang bakso tadi".
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
	Start time.Time // inklusif
	End   time.Time // eksklusif
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
