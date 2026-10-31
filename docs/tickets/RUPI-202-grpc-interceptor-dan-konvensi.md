# RUPI-202 · Interceptor gRPC dan konvensi error

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 2 · Microservices gRPC | P0 | M | Aturan umum Kontrak gRPC, SEC-02 |

## Cakupan
- Paket bersama `internal/grpcx` berisi interceptor server dan client:
  - meneruskan metadata `x-request-id` dan `x-user-id` ke context, log, dan trace;
  - logging terstruktur per RPC (method, durasi, kode status) tanpa isi pesan;
  - recovery dari panic.
- Pemetaan error domain ke status gRPC: `INVALID_ARGUMENT`, `NOT_FOUND`,
  `ALREADY_EXISTS`, `DEADLINE_EXCEEDED`, `UNAVAILABLE`.
- Kebijakan retry: hanya RPC idempoten, maksimal 3 kali dengan exponential backoff.
- Server health check `grpc.health.v1`.

## Kriteria penerimaan
- [ ] Request tanpa `x-user-id` ke Transaction/Report ditolak `UNAUTHENTICATED`.
- [ ] `x-request-id` yang sama muncul di log semua service untuk satu pesan.
- [ ] Unit test untuk retry: RPC non-idempoten tidak pernah di-retry.

## Dependensi
RUPI-201.
