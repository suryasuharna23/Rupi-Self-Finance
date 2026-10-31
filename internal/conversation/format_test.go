package conversation

import (
	"testing"
	"time"
)

func TestRupiah(t *testing.T) {
	cases := map[int64]string{0: "Rp0", 500: "Rp500", 15_000: "Rp15.000", 1_500_000: "Rp1.500.000", -3_500: "-Rp3.500"}
	for in, want := range cases {
		if got := Rupiah(in); got != want {
			t.Errorf("Rupiah(%d) = %q, want %q", in, got, want)
		}
	}
}

func TestRateLimiter(t *testing.T) {
	r := newRateLimiter(2, time.Minute)
	now := time.Now()
	if !r.Allow(1, now) || !r.Allow(1, now) || r.Allow(1, now) {
		t.Fatal("pesan ketiga dalam satu jendela harus ditolak")
	}
	if !r.Allow(2, now) {
		t.Fatal("pengguna lain tidak boleh terpengaruh")
	}
	if !r.Allow(1, now.Add(time.Minute)) {
		t.Fatal("jendela baru harus mengizinkan lagi")
	}
}

func TestSessionExpires(t *testing.T) {
	s := newSessionStore(30 * time.Minute)
	now := time.Now()
	s.Update("u1", now, func(st *session) { st.pendingText = "makan siang" })
	if got := s.Get("u1", now.Add(29*time.Minute)); got.pendingText != "makan siang" {
		t.Fatalf("got %+v", got)
	}
	if got := s.Get("u1", now.Add(31*time.Minute)); got.pendingText != "" {
		t.Fatalf("sesi harus kedaluwarsa, got %+v", got)
	}
}
