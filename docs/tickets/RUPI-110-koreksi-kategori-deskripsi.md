# RUPI-110 · Koreksi kategori, deskripsi, dan tanggal

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 1 · Inti percakapan | P0 | M | FR-13, FR-14 |

## Latar
Koreksi saat ini hanya bisa mengubah nominal, dan tombol "Ubah" hanya memberi
petunjuk teks.

## Cakupan
- Dukung kalimat seperti "yang bensin tadi masuk transportasi" dan "yang bakso tadi kemarin".
- Tombol "Ubah" menampilkan pilihan Nominal / Kategori / Tanggal; kategori dipilih lewat tombol.
- Setiap perubahan tercatat di `transaction_edits` dengan nama field-nya.

## Kriteria penerimaan
- [ ] Koreksi kategori, deskripsi, dan tanggal bisa lewat kalimat maupun tombol.
- [ ] Ada jejak audit per field.
- [ ] Contoh kalimat koreksi baru masuk ke dataset evaluasi.

## Dependensi
RUPI-101.
