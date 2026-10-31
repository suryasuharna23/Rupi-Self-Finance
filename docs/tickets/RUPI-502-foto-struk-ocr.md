# RUPI-502 · Foto struk → OCR → transaksi per item

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 5 · Fitur lanjutan | P2 | L | FR-18 |

## Cakupan
- Channel Adapter menerima foto dari Telegram.
- OCR (layanan eksternal atau model lokal) lalu ekstraksi item dan nominal.
- Bot menampilkan usulan transaksi per item untuk dikonfirmasi sebelum disimpan.
- Foto tidak disimpan setelah diproses.

## Kriteria penerimaan
- [ ] Struk minimarket umum menghasilkan usulan item dengan total yang cocok.
- [ ] Tidak ada transaksi yang tersimpan tanpa konfirmasi pengguna.

## Dependensi
Gerbang Fase 4.
