//go:build integration

package transaction_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/suryasuharna23/rupi/internal/parser"
	"github.com/suryasuharna23/rupi/internal/storage/storagetest"
	"github.com/suryasuharna23/rupi/internal/transaction"
)

func TestMain(m *testing.M) { os.Exit(storagetest.Run(m)) }

func expense(amount int64, desc, category string) parser.ParsedTransaction {
	return parser.ParsedTransaction{Type: parser.TypeExpense, Amount: amount, Description: desc, Category: category, OccurredAt: time.Now()}
}

// countByMessage counts live transactions stored for one source message.
func countByMessage(t *testing.T, userID, msgID string) int {
	t.Helper()
	var n int
	err := storagetest.DB(t).QueryRow(context.Background(), `
		SELECT count(*) FROM transaction_svc.transactions
		WHERE user_id = $1 AND source_message_id = $2 AND deleted_at IS NULL`, userID, msgID).Scan(&n)
	if err != nil {
		t.Fatal(err)
	}
	return n
}

func TestCreate(t *testing.T) {
	ctx := context.Background()
	db := storagetest.DB(t)
	store := transaction.NewStore(db)

	t.Run("maps category names to IDs by type", func(t *testing.T) {
		user := storagetest.NewUser(t)
		income := expense(500000, "dari ortu", "Lainnya")
		income.Type = parser.TypeIncome
		txs, err := store.Create(ctx, user, "m1", parser.SourceRule, []parser.ParsedTransaction{
			expense(15000, "bakso", "Makanan & Minuman"),
			expense(2000, "parkir", "Lainnya"),
			income,
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(txs) != 3 {
			t.Fatalf("got %d transactions, want 3", len(txs))
		}
		// "Lainnya" exists as both an expense and an income category, so the
		// stored category must match the transaction type too.
		for _, tx := range txs {
			var name string
			var typ int
			err := db.QueryRow(ctx, `
				SELECT c.name, c.type FROM transaction_svc.transactions t
				JOIN transaction_svc.categories c ON c.id = t.category_id
				WHERE t.id = $1`, tx.ID).Scan(&name, &typ)
			if err != nil {
				t.Fatal(err)
			}
			if name != tx.Category || typ != int(tx.Type) {
				t.Errorf("%s: stored category (%s, %d), want (%s, %d)", tx.Description, name, typ, tx.Category, tx.Type)
			}
		}
	})

	t.Run("prefers the user's own category over the default", func(t *testing.T) {
		user := storagetest.NewUser(t)
		var own string
		err := db.QueryRow(ctx, `
			INSERT INTO transaction_svc.categories (user_id, name, type) VALUES ($1, 'Belanja', 1)
			RETURNING id`, user).Scan(&own)
		if err != nil {
			t.Fatal(err)
		}
		txs, err := store.Create(ctx, user, "m1", parser.SourceRule, []parser.ParsedTransaction{expense(30000, "sabun", "Belanja")})
		if err != nil {
			t.Fatal(err)
		}
		var got string
		if err := db.QueryRow(ctx, `SELECT category_id FROM transaction_svc.transactions WHERE id = $1`, txs[0].ID).Scan(&got); err != nil {
			t.Fatal(err)
		}
		if got != own {
			t.Errorf("category_id = %s, want the user's own %s", got, own)
		}
	})

	t.Run("is atomic", func(t *testing.T) {
		user := storagetest.NewUser(t)
		_, err := store.Create(ctx, user, "m1", parser.SourceRule, []parser.ParsedTransaction{
			expense(15000, "bakso", "Makanan & Minuman"),
			expense(5000, "???", "Kategori Tidak Ada"),
		})
		if err == nil {
			t.Fatal("want an error for an unknown category")
		}
		if n := countByMessage(t, user, "m1"); n != 0 {
			t.Errorf("%d transactions saved after a failed Create, want 0", n)
		}
	})
}

func TestUpdateAmountWritesAuditTrail(t *testing.T) {
	ctx := context.Background()
	db := storagetest.DB(t)
	store := transaction.NewStore(db)
	user := storagetest.NewUser(t)

	txs, err := store.Create(ctx, user, "m1", parser.SourceRule, []parser.ParsedTransaction{expense(15000, "bakso", "Makanan & Minuman")})
	if err != nil {
		t.Fatal(err)
	}
	got, err := store.UpdateAmount(ctx, user, txs[0].ID, 20000)
	if err != nil {
		t.Fatal(err)
	}
	if got.Amount != 20000 || got.Category != "Makanan & Minuman" || got.Description != "bakso" {
		t.Errorf("UpdateAmount returned %+v", got)
	}

	var field, oldValue, newValue string
	err = db.QueryRow(ctx, `
		SELECT field, old_value, new_value FROM transaction_svc.transaction_edits
		WHERE transaction_id = $1`, txs[0].ID).Scan(&field, &oldValue, &newValue)
	if err != nil {
		t.Fatal(err)
	}
	if field != "amount" || oldValue != "15000" || newValue != "20000" {
		t.Errorf("audit row = (%s, %s, %s), want (amount, 15000, 20000)", field, oldValue, newValue)
	}
}

func TestDeleteLastRespectsUndoWindow(t *testing.T) {
	ctx := context.Background()
	db := storagetest.DB(t)
	store := transaction.NewStore(db)
	user := storagetest.NewUser(t)

	txs, err := store.Create(ctx, user, "old", parser.SourceRule, []parser.ParsedTransaction{expense(15000, "bakso", "Makanan & Minuman")})
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(ctx, `UPDATE transaction_svc.transactions SET created_at = now() - interval '25 hours' WHERE id = $1`, txs[0].ID)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := store.DeleteLast(ctx, user); !errors.Is(err, transaction.ErrNotFound) {
		t.Fatalf("DeleteLast on a 25-hour-old transaction: err = %v, want ErrNotFound", err)
	}
	if n := countByMessage(t, user, "old"); n != 1 {
		t.Errorf("old transaction was deleted")
	}

	recent, err := store.Create(ctx, user, "new", parser.SourceRule, []parser.ParsedTransaction{expense(8000, "es teh", "Makanan & Minuman")})
	if err != nil {
		t.Fatal(err)
	}
	got, err := store.DeleteLast(ctx, user)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != recent[0].ID {
		t.Errorf("DeleteLast deleted %s, want %s", got.ID, recent[0].ID)
	}
}

// TestUserIsolation checks SEC-02: user B cannot touch user A's data.
func TestUserIsolation(t *testing.T) {
	ctx := context.Background()
	store := transaction.NewStore(storagetest.DB(t))
	alice, bob := storagetest.NewUser(t), storagetest.NewUser(t)

	txs, err := store.Create(ctx, alice, "m1", parser.SourceRule, []parser.ParsedTransaction{expense(15000, "bakso", "Makanan & Minuman")})
	if err != nil {
		t.Fatal(err)
	}

	if recent, err := store.FindRecent(ctx, bob, 10); err != nil || len(recent) != 0 {
		t.Errorf("FindRecent as bob = %v, %v; want nothing", recent, err)
	}
	if _, err := store.UpdateAmount(ctx, bob, txs[0].ID, 1); !errors.Is(err, transaction.ErrNotFound) {
		t.Errorf("UpdateAmount as bob: err = %v, want ErrNotFound", err)
	}
	if n, err := store.DeleteByMessage(ctx, bob, "m1"); err != nil || n != 0 {
		t.Errorf("DeleteByMessage as bob deleted %d (err %v), want 0", n, err)
	}
	if _, err := store.DeleteLast(ctx, bob); !errors.Is(err, transaction.ErrNotFound) {
		t.Errorf("DeleteLast as bob: err = %v, want ErrNotFound", err)
	}
	if err := store.DeleteAllForUser(ctx, bob); err != nil {
		t.Fatal(err)
	}

	recent, err := store.FindRecent(ctx, alice, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(recent) != 1 || recent[0].Amount != 15000 {
		t.Errorf("alice's transactions after bob's attempts = %+v, want the original one", recent)
	}
}
