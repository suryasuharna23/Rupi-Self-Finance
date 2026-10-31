package parser

import (
	"regexp"
	"strings"
	"time"
)

var (
	smallTalkRe        = regexp.MustCompile(`(?i)^\s*(halo|hai|hi|hello|hey|pagi|siang|sore|malam|makasih|terima\s*kasih|thanks|thx|tq|ok|oke|sip)\b`)
	segmentSepRe       = regexp.MustCompile(`(?i)\s+(?:sama|dan|terus|trus|plus|\+)\s+|,\s+|;|\n`)
	fillerRe           = regexp.MustCompile(`(?i)\b(beli|bayar|jajan|abis|habis|buat|untuk|utk|dapet|dapat|sebesar|seharga|harga)\b`)
	spacesRe           = regexp.MustCompile(`\s+`)
	correctionRe       = regexp.MustCompile(`(?i)^\s*yang\s+`)
	correctionFillerRe = regexp.MustCompile(`(?i)\b(yang|tadi|harusnya|seharusnya|jadi|diganti|ganti|ubah)\b`)
)

var (
	queryWords   = []string{"berapa", "total", "rekap", "ringkasan"}
	expenseVerbs = []string{"beli", "bayar", "jajan", "isi", "makan", "minum"}
)

// ParseRules tries to understand a message using rules only.
// ok=false means the message must be passed on to the LLM.
func ParseRules(req Request) (Result, bool) {
	text := strings.TrimSpace(req.Text)
	now := req.ReceivedAt.In(location(req))
	amounts := FindAmounts(text)

	switch {
	case containsAny(text, queryWords):
		return Result{Kind: KindQuery, Query: buildQuery(text, now), Source: SourceRule}, true
	case correctionRe.MatchString(text) && len(amounts) == 1:
		return parseCorrection(text, amounts[0], req.Context), true
	case len(amounts) == 0 && smallTalkRe.MatchString(text):
		return Result{Kind: KindSmallTalk, Source: SourceRule}, true
	case len(amounts) == 0 && (containsAny(text, expenseVerbs) || IsIncome(text)):
		return clarify("Berapa nominalnya?"), true
	case len(amounts) == 0:
		return Result{}, false
	}

	occurredAt, rest := ResolveTime(text, now)
	var txs []ParsedTransaction
	pending := ""
	for _, seg := range segmentSepRe.Split(rest, -1) {
		seg = strings.TrimSpace(pending + " " + seg)
		segAmounts := FindAmounts(seg)
		if len(segAmounts) == 0 {
			pending = seg // merge with the next segment: "makan sama temen 20rb"
			continue
		}
		pending = ""
		txs = append(txs, buildTransaction(seg, segAmounts, occurredAt))
	}
	if len(txs) == 0 {
		return Result{}, false
	}
	return Result{Kind: KindTransactions, Transactions: txs, Source: SourceRule}, true
}

func buildTransaction(seg string, amounts []Amount, occurredAt time.Time) ParsedTransaction {
	a := amounts[0]
	typ := TypeExpense
	if IsIncome(seg) {
		typ = TypeIncome
	}
	desc := seg[:a.Start] + " " + seg[a.End:]
	category, known := Categorize(desc, typ)
	desc = cleanDescription(fillerRe.ReplaceAllString(desc, " "))

	confidence := 0.95
	switch {
	case a.Ambiguous || len(amounts) > 1:
		confidence = 0.5
	case desc == "":
		confidence = 0.6
	case !known:
		confidence = 0.8
	}
	return ParsedTransaction{
		Type:        typ,
		Amount:      a.Value,
		Description: desc,
		Category:    category,
		OccurredAt:  occurredAt,
		Confidence:  confidence,
	}
}

func parseCorrection(text string, a Amount, recent []RecentTransaction) Result {
	words := strings.Fields(strings.ToLower(correctionFillerRe.ReplaceAllString(text[:a.Start]+" "+text[a.End:], " ")))
	for _, tx := range recent { // recent is ordered newest first
		desc := strings.ToLower(tx.Description)
		for _, w := range words {
			if len(w) >= 3 && strings.Contains(desc, w) {
				return Result{
					Kind:       KindCorrection,
					Correction: &Correction{TransactionID: tx.ID, NewAmount: a.Value},
					Source:     SourceRule,
				}
			}
		}
	}
	return clarify("Transaksi mana yang mau diubah? Belum ketemu di catatan terakhirmu.")
}

func buildQuery(text string, now time.Time) *ReportQuery {
	q := &ReportQuery{Metric: MetricSum, Type: TypeExpense, Range: monthRange(now)}
	lower := strings.ToLower(text)
	switch {
	case strings.Contains(lower, "hari ini"):
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		q.Range = DateRange{Start: start, End: start.AddDate(0, 0, 1)}
	case strings.Contains(lower, "minggu ini"):
		q.Range = weekRange(now)
	case strings.Contains(lower, "bulan lalu"):
		r := monthRange(now)
		q.Range = DateRange{Start: r.Start.AddDate(0, -1, 0), End: r.Start}
	}
	if containsAny(text, []string{"pemasukan", "dapet", "dapat", "masuk"}) {
		q.Type = TypeIncome
	}
	if cat, ok := Categorize(text, q.Type); ok {
		q.Categories = []string{cat}
	}
	return q
}

func monthRange(now time.Time) DateRange {
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return DateRange{Start: start, End: start.AddDate(0, 1, 0)}
}

// weekRange treats Monday as the start of the week.
func weekRange(now time.Time) DateRange {
	offset := (int(now.Weekday()) + 6) % 7
	start := time.Date(now.Year(), now.Month(), now.Day()-offset, 0, 0, 0, 0, now.Location())
	return DateRange{Start: start, End: start.AddDate(0, 0, 7)}
}

func cleanDescription(s string) string {
	return strings.Trim(spacesRe.ReplaceAllString(s, " "), " -:")
}

func clarify(question string) Result {
	return Result{Kind: KindClarification, Clarification: question, Source: SourceRule}
}

func location(req Request) *time.Location {
	if req.Location != nil {
		return req.Location
	}
	return time.UTC
}
