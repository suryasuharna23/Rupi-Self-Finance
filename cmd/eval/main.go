// Command eval mengukur akurasi parser terhadap dataset berlabel (NFR-03)
// dan porsi pesan yang selesai lewat jalur aturan (NFR-04).
//
//	go run ./cmd/eval -dataset testdata/eval/dataset.jsonl
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/suryasuharna23/rupi/internal/parser"
)

type expectedTx struct {
	Type     string `json:"type"`
	Amount   int64  `json:"amount"`
	Category string `json:"category"`
}

type example struct {
	Text         string       `json:"text"`
	Kind         string       `json:"kind"`
	Transactions []expectedTx `json:"transactions"`
}

var kindNames = map[parser.Kind]string{
	parser.KindTransactions:  "transactions",
	parser.KindQuery:         "query",
	parser.KindCorrection:    "correction",
	parser.KindClarification: "clarification",
	parser.KindSmallTalk:     "small_talk",
}

func main() {
	path := flag.String("dataset", "testdata/eval/dataset.jsonl", "file JSONL berlabel")
	minAccuracy := flag.Float64("min-accuracy", 0.95, "batas akurasi ekstraksi; di bawahnya exit 1")
	verbose := flag.Bool("v", false, "tampilkan setiap kesalahan")
	flag.Parse()

	examples, err := load(*path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	now := time.Now().In(time.FixedZone("WIB", 7*3600))
	var correct, txTotal, txCorrect, catCorrect, ruleHandled int
	for _, ex := range examples {
		res, ok := parser.ParseRules(parser.Request{Text: ex.Text, Location: now.Location(), ReceivedAt: now})
		if ok {
			ruleHandled++
		}
		if kindNames[res.Kind] != ex.Kind {
			report(*verbose, ex.Text, "jenis %q, harusnya %q", kindNames[res.Kind], ex.Kind)
			txTotal += len(ex.Transactions)
			continue
		}
		allOK := len(res.Transactions) == len(ex.Transactions)
		for i, want := range ex.Transactions {
			txTotal++
			if i >= len(res.Transactions) {
				continue
			}
			got := res.Transactions[i]
			if got.Amount == want.Amount && got.Type.String() == want.Type {
				txCorrect++
			} else {
				allOK = false
				report(*verbose, ex.Text, "tx[%d] %s %d, harusnya %s %d", i, got.Type, got.Amount, want.Type, want.Amount)
			}
			if got.Category == want.Category {
				catCorrect++
			} else {
				report(*verbose, ex.Text, "tx[%d] kategori %q, harusnya %q", i, got.Category, want.Category)
			}
		}
		if allOK {
			correct++
		}
	}

	n := len(examples)
	extraction := ratio(txCorrect, txTotal)
	fmt.Printf("contoh               : %d (target dataset MVP: 300)\n", n)
	fmt.Printf("pesan benar          : %d/%d (%.1f%%)\n", correct, n, 100*ratio(correct, n))
	fmt.Printf("ekstraksi transaksi  : %d/%d (%.1f%%, target >= 95%%)\n", txCorrect, txTotal, 100*extraction)
	fmt.Printf("kategori             : %d/%d (%.1f%%, target >= 85%%)\n", catCorrect, txTotal, 100*ratio(catCorrect, txTotal))
	fmt.Printf("selesai lewat aturan : %d/%d (%.1f%%, target >= 60%%)\n", ruleHandled, n, 100*ratio(ruleHandled, n))
	if extraction < *minAccuracy {
		os.Exit(1)
	}
}

func load(path string) ([]example, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var out []example
	sc := bufio.NewScanner(f)
	for line := 1; sc.Scan(); line++ {
		if len(sc.Bytes()) == 0 {
			continue
		}
		var ex example
		if err := json.Unmarshal(sc.Bytes(), &ex); err != nil {
			return nil, fmt.Errorf("%s:%d: %w", path, line, err)
		}
		out = append(out, ex)
	}
	return out, sc.Err()
}

func report(verbose bool, text, format string, args ...any) {
	if verbose {
		fmt.Printf("  ✗ %q: %s\n", text, fmt.Sprintf(format, args...))
	}
}

func ratio(a, b int) float64 {
	if b == 0 {
		return 1
	}
	return float64(a) / float64(b)
}
