# RUPI-302 · Budget Service dan pengaturan budget lewat chat

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 3 · Event, budget, laporan | P1 | M | FR-30, `BudgetService` |

## Cakupan
- `cmd/budget-service` dengan `SetBudget` dan `GetBudgetStatus`; schema `budget_svc`
  dengan tabel `budgets` dan `budget_alerts`.
- Parser mengenali kalimat budget: "budget makan 500rb sebulan", "budget total 1jt per minggu".
- Perintah `/budget` menampilkan pemakaian setiap budget.

## Kriteria penerimaan
- [ ] Budget per kategori maupun total, per minggu atau per bulan, bisa diatur lewat chat.
- [ ] `/budget` menampilkan pemakaian dan sisa yang benar.
- [ ] Contoh kalimat budget masuk ke dataset evaluasi.

## Dependensi
RUPI-301.
