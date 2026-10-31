# RUPI-505 · Ekspor transaksi ke CSV

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 5 · Fitur lanjutan | P2 | S | FR-24 |

## Cakupan
- Perintah `/ekspor` dengan rentang tanggal ("/ekspor bulan ini", "/ekspor 1-15 okt").
- Kirim file CSV lewat Telegram (`sendDocument`).
- Kolom: tanggal, jenis, kategori, deskripsi, nominal.

## Kriteria penerimaan
- [ ] File bisa dibuka benar di Excel dan Google Sheets (UTF-8 dengan BOM).
- [ ] Hanya berisi transaksi milik pengguna yang meminta.

## Dependensi
RUPI-206.
