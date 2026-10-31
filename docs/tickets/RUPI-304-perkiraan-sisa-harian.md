# RUPI-304 · Perkiraan sisa uang harian

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 3 · Event, budget, laporan | P1 | S | FR-32 |

## Cakupan
- Hitung sisa budget dibagi sisa hari periode, misalnya
  "sisa Rp340rb untuk 12 hari, aman ~Rp28rb/hari".
- Tampilkan di `/budget`, di `/ringkas`, dan di balasan peringatan budget.

## Kriteria penerimaan
- [ ] Unit test untuk akhir bulan, tahun kabisat, dan budget yang sudah terlampaui.
- [ ] Pembulatan ke ribuan terdekat ke bawah.

## Dependensi
RUPI-302.
