# RUPI-309 · Query lanjutan: jumlah, terbesar, dan rincian

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 3 · Event, budget, laporan | P0 | M | FR-20, `ReportQuery.Metric` |

## Latar
Saat ini hanya metrik SUM yang didukung. `ReportQuery` di dokumen juga
mendefinisikan COUNT, TOP_N, BREAKDOWN, dan GroupBy.

## Cakupan
- Parser mengenali: "berapa kali jajan kopi minggu ini?" (COUNT),
  "pengeluaran terbesar bulan ini apa?" (TOP_N),
  "rincian pengeluaran bulan ini per kategori" (BREAKDOWN + GroupBy).
- Report Service mengimplementasikan setiap metrik.
- Rentang waktu tambahan: "kemarin", "minggu lalu", "dari tgl 1 sampai 15".

## Kriteria penerimaan
- [ ] Setiap metrik punya unit test dan contoh di dataset evaluasi.
- [ ] Balasan BREAKDOWN terurut dari kategori terbesar.

## Dependensi
RUPI-206.
