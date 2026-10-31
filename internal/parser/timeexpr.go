package parser

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var weekdays = map[string]time.Weekday{
	"minggu": time.Sunday, "senin": time.Monday, "selasa": time.Tuesday,
	"rabu": time.Wednesday, "kamis": time.Thursday, "jumat": time.Friday,
	"jum'at": time.Friday, "sabtu": time.Saturday,
}

var partOfDayHour = map[string]int{"pagi": 8, "siang": 12, "sore": 16, "malam": 19}

type timeRule struct {
	re      *regexp.Regexp
	resolve func(m []string, now time.Time) time.Time
}

// Order matters: more specific patterns are tried first.
var timeRules = []timeRule{
	{regexp.MustCompile(`(?i)\bkemarin\s+lusa\b`), func(_ []string, now time.Time) time.Time {
		return atNoon(now.AddDate(0, 0, -2))
	}},
	// "minggu lalu" means last week, not last Sunday.
	{regexp.MustCompile(`(?i)\bminggu\s+lalu\b`), func(_ []string, now time.Time) time.Time {
		return atNoon(now.AddDate(0, 0, -7))
	}},
	{regexp.MustCompile(`(?i)\b(senin|selasa|rabu|kamis|jum'?at|sabtu|minggu)\s+(?:kemarin|lalu)\b`), func(m []string, now time.Time) time.Time {
		target := weekdays[strings.ToLower(m[1])]
		diff := (int(now.Weekday()) - int(target) + 7) % 7
		if diff == 0 {
			diff = 7
		}
		return atNoon(now.AddDate(0, 0, -diff))
	}},
	{regexp.MustCompile(`(?i)\bkemarin\b`), func(_ []string, now time.Time) time.Time {
		return atNoon(now.AddDate(0, 0, -1))
	}},
	{regexp.MustCompile(`(?i)\b(?:tgl\.?|tanggal)\s*(\d{1,2})\b`), func(m []string, now time.Time) time.Time {
		day, _ := strconv.Atoi(m[1])
		year, month := now.Year(), now.Month()
		if day > now.Day() {
			month--
		}
		// Clamp to the month's last day: "tgl 31" in a 30-day month.
		last := time.Date(year, month+1, 0, 0, 0, 0, 0, now.Location()).Day()
		return time.Date(year, month, min(day, last), 12, 0, 0, 0, now.Location())
	}},
	{regexp.MustCompile(`(?i)\btadi\s+(pagi|siang|sore|malam)\b`), func(m []string, now time.Time) time.Time {
		h := partOfDayHour[strings.ToLower(m[1])]
		t := time.Date(now.Year(), now.Month(), now.Day(), h, 0, 0, 0, now.Location())
		if t.After(now) {
			t = t.AddDate(0, 0, -1) // "tadi malam" sent after midnight
		}
		return t
	}},
	{regexp.MustCompile(`(?i)\btadi\b`), func(_ []string, now time.Time) time.Time { return now }},
}

// ResolveTime finds the first relative time expression in the text.
// It returns when the transaction happened and the text with the expression removed.
// Without a time expression, the result is now.
func ResolveTime(text string, now time.Time) (time.Time, string) {
	for _, r := range timeRules {
		loc := r.re.FindStringSubmatchIndex(text)
		if loc == nil {
			continue
		}
		m := make([]string, len(loc)/2)
		for i := range m {
			if loc[2*i] >= 0 {
				m[i] = text[loc[2*i]:loc[2*i+1]]
			}
		}
		t := r.resolve(m, now)
		if t.After(now) {
			t = now
		}
		return t, text[:loc[0]] + " " + text[loc[1]:]
	}
	return now, text
}

func atNoon(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 12, 0, 0, 0, t.Location())
}
