# RUPI-204 · Pecah Transaction Service

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 2 · Microservices gRPC | P0 | M | `TransactionService`, SEC-02, NFR-08 |

## Cakupan
- `cmd/transaction-service` dengan RPC `CreateTransactions`, `UpdateTransaction`,
  `DeleteTransaction`, `ListTransactions`, `FindRecent`.
- Hanya service ini yang mengakses schema `transaction_svc`.
- `user_id` selalu diambil dari metadata `x-user-id`, tidak pernah dari body request.
- `CreateTransactions` idempoten berdasarkan `source_message_id` (unique index).

## Kriteria penerimaan
- [ ] Memanggil `CreateTransactions` dua kali dengan `source_message_id` yang sama
      menghasilkan `ALREADY_EXISTS`, bukan transaksi ganda.
- [ ] Pengguna A tidak bisa mengakses transaksi pengguna B lewat RPC apa pun.
- [ ] Integration test dari RUPI-102 dipindah dan tetap lulus.

## Dependensi
RUPI-201, RUPI-202, RUPI-102.
