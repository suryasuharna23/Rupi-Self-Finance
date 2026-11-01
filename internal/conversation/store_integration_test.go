//go:build integration

package conversation_test

import (
	"context"
	"os"
	"testing"

	"github.com/suryasuharna23/rupi/internal/conversation"
	"github.com/suryasuharna23/rupi/internal/storage/storagetest"
)

func TestMain(m *testing.M) { os.Exit(storagetest.Run(m)) }

// TestMarkProcessed checks NFR-07: a redelivered message is processed once.
func TestMarkProcessed(t *testing.T) {
	ctx := context.Background()
	store := conversation.NewIdempotencyStore(storagetest.DB(t))

	steps := []struct {
		name    string
		channel string
		msgID   string
		want    bool
	}{
		{"first delivery", "telegram", "1001", true},
		{"redelivery", "telegram", "1001", false},
		{"same ID on another channel", "whatsapp", "1001", true},
		{"another message", "telegram", "1002", true},
	}
	for _, s := range steps {
		got, err := store.MarkProcessed(ctx, s.channel, s.msgID)
		if err != nil {
			t.Fatal(err)
		}
		if got != s.want {
			t.Errorf("%s: MarkProcessed = %v, want %v", s.name, got, s.want)
		}
	}

	if err := store.Unmark(ctx, "telegram", "1001"); err != nil {
		t.Fatal(err)
	}
	if got, err := store.MarkProcessed(ctx, "telegram", "1001"); err != nil || !got {
		t.Errorf("MarkProcessed after Unmark = %v, %v; want true", got, err)
	}
}
