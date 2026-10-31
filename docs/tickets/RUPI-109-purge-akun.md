# RUPI-109 · Purge permanen akun setelah 30 hari

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 1 · Inti percakapan | P0 | S | FR-03, SEC-10 |

## Cakupan
- Job terjadwal (goroutine harian di Fase 1, CronJob di Fase 2) yang menghapus
  permanen baris `users` dengan `deleted_at` lebih dari 30 hari.
- Pastikan tidak ada sisa data pengguna di schema lain (transaksi, kategori kustom,
  parser_logs, budget).

## Kriteria penerimaan
- [ ] Akun yang dihapus 31 hari lalu hilang total; yang dihapus 29 hari lalu masih ada.
- [ ] Integration test yang memeriksa semua tabel milik pengguna.

## Dependensi
RUPI-102.
