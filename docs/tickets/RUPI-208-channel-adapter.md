# RUPI-208 · Channel Adapter sebagai service terpisah

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 2 · Microservices gRPC | P0 | M | SEC-01, NFR-07, alur satu pesan |

## Cakupan
- `cmd/channel-adapter` berisi `internal/telegram`; satu-satunya service yang
  terbuka ke internet (Ingress).
- Cek duplikasi ID pesan di Redis sebelum memanggil `ConversationService.HandleMessage`.
- Berlangganan `SubscribeOutbound` untuk mengirim peringatan dan laporan ke Telegram.
- Mode webhook dengan verifikasi `secret_token` sebagai default di cluster.

## Kriteria penerimaan
- [ ] Webhook tanpa `secret_token` valid ditolak 403.
- [ ] Webhook yang sama dikirim dua kali hanya menghasilkan satu set transaksi.
- [ ] Pesan proaktif dari stream terkirim ke chat yang benar.

## Dependensi
RUPI-207.
