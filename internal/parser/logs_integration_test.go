//go:build integration

package parser_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/suryasuharna23/rupi/internal/parser"
	"github.com/suryasuharna23/rupi/internal/storage/storagetest"
)

func TestMain(m *testing.M) { os.Exit(storagetest.Run(m)) }

func TestLogStore(t *testing.T) {
	ctx := context.Background()
	db := storagetest.DB(t)
	store := parser.NewLogStore(db)
	p := parser.New(nil).WithLogs(store)
	wib := time.FixedZone("WIB", 7*3600)

	p.Parse(ctx, parser.Request{UserID: "user-1", Text: "transfer ke Budi 50rb", Location: wib, ReceivedAt: time.Now()})
	p.Parse(ctx, parser.Request{UserID: "user-1", Text: "apa kabar dunia", Location: wib, ReceivedAt: time.Now()})

	rows, err := db.Query(ctx, `
		SELECT user_id_hash, text_redacted, result, source, latency_ms, tokens
		FROM parser_svc.parser_logs ORDER BY created_at, source DESC`)
	if err != nil {
		t.Fatal(err)
	}
	type row struct {
		hash, text, source string
		result             []byte
		latency, tokens    int
	}
	var got []row
	for rows.Next() {
		var r row
		if err := rows.Scan(&r.hash, &r.text, &r.result, &r.source, &r.latency, &r.tokens); err != nil {
			t.Fatal(err)
		}
		got = append(got, r)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("%d rows, want one per Parse call (2)", len(got))
	}

	byText := map[string]row{}
	for _, r := range got {
		if r.hash != parser.HashUserID("user-1") {
			t.Errorf("user_id_hash = %q", r.hash)
		}
		byText[r.text] = r
	}
	tx, ok := byText["transfer ke <name> 50rb"]
	if !ok {
		t.Fatalf("redacted text not stored; rows: %+v", got)
	}
	if tx.source != parser.LogSourceRule {
		t.Errorf("source = %s, want RULE", tx.source)
	}
	var res parser.Result
	if err := json.Unmarshal(tx.result, &res); err != nil {
		t.Fatal(err)
	}
	if res.Kind != parser.KindTransactions || len(res.Transactions) != 1 || res.Transactions[0].Amount != 50_000 {
		t.Errorf("stored result = %+v", res)
	}
	if r := byText["apa kabar dunia"]; r.source != parser.LogSourceNoLLM {
		t.Errorf("fallback source = %q, want NO_LLM", r.source)
	}
}

func TestLogStoreRetention(t *testing.T) {
	ctx := context.Background()
	db := storagetest.DB(t)
	store := parser.NewLogStore(db)

	for _, age := range []string{"31 days", "29 days"} {
		_, err := db.Exec(ctx, `
			INSERT INTO parser_svc.parser_logs (user_id_hash, text_redacted, result, source, latency_ms, created_at)
			VALUES ('retention', $1, '{}', 'RULE', 1, now() - $2::interval)`, age, age)
		if err != nil {
			t.Fatal(err)
		}
	}

	if _, err := store.DeleteOlderThan(ctx, time.Now().Add(-parser.LogRetention)); err != nil {
		t.Fatal(err)
	}
	var left []string
	rows, err := db.Query(ctx, `SELECT text_redacted FROM parser_svc.parser_logs WHERE user_id_hash = 'retention'`)
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var s string
		rows.Scan(&s)
		left = append(left, s)
	}
	if len(left) != 1 || left[0] != "29 days" {
		t.Errorf("rows left = %v, want only the 29-day-old one", left)
	}
}
