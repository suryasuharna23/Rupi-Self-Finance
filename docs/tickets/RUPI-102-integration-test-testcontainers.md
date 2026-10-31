# RUPI-102 · Integration test store dengan Testcontainers

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 1 · Inti percakapan | P0 | M | NFR-07, NFR-08, SEC-02, bagian Pengujian |

## Latar
Semua query SQL belum punya tes otomatis. Dokumen menetapkan integration test
memakai Postgres sungguhan lewat Testcontainers.

## Cakupan
- Helper test yang menyalakan `postgres:16-alpine` dan menjalankan `storage.Migrate`.
- Tes untuk `transaction.Store`, `user.Store`, `report.Store`, `conversation.IdempotencyStore`.
- Pisahkan dengan build tag `integration` agar `go test ./...` tetap cepat.
- Job CI terpisah untuk integration test.

## Kriteria penerimaan
- [ ] `Create` menyimpan beberapa transaksi secara atomik dan memetakan nama kategori ke ID.
- [ ] `UpdateAmount` menulis jejak di `transaction_edits`.
- [ ] `DeleteLast` tidak menghapus transaksi yang dibuat lebih dari 24 jam lalu.
- [ ] Pengguna A tidak bisa membaca, mengubah, atau menghapus transaksi pengguna B (SEC-02).
- [ ] `MarkProcessed` mengembalikan `false` untuk pesan yang sama kedua kali (NFR-07).
- [ ] `GetOrCreate` membuat pengguna baru setelah akun lama di-soft delete.
- [ ] `Summary` benar untuk minggu yang melewati awal bulan.

## Dependensi
RUPI-101.
