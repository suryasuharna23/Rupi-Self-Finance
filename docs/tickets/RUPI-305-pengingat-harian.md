# RUPI-305 · Pengingat harian opsional

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 3 · Event, budget, laporan | P1 | S | FR-33 |

## Cakupan
- Perintah `/pengingat on|off` dan jam pengingat (default 20.00 waktu pengguna).
- Pengingat hanya dikirim bila pengguna belum mencatat apa pun hari itu.

## Kriteria penerimaan
- [ ] Pengguna yang sudah mencatat hari ini tidak menerima pengingat.
- [ ] Jam pengingat mengikuti zona waktu pengguna.
- [ ] Default mati; pengguna harus mengaktifkannya sendiri.

## Dependensi
RUPI-306 (memakai mekanisme penjadwalan yang sama).
