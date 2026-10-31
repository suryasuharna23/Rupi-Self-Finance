package conversation

import (
	"strconv"
	"strings"
)

// Rupiah memformat nominal seperti "Rp15.000".
func Rupiah(amount int64) string {
	sign := ""
	if amount < 0 {
		sign, amount = "-", -amount
	}
	s := strconv.FormatInt(amount, 10)
	var b strings.Builder
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteRune(r)
	}
	return sign + "Rp" + b.String()
}
