package parser

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Placeholders that replace personal data in parser_logs (SEC-08).
const (
	RedactedEmail   = "<email>"
	RedactedPhone   = "<phone>"
	RedactedAccount = "<account>"
	RedactedName    = "<name>"
)

var (
	emailRe = regexp.MustCompile(`[\w.+-]+@[\w-]+(?:\.[\w-]+)+`)
	// Indonesian mobile numbers: 08xx, 628xx, or +62 8xx, optionally split by spaces or dashes.
	phoneRe = regexp.MustCompile(`(?:\+62[\s-]?|\b62|\b0)8\d{1,3}(?:[\s-]?\d{2,4}){2,3}\b`)
	// Ten or more digits, optionally dash-separated, are treated as an account
	// number. Amounts stay below this: Rp1 miliar (the maximum) has ten digits
	// but is written with dots or a unit in practice.
	accountRe = regexp.MustCompile(`\b\d(?:-?\d){9,}\b`)
	wordRe    = regexp.MustCompile(`\p{L}+`)
	// "a.n. Budi Santoso" / "atas nama Budi Santoso" on transfer receipts.
	onBehalfRe = regexp.MustCompile(`(?i)\b(?:a\.?n\.?|atas\s+nama)\s+\p{L}+(?:\s+\p{L}+){0,2}`)
)

// honorifics are usually followed by a name: "transfer ke mas budi".
var honorifics = setOf("pak", "bu", "mas", "mbak", "mba", "kak", "bang", "om", "tante", "dek")

// notNames are capitalised or honorific-following words that are not names.
var notNames = setOf(
	"rp", "yang", "uang", "kos", "kost", "kantin", "warung", "toko", "dosen", "rt", "rw",
	"buat", "untuk", "sama", "dan", "ke", "dari", "di", "tadi", "kemarin", "hari", "ini",
	"bayar", "beli", "transfer", "tf", "kirim", "utang", "pinjam", "pinjem",
	"senin", "selasa", "rabu", "kamis", "jumat", "sabtu", "minggu",
	"januari", "februari", "maret", "april", "mei", "juni", "juli", "agustus",
	"september", "oktober", "november", "desember",
)

func init() {
	for _, m := range []map[string][]string{expenseKeywords, incomeKeywords} {
		for _, words := range m {
			for _, w := range words {
				for _, f := range strings.Fields(w) {
					notNames[f] = true
				}
			}
		}
	}
	for _, w := range incomeWords {
		notNames[w] = true
	}
}

// Redact removes personal data from a message before it is stored in
// parser_logs: emails, phone numbers, bank account numbers, and likely names.
// Amounts ("15rb", "Rp1.500.000") are kept so the logs stay useful for
// measuring parser accuracy. Name detection is a heuristic: capitalised words
// that are not at the start of a sentence, words after honorifics ("mas",
// "bu"), and names after "a.n."; a lowercase name without an honorific is missed.
func Redact(text string) string {
	text = emailRe.ReplaceAllString(text, RedactedEmail)
	text = onBehalfRe.ReplaceAllStringFunc(text, redactOnBehalf)
	text = redactNames(text)
	text = phoneRe.ReplaceAllString(text, RedactedPhone)
	return accountRe.ReplaceAllString(text, RedactedAccount)
}

// redactOnBehalf turns "a.n. Andi Wijaya" into "a.n. <name>", stopping at
// the first word that is not a name ("a.n. Rina transfer" keeps "transfer").
func redactOnBehalf(m string) string {
	locs := wordRe.FindAllStringIndex(m, -1)
	// Skip the marker: "a"+"n" or "atas"+"nama" is two words, "an" is one.
	skip := 2
	if strings.EqualFold(m[locs[0][0]:locs[0][1]], "an") {
		skip = 1
	}
	start, end := -1, -1
	for _, loc := range locs[skip:] {
		if notNames[strings.ToLower(m[loc[0]:loc[1]])] {
			break
		}
		if start < 0 {
			start = loc[0]
		}
		end = loc[1]
	}
	if start < 0 {
		return m
	}
	return m[:start] + RedactedName + m[end:]
}

func redactNames(text string) string {
	var b strings.Builder
	last := 0
	prev := "" // previous word, lowercased
	for _, loc := range wordRe.FindAllStringIndex(text, -1) {
		word := text[loc[0]:loc[1]]
		lower := strings.ToLower(word)
		if isPlaceholder(text, loc[0]) {
			prev = lower
			continue
		}
		first, _ := utf8.DecodeRuneInString(word)
		capitalised := unicode.IsUpper(first) && !sentenceStart(text[:loc[0]])
		if !notNames[lower] && (capitalised || honorifics[prev]) && !honorifics[lower] {
			b.WriteString(text[last:loc[0]])
			b.WriteString(RedactedName)
			last = loc[1]
		}
		prev = lower
	}
	b.WriteString(text[last:])
	return b.String()
}

// sentenceStart reports whether a word preceded by before starts a sentence.
func sentenceStart(before string) bool {
	before = strings.TrimRight(before, " \t")
	if before == "" {
		return true
	}
	r, _ := utf8.DecodeLastRuneInString(before)
	return r == '.' || r == '!' || r == '?' || r == '\n'
}

// isPlaceholder reports whether the word at i is inside an earlier placeholder.
func isPlaceholder(text string, i int) bool { return i > 0 && text[i-1] == '<' }

func setOf(words ...string) map[string]bool {
	m := make(map[string]bool, len(words))
	for _, w := range words {
		m[w] = true
	}
	return m
}
