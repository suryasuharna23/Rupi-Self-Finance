# RUPI-301 · Event transaksi lewat NATS JetStream

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 3 · Event, budget, laporan | P1 | M | Arsitektur (jalur event), Infrastruktur NATS |

## Cakupan
- NATS JetStream sebagai StatefulSet (1–3 replika) di Helm chart.
- Transaction Service menerbitkan `transaction.created`, `transaction.updated`,
  `transaction.deleted` dalam format protobuf.
- Pola outbox agar event tidak hilang bila publish gagal setelah commit database.

## Kriteria penerimaan
- [ ] Setiap perubahan transaksi menghasilkan tepat satu event.
- [ ] Mematikan NATS sebentar tidak menghilangkan event (terkirim setelah pulih).
- [ ] Konsumen bisa memproses ulang event dengan aman (idempoten).

## Dependensi
Gerbang MVP (RUPI-211).
