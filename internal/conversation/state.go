package conversation

import (
	"sync"
	"time"
)

// session holds the ongoing conversation context. Phase 1 keeps it in
// process memory; Phase 2 moves it to Redis (key conv:<user_id>).
type session struct {
	awaitingName bool
	pendingText  string // message waiting for an answer to "Berapa nominalnya?"
	expiresAt    time.Time
}

type sessionStore struct {
	mu   sync.Mutex
	ttl  time.Duration
	data map[string]session
}

func newSessionStore(ttl time.Duration) *sessionStore {
	return &sessionStore{ttl: ttl, data: map[string]session{}}
}

func (s *sessionStore) Get(userID string, now time.Time) session {
	s.mu.Lock()
	defer s.mu.Unlock()
	st, ok := s.data[userID]
	if !ok || now.After(st.expiresAt) {
		return session{}
	}
	return st
}

func (s *sessionStore) Update(userID string, now time.Time, fn func(*session)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	st := s.data[userID]
	if now.After(st.expiresAt) {
		st = session{}
	}
	fn(&st)
	st.expiresAt = now.Add(s.ttl)
	s.data[userID] = st
}

func (s *sessionStore) Delete(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, userID)
}

// rateLimiter is a per-user fixed window.
type rateLimiter struct {
	mu     sync.Mutex
	limit  int
	window time.Duration
	counts map[int64]windowCount
}

type windowCount struct {
	start time.Time
	n     int
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{limit: limit, window: window, counts: map[int64]windowCount{}}
}

func (r *rateLimiter) Allow(key int64, now time.Time) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	c := r.counts[key]
	if now.Sub(c.start) >= r.window {
		c = windowCount{start: now}
	}
	c.n++
	r.counts[key] = c
	return c.n <= r.limit
}
