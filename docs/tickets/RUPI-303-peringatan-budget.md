# RUPI-303 · Peringatan budget 80% dan 100%

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 3 · Event, budget, laporan | P1 | M | FR-31, gerbang Fase 3 |

## Cakupan
- Budget Service mengonsumsi event transaksi dan menghitung ulang pemakaian.
- Peringatan dikirim lewat `SubscribeOutbound` saat melewati 80% dan 100%.
- Tabel `budget_alerts` (budget_id, period_start, level) dengan unique constraint
  mencegah peringatan ganda.
- Transaksi yang dihapus atau dikoreksi ikut menurunkan pemakaian.

## Kriteria penerimaan
- [ ] Peringatan 80% terkirim tepat satu kali per periode, walau event diproses ulang.
- [ ] Koreksi yang menurunkan pemakaian di bawah ambang tidak memicu peringatan baru.
- [ ] Periode baru mereset status peringatan.

## Dependensi
RUPI-302.
