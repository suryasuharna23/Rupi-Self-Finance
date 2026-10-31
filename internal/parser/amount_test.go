package parser

import "testing"

func TestFindAmounts(t *testing.T) {
	cases := []struct {
		in        string
		want      int64
		ambiguous bool
	}{
		{"15000", 15_000, false},
		{"15.000", 15_000, false},
		{"15,000", 15_000, false},
		{"Rp15.000", 15_000, false},
		{"rp 15.000", 15_000, false},
		{"15rb", 15_000, false},
		{"15 rb", 15_000, false},
		{"15ribu", 15_000, false},
		{"15k", 15_000, false},
		{"1,5jt", 1_500_000, false},
		{"1.5jt", 1_500_000, false},
		{"1,5 juta", 1_500_000, false},
		{"1500k", 1_500_000, false},
		{"goceng", 5_000, false},
		{"ceban", 10_000, false},
		{"goban", 50_000, false},
		{"cepek", 100, true},
		{"gocap", 50, true},
		{"sejuta", 1_000_000, false},
		{"setengah juta", 500_000, false},
		{"print tugas 3500", 3_500, false},
		{"bakso 15", 15, true},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			got := FindAmounts(c.in)
			if len(got) != 1 {
				t.Fatalf("FindAmounts(%q) = %d amounts, want 1: %+v", c.in, len(got), got)
			}
			if got[0].Value != c.want || got[0].Ambiguous != c.ambiguous {
				t.Errorf("FindAmounts(%q) = %d (ambiguous=%v), want %d (ambiguous=%v)",
					c.in, got[0].Value, got[0].Ambiguous, c.want, c.ambiguous)
			}
		})
	}
}

func TestFindAmountsMultiple(t *testing.T) {
	got := FindAmounts("bakso 15rb sama es teh goceng")
	if len(got) != 2 || got[0].Value != 15_000 || got[1].Value != 5_000 {
		t.Fatalf("got %+v", got)
	}
}

func TestFindAmountsRejectsOutOfRange(t *testing.T) {
	if got := FindAmounts("2000jt"); len(got) != 0 {
		t.Fatalf("amounts above Rp1 billion must be rejected, got %+v", got)
	}
}
