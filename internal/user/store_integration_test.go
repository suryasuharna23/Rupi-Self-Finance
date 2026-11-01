//go:build integration

package user_test

import (
	"context"
	"os"
	"testing"

	"github.com/suryasuharna23/rupi/internal/storage/storagetest"
	"github.com/suryasuharna23/rupi/internal/user"
)

func TestMain(m *testing.M) { os.Exit(storagetest.Run(m)) }

func TestGetOrCreate(t *testing.T) {
	ctx := context.Background()
	store := user.NewStore(storagetest.DB(t))
	tgID := storagetest.TelegramID()

	first, created, err := store.GetOrCreate(ctx, tgID)
	if err != nil {
		t.Fatal(err)
	}
	if !created || first.TelegramID != tgID || first.Timezone != "Asia/Jakarta" || first.Onboarded {
		t.Fatalf("first GetOrCreate = %+v, created %v", first, created)
	}

	again, created, err := store.GetOrCreate(ctx, tgID)
	if err != nil {
		t.Fatal(err)
	}
	if created || again.ID != first.ID {
		t.Errorf("second GetOrCreate = %s (created %v), want existing %s", again.ID, created, first.ID)
	}

	if err := store.CompleteOnboarding(ctx, first.ID, "Dina"); err != nil {
		t.Fatal(err)
	}
	onboarded, _, err := store.GetOrCreate(ctx, tgID)
	if err != nil {
		t.Fatal(err)
	}
	if !onboarded.Onboarded || onboarded.DisplayName != "Dina" {
		t.Errorf("after onboarding = %+v", onboarded)
	}
}

func TestGetOrCreateAfterSoftDelete(t *testing.T) {
	ctx := context.Background()
	store := user.NewStore(storagetest.DB(t))
	tgID := storagetest.TelegramID()

	old, _, err := store.GetOrCreate(ctx, tgID)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.CompleteOnboarding(ctx, old.ID, "Dina"); err != nil {
		t.Fatal(err)
	}
	if err := store.Delete(ctx, old.ID); err != nil {
		t.Fatal(err)
	}

	fresh, created, err := store.GetOrCreate(ctx, tgID)
	if err != nil {
		t.Fatal(err)
	}
	if !created || fresh.ID == old.ID {
		t.Errorf("GetOrCreate after delete = %s (created %v), want a new user, not %s", fresh.ID, created, old.ID)
	}
	if fresh.Onboarded || fresh.DisplayName != "" {
		t.Errorf("new user inherited the old profile: %+v", fresh)
	}
}
