# RUPI-103 · Unit test Conversation Service dengan fake

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 1 · Inti percakapan | P0 | M | FR-01 – FR-16, NFR-12 |

## Latar
`conversation.Service` sudah memakai interface untuk semua dependensinya, tetapi
alurnya (onboarding, klarifikasi, koreksi, callback) belum dites. Cakupan paket
saat ini 19,7%.

## Cakupan
- Fake in-memory untuk `Users`, `Transactions`, `Reports`, `Idempotency`.
- Tes tabel untuk setiap cabang `HandleMessage` dan `handleCallback`.

## Kriteria penerimaan
- [ ] Pesan duplikat tidak menghasilkan balasan maupun transaksi kedua.
- [ ] Error dari store membuat pesan di-unmark agar retry bisa berhasil.
- [ ] Transaksi dengan confidence < 0,7 tidak disimpan dan bot bertanya balik (FR-12).
- [ ] Jawaban nominal setelah "Berapa nominalnya?" digabung dengan pesan sebelumnya.
- [ ] Pesan ke-31 dalam satu menit ditolak (SEC-07).
- [ ] Cakupan paket `conversation` ≥ 70%.

## Dependensi
Tidak ada.
