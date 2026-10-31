package parser

import (
	"regexp"
	"strconv"
	"strings"
)

// Amount is an amount found in the text, along with its position.
type Amount struct {
	Value int64
	// Ambiguous is true when the amount needs confirmation, e.g. "15" without
	// a unit, or old slang like "cepek".
	Ambiguous  bool
	Start, End int
}

const (
	MinAmount int64 = 1
	MaxAmount int64 = 1_000_000_000
)

var (
	numericAmountRe = regexp.MustCompile(`(?i)(?:\brp\.?\s*|\b)(\d+(?:[.,]\d+)*)(?:\s*(rb|ribu|k|jt|juta)\b)?`)
	slangAmountRe   = regexp.MustCompile(`(?i)\b(setengah\s+juta|sejuta|goceng|ceban|goban|cepek|gocap)\b`)
)

var slangValues = map[string]Amount{
	"goceng":        {Value: 5_000},
	"ceban":         {Value: 10_000},
	"goban":         {Value: 50_000},
	"cepek":         {Value: 100, Ambiguous: true},
	"gocap":         {Value: 50, Ambiguous: true},
	"sejuta":        {Value: 1_000_000},
	"setengah juta": {Value: 500_000},
}

var multipliers = map[string]int64{
	"":     1,
	"rb":   1_000,
	"ribu": 1_000,
	"k":    1_000,
	"jt":   1_000_000,
	"juta": 1_000_000,
}

// FindAmounts returns every amount in the text, ordered by position.
func FindAmounts(text string) []Amount {
	var out []Amount
	for _, m := range slangAmountRe.FindAllStringSubmatchIndex(text, -1) {
		key := strings.Join(strings.Fields(strings.ToLower(text[m[2]:m[3]])), " ")
		a := slangValues[key]
		a.Start, a.End = m[0], m[1]
		out = append(out, a)
	}
	for _, m := range numericAmountRe.FindAllStringSubmatchIndex(text, -1) {
		suffix := ""
		if m[4] >= 0 {
			suffix = strings.ToLower(text[m[4]:m[5]])
		}
		hasRp := strings.HasPrefix(strings.ToLower(text[m[0]:m[1]]), "rp")
		value, ambiguous, ok := parseNumber(text[m[2]:m[3]], suffix, hasRp)
		if !ok {
			continue
		}
		out = insertSorted(out, Amount{Value: value, Ambiguous: ambiguous, Start: m[0], End: m[1]})
	}
	return out
}

// parseNumber handles "15000", "15.000", "15,000", "1,5" + unit, and so on.
func parseNumber(num, suffix string, hasRp bool) (value int64, ambiguous, ok bool) {
	mult := multipliers[suffix]
	intPart, fracPart := num, ""

	groups := strings.FieldsFunc(num, func(r rune) bool { return r == '.' || r == ',' })
	if len(groups) > 1 {
		if isThousandsGrouping(groups) && !(suffix != "" && len(groups) == 2) {
			intPart = strings.Join(groups, "")
		} else if len(groups) == 2 && suffix != "" {
			intPart, fracPart = groups[0], groups[1]
		} else {
			return 0, false, false
		}
	}

	i, err := strconv.ParseInt(intPart, 10, 64)
	if err != nil {
		return 0, false, false
	}
	value = i * mult
	if fracPart != "" {
		f, err := strconv.ParseInt(fracPart, 10, 64)
		if err != nil {
			return 0, false, false
		}
		value += f * mult / pow10(len(fracPart))
	}
	if value < MinAmount || value > MaxAmount {
		return 0, false, false
	}
	// A small number without a unit ("15") could mean Rp15 or Rp15.000.
	ambiguous = suffix == "" && !hasRp && value < 1_000
	return value, ambiguous, true
}

func isThousandsGrouping(groups []string) bool {
	if len(groups[0]) == 0 || len(groups[0]) > 3 {
		return false
	}
	for _, g := range groups[1:] {
		if len(g) != 3 {
			return false
		}
	}
	return true
}

func pow10(n int) int64 {
	p := int64(1)
	for range n {
		p *= 10
	}
	return p
}

func insertSorted(list []Amount, a Amount) []Amount {
	for i, x := range list {
		if a.Start < x.End && x.Start < a.End {
			return list // overlaps slang that was already found
		}
		if a.Start < x.Start {
			return append(list[:i], append([]Amount{a}, list[i:]...)...)
		}
	}
	return append(list, a)
}
