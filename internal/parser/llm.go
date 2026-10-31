package parser

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"
)

// LLM adalah lapis kedua parser. Implementasi hanya menerima teks pesan dan
// daftar kategori (SEC-04), dan hanya mengembalikan data terstruktur: LLM
// tidak pernah menulis ke database atau menjalankan aksi (SEC-06).
//
// Penyedia LLM masih pertanyaan terbuka di dokumen kebutuhan, jadi belum ada
// implementasi konkret.
type LLM interface {
	Parse(ctx context.Context, text string, now time.Time, categories []string) (Result, error)
}

var ErrInvalidLLMOutput = errors.New("keluaran LLM tidak valid")

// Validate memeriksa keluaran LLM terhadap batas wajar (SEC-05).
func Validate(res Result, now time.Time) error {
	switch res.Kind {
	case KindTransactions:
		if len(res.Transactions) == 0 {
			return fmt.Errorf("%w: daftar transaksi kosong", ErrInvalidLLMOutput)
		}
		for _, tx := range res.Transactions {
			if tx.Amount < MinAmount || tx.Amount > MaxAmount {
				return fmt.Errorf("%w: nominal %d di luar batas", ErrInvalidLLMOutput, tx.Amount)
			}
			if tx.Type != TypeExpense && tx.Type != TypeIncome {
				return fmt.Errorf("%w: jenis transaksi %s", ErrInvalidLLMOutput, tx.Type)
			}
			if tx.OccurredAt.After(now.Add(time.Minute)) || tx.OccurredAt.Before(now.AddDate(-1, 0, 0)) {
				return fmt.Errorf("%w: tanggal %s tidak wajar", ErrInvalidLLMOutput, tx.OccurredAt)
			}
			valid := ExpenseCategories
			if tx.Type == TypeIncome {
				valid = IncomeCategories
			}
			if !slices.Contains(valid, tx.Category) {
				return fmt.Errorf("%w: kategori %q", ErrInvalidLLMOutput, tx.Category)
			}
		}
	case KindQuery:
		if res.Query == nil || !res.Query.Range.Start.Before(res.Query.Range.End) {
			return fmt.Errorf("%w: query tidak lengkap", ErrInvalidLLMOutput)
		}
	case KindClarification, KindSmallTalk:
	default:
		// Koreksi selalu lewat jalur aturan karena butuh ID transaksi.
		return fmt.Errorf("%w: jenis hasil %d", ErrInvalidLLMOutput, res.Kind)
	}
	return nil
}
