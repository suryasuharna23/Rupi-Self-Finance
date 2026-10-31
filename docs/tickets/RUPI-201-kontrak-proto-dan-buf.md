# RUPI-201 · Kontrak protobuf dan buf

| Fase | Prioritas | Ukuran | Referensi |
| --- | --- | --- | --- |
| 2 · Microservices gRPC | P0 | M | bagian Kontrak gRPC |

## Cakupan
- Direktori `proto/` dengan paket berversi `rupi.<service>.v1` untuk Conversation,
  Parser, Transaction, Report, Budget, dan User.
- Pesan inti sesuai dokumen: `ParseRequest`, `ParsedTransaction`, `ParseResult` (oneof),
  `ReportQuery`, enum `TransactionType`, `ParserSource`.
- `buf.yaml` + `buf.gen.yaml`; kode Go hasil generate di `gen/`.
- CI menjalankan `buf lint` dan `buf breaking` terhadap `main`.

## Kriteria penerimaan
- [ ] `buf generate` menghasilkan kode yang ter-compile.
- [ ] Mengubah nomor field yang sudah ada membuat CI gagal.
- [ ] Nominal memakai `int64`, waktu memakai `google.protobuf.Timestamp`.

## Dependensi
Gerbang Fase 1 (RUPI-104) sebaiknya sudah tercapai.
