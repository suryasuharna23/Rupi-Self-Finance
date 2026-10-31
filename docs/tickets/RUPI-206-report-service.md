# RUPI-206 · Pecah Report Service

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 2 · Microservices gRPC | P0 | M | `ReportService.RunQuery`, FR-20, FR-21 |

## Latar
Dokumen menyatakan service tidak membaca tabel milik service lain secara
langsung, tetapi diagram arsitektur menunjukkan Report Service membaca
PostgreSQL. Putuskan dulu pendekatannya.

## Cakupan
- Putuskan: (a) read-only role ke schema transaksi, (b) memanggil
  `TransactionService.ListTransactions`, atau (c) proyeksi sendiri dari event NATS (Fase 3).
  Catat keputusannya di `docs/`.
- `cmd/report-service` dengan `RunQuery(ReportQuery) → ReportResult`.

## Kriteria penerimaan
- [ ] `/ringkas` dan pertanyaan total memberi angka yang sama seperti versi monolit.
- [ ] Setiap query difilter `x-user-id` (SEC-02).

## Dependensi
RUPI-201, RUPI-202, RUPI-204.
