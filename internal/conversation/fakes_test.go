package conversation

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/suryasuharna23/rupi/internal/parser"
	"github.com/suryasuharna23/rupi/internal/report"
	"github.com/suryasuharna23/rupi/internal/transaction"
	"github.com/suryasuharna23/rupi/internal/user"
)

// In-memory fakes for the Service dependencies (RUPI-103). Each has an err
// field (or per-method error) that, when set, is returned instead of a result.

type fakeUsers struct {
	mu         sync.Mutex
	byTG       map[int64]*user.User
	nextID     int
	getErr     error
	onboardErr error
	deleteErr  error
	deleted    []string
}

func newFakeUsers() *fakeUsers { return &fakeUsers{byTG: map[int64]*user.User{}} }

// add registers an existing user, as if they had signed up before.
func (f *fakeUsers) add(tgID int64, name string, onboarded bool) *user.User {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.nextID++
	u := &user.User{ID: fmt.Sprintf("user-%d", f.nextID), TelegramID: tgID, DisplayName: name, Timezone: "Asia/Jakarta", Onboarded: onboarded}
	f.byTG[tgID] = u
	return u
}

func (f *fakeUsers) GetOrCreate(_ context.Context, tgID int64) (user.User, bool, error) {
	if f.getErr != nil {
		return user.User{}, false, f.getErr
	}
	f.mu.Lock()
	u, ok := f.byTG[tgID]
	f.mu.Unlock()
	if ok {
		return *u, false, nil
	}
	return *f.add(tgID, "", false), true, nil
}

func (f *fakeUsers) CompleteOnboarding(_ context.Context, userID, name string) error {
	if f.onboardErr != nil {
		return f.onboardErr
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, u := range f.byTG {
		if u.ID == userID {
			u.DisplayName, u.Onboarded = name, true
		}
	}
	return nil
}

func (f *fakeUsers) Delete(_ context.Context, userID string) error {
	if f.deleteErr != nil {
		return f.deleteErr
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	for tg, u := range f.byTG {
		if u.ID == userID {
			delete(f.byTG, tg)
		}
	}
	f.deleted = append(f.deleted, userID)
	return nil
}

type storedTx struct {
	transaction.Transaction
	userID, messageID string
	deleted           bool
}

type fakeTxs struct {
	mu                                                                         sync.Mutex
	txs                                                                        []*storedTx
	createErr, updateErr, deleteMsgErr, deleteLastErr, recentErr, deleteAllErr error
	deletedAllFor                                                              []string
}

func (f *fakeTxs) Create(_ context.Context, userID, msgID string, _ parser.Source, in []parser.ParsedTransaction) ([]transaction.Transaction, error) {
	if f.createErr != nil {
		return nil, f.createErr
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]transaction.Transaction, 0, len(in))
	for _, p := range in {
		t := transaction.Transaction{ID: fmt.Sprintf("tx-%d", len(f.txs)+1), Type: p.Type, Amount: p.Amount, Description: p.Description, Category: p.Category, OccurredAt: p.OccurredAt}
		f.txs = append(f.txs, &storedTx{Transaction: t, userID: userID, messageID: msgID})
		out = append(out, t)
	}
	return out, nil
}

// live returns the user's transactions that have not been deleted, newest last.
func (f *fakeTxs) live(userID string) []*storedTx {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*storedTx
	for _, t := range f.txs {
		if t.userID == userID && !t.deleted {
			out = append(out, t)
		}
	}
	return out
}

func (f *fakeTxs) UpdateAmount(_ context.Context, userID, txID string, amount int64) (transaction.Transaction, error) {
	if f.updateErr != nil {
		return transaction.Transaction{}, f.updateErr
	}
	for _, t := range f.live(userID) {
		if t.ID == txID {
			t.Amount = amount
			return t.Transaction, nil
		}
	}
	return transaction.Transaction{}, transaction.ErrNotFound
}

func (f *fakeTxs) DeleteByMessage(_ context.Context, userID, msgID string) (int64, error) {
	if f.deleteMsgErr != nil {
		return 0, f.deleteMsgErr
	}
	var n int64
	for _, t := range f.live(userID) {
		if t.messageID == msgID {
			t.deleted = true
			n++
		}
	}
	return n, nil
}

func (f *fakeTxs) DeleteLast(_ context.Context, userID string) (transaction.Transaction, error) {
	if f.deleteLastErr != nil {
		return transaction.Transaction{}, f.deleteLastErr
	}
	live := f.live(userID)
	if len(live) == 0 {
		return transaction.Transaction{}, transaction.ErrNotFound
	}
	last := live[len(live)-1]
	last.deleted = true
	return last.Transaction, nil
}

func (f *fakeTxs) FindRecent(_ context.Context, userID string, limit int) ([]parser.RecentTransaction, error) {
	if f.recentErr != nil {
		return nil, f.recentErr
	}
	var out []parser.RecentTransaction
	live := f.live(userID)
	for i := len(live) - 1; i >= 0 && len(out) < limit; i-- {
		t := live[i]
		out = append(out, parser.RecentTransaction{ID: t.ID, Description: t.Description, Amount: t.Amount, Category: t.Category, OccurredAt: t.OccurredAt})
	}
	return out, nil
}

func (f *fakeTxs) DeleteAllForUser(_ context.Context, userID string) error {
	if f.deleteAllErr != nil {
		return f.deleteAllErr
	}
	for _, t := range f.live(userID) {
		t.deleted = true
	}
	f.deletedAllFor = append(f.deletedAllFor, userID)
	return nil
}

type fakeReports struct {
	sum        int64
	periods    []report.Period
	err        error
	lastQuery  parser.ReportQuery
	summaryNow time.Time
}

func (f *fakeReports) Sum(_ context.Context, _ string, q parser.ReportQuery) (int64, error) {
	f.lastQuery = q
	return f.sum, f.err
}

func (f *fakeReports) Summary(_ context.Context, _ string, now time.Time) ([]report.Period, error) {
	f.summaryNow = now
	return f.periods, f.err
}

type fakeIdempotency struct {
	mu        sync.Mutex
	seen      map[string]bool
	markErr   error
	unmarkErr error
	unmarked  []string
}

func newFakeIdempotency() *fakeIdempotency { return &fakeIdempotency{seen: map[string]bool{}} }

func (f *fakeIdempotency) MarkProcessed(_ context.Context, channel, msgID string) (bool, error) {
	if f.markErr != nil {
		return false, f.markErr
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	key := channel + "/" + msgID
	if f.seen[key] {
		return false, nil
	}
	f.seen[key] = true
	return true, nil
}

func (f *fakeIdempotency) Unmark(_ context.Context, channel, msgID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.unmarked = append(f.unmarked, msgID)
	if f.unmarkErr != nil {
		return f.unmarkErr
	}
	delete(f.seen, channel+"/"+msgID)
	return nil
}

// fakeParser returns whatever result is set and records every request.
type fakeParser struct {
	result   parser.Result
	requests []parser.Request
}

func (f *fakeParser) Parse(_ context.Context, req parser.Request) parser.Result {
	f.requests = append(f.requests, req)
	return f.result
}

func (f *fakeParser) lastText() string {
	if len(f.requests) == 0 {
		return ""
	}
	return f.requests[len(f.requests)-1].Text
}
