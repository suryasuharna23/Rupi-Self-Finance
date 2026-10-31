# RUPI-107 · Tabel parser_logs dengan teks disamarkan

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 1 · Inti percakapan | P0 | M | SEC-08, model data `parser_logs` |

## Cakupan
- Migrasi schema `parser_svc` dengan tabel `parser_logs` (user_id hash, text_redacted,
  result JSON, source, latency_ms, tokens, created_at).
- Fungsi penyamaran: angka panjang (rekening, HP), email, dan nama orang.
- Hapus baris berumur > 30 hari secara terjadwal.
- Pastikan log aplikasi level INFO tidak memuat isi pesan mentah.

## Kriteria penerimaan
- [ ] Setiap pemanggilan parser menulis satu baris log.
- [ ] Unit test penyamaran untuk nomor HP, nomor rekening, dan email.
- [ ] Tidak ada log level INFO yang memuat teks pesan.

## Dependensi
Tidak ada.
