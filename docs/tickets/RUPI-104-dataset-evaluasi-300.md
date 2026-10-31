# RUPI-104 · Dataset evaluasi parser 300 kalimat

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 1 · Inti percakapan | P0 | L | NFR-03, NFR-04, gerbang Fase 1 |

## Latar
Gerbang Fase 1 adalah akurasi ekstraksi ≥ 95% dan kategori ≥ 85% pada dataset
berlabel minimal 300 kalimat. `testdata/eval/dataset.jsonl` baru berisi 20.

## Cakupan
- Kumpulkan kalimat nyata dari calon pengguna uji (pertanyaan terbuka di dokumen).
- Sebar contoh merata: semua format nominal, waktu relatif, pemasukan, multi-transaksi,
  koreksi, query, small talk, typo, dan kalimat yang sengaja ambigu.
- Tambahkan label tanggal di dataset dan periksa tanggal di `cmd/eval`.
- Tambahkan skor per jenis kesalahan (nominal, jenis, tanggal, kategori).
- Perbaiki aturan parser sampai target tercapai.

## Kriteria penerimaan
- [ ] Dataset ≥ 300 kalimat; tiap kategori bawaan punya ≥ 10 contoh.
- [ ] `go run ./cmd/eval` melaporkan ekstraksi ≥ 95%, kategori ≥ 85%, jalur aturan ≥ 60%.
- [ ] Laporan memuat rincian per jenis kesalahan.
- [ ] Dataset tidak memuat data pribadi (nama asli, nomor rekening).

## Dependensi
Tidak ada. Bisa dikerjakan paralel dengan tiket lain.
