# RUPI-306 · Laporan mingguan dan bulanan otomatis

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 3 · Event, budget, laporan | P1 | M | FR-22, Infrastruktur (CronJob) |

## Cakupan
- Kubernetes CronJob: Senin 08.00 (laporan mingguan) dan tanggal 1 (laporan bulanan), WIB.
- Job memanggil Report Service per pengguna aktif lalu mengirim lewat Conversation Service.
- Laporan berisi total keluar/masuk, 3 kategori terbesar, dan transaksi terbesar.

## Kriteria penerimaan
- [ ] Job yang dijalankan ulang tidak mengirim laporan ganda.
- [ ] Pengguna tanpa transaksi pada periode itu tidak menerima laporan kosong.
- [ ] Job selesai dalam batas resource 200m / 256Mi untuk 1.000 pengguna simulasi.

## Dependensi
RUPI-206, RUPI-207.
