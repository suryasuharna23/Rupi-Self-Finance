# RUPI-203 · Pecah Parser Service

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 2 · Microservices gRPC | P0 | M | `ParserService.Parse`, NFR-01, NFR-02 |

## Cakupan
- `cmd/parser-service` yang membungkus `internal/parser` di balik `ParserService.Parse`.
- Konversi `parser.Result` ↔ `ParseResult` protobuf, termasuk `source` untuk metrik.
- Deadline 5 detik dari pemanggil; LLM API key hanya dimiliki service ini.

## Kriteria penerimaan
- [ ] Dataset evaluasi memberi hasil yang sama lewat gRPC maupun pemanggilan langsung.
- [ ] Panggilan melewati deadline mengembalikan `DEADLINE_EXCEEDED`.
- [ ] Service stateless dan bisa berjalan dengan beberapa replika.

## Dependensi
RUPI-201, RUPI-202.
