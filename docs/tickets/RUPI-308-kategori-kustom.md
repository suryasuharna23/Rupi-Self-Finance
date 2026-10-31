# RUPI-308 · Kategori kustom dan aturan kata kunci

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 3 · Event, budget, laporan | P1 | M | FR-17 |

## Cakupan
- Kalimat seperti "bikin kategori Kopi" dan "kalau ada kata 'kopken' masukkan ke Jajan".
- Simpan di `transaction_svc.categories` (kolom `user_id` dan `keywords` sudah ada).
- Parser menerima daftar kategori kustom pengguna dan mendahulukan kata kuncinya
  dibanding kamus bawaan.
- Perintah `/kategori` untuk melihat dan menghapus kategori kustom.

## Kriteria penerimaan
- [ ] Setelah aturan dibuat, "kopken 25rb" masuk ke kategori kustom tersebut.
- [ ] Kategori kustom pengguna A tidak terlihat oleh pengguna B.
- [ ] Nama kategori bawaan tidak bisa ditimpa.

## Dependensi
RUPI-203, RUPI-204.
