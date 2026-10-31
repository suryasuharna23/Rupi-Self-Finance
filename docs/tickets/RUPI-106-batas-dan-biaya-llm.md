# RUPI-106 · Batas harian dan pencatatan biaya LLM

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 1 · Inti percakapan | P0 | S | SEC-07, NFR-11 |

## Cakupan
- Batas 500 panggilan LLM per pengguna per hari; setelahnya hanya jalur aturan.
- Simpan akumulasi token per pengguna per hari untuk menghitung biaya per pengguna aktif
  (target < Rp50 per pengguna aktif per bulan).

## Kriteria penerimaan
- [ ] Panggilan ke-501 dalam sehari tidak memanggil LLM dan pengguna mendapat tanya balik.
- [ ] Ada cara melihat total token per hari (query SQL atau log terstruktur).

## Dependensi
RUPI-105.
