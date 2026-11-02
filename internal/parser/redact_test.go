package parser

import "testing"

func TestRedact(t *testing.T) {
	cases := []struct{ in, want string }{
		// Phone numbers.
		{"pulsa ke 081234567890 25rb", "pulsa ke <phone> 25rb"},
		{"isi pulsa 0812-3456-7890 50rb", "isi pulsa <phone> 50rb"},
		{"kirim ke +62 812 3456 7890", "kirim ke <phone>"},
		{"nomornya 6281234567890", "nomornya <phone>"},
		// Bank account numbers.
		{"transfer ke rek 1234567890 500rb", "transfer ke rek <account> 500rb"},
		{"tf ke 123-456-789-012 100rb", "tf ke <account> 100rb"},
		{"va 8808123456789012 bayar ukt", "va <account> bayar ukt"},
		// Emails.
		{"bayar ke budi.santoso+kos@gmail.com 1jt", "bayar ke <email> 1jt"},
		// Names.
		{"transfer ke Budi 50rb", "transfer ke <name> 50rb"},
		{"utang ke Siti Rahma 20rb", "utang ke <name> <name> 20rb"},
		{"bayar mas budi 15rb", "bayar mas <name> 15rb"},
		{"Kak Dina traktir 30rb", "Kak <name> traktir 30rb"},
		{"tf 200rb a.n. Andi Wijaya", "tf 200rb a.n. <name>"},
		{"transfer atas nama Rina 75rb", "transfer atas nama <name> 75rb"},
		{"tf an budi santoso 50rb", "tf an <name> 50rb"},
		{"a.n. Rina buat kado 100rb", "a.n. <name> buat kado 100rb"},
		// Not personal data: amounts, keywords, brands, sentence starts.
		{"beli bakso 15rb", "beli bakso 15rb"},
		{"Makan siang Rp25.000", "Makan siang Rp25.000"},
		{"dapet transferan ortu 1,5jt", "dapet transferan ortu 1,5jt"},
		{"gaji 1.500.000 sama bonus 150000", "gaji 1.500.000 sama bonus 150000"},
		{"beli Indomie 3500. Terus naik Gojek 12rb", "beli Indomie 3500. Terus naik Gojek 12rb"},
		{"bayar bu kos 800rb", "bayar bu kos 800rb"},
		{"dikasih ibu uang saku 500rb", "dikasih ibu uang saku 500rb"},
		{"bulan Oktober habis berapa?", "bulan Oktober habis berapa?"},
		{"bakso 15rb sama es teh goceng", "bakso 15rb sama es teh goceng"},
	}
	for _, c := range cases {
		if got := Redact(c.in); got != c.want {
			t.Errorf("Redact(%q)\n got  %q\n want %q", c.in, got, c.want)
		}
	}
}
