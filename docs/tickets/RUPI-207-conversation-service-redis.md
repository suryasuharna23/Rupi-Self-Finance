# RUPI-207 · Conversation Service dengan konteks di Redis

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 2 · Microservices gRPC | P0 | L | `ConversationService`, model data Redis |

## Cakupan
- `cmd/conversation-service` dengan `HandleMessage` (unary) dan
  `SubscribeOutbound` (server streaming) untuk pesan proaktif.
- Pindahkan `sessionStore` ke Redis: kunci `conv:<user_id>` berisi 10 transaksi
  terakhir dan pertanyaan klarifikasi yang menunggu, TTL 30 menit.
- Ganti panggilan paket langsung dengan client gRPC ke Parser, Transaction,
  Report, dan User.
- Rate limit per pengguna dipindah ke Redis agar berlaku lintas replika.

## Kriteria penerimaan
- [ ] Alur "makan siang" → "15rb" tetap jalan walau dua pesan ditangani replika berbeda.
- [ ] Service stateless; restart Pod tidak menghapus konteks percakapan.
- [ ] Klien `SubscribeOutbound` menerima pesan yang dikirim service lain.

## Dependensi
RUPI-203, RUPI-204, RUPI-205, RUPI-206.
