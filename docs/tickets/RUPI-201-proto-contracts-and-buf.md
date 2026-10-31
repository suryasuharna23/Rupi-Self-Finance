# RUPI-201 · Protobuf contracts and buf

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 2 · gRPC microservices | P0 | M | gRPC contracts section |

## Scope
- A `proto/` directory with versioned `rupi.<service>.v1` packages for Conversation,
  Parser, Transaction, Report, Budget, and User.
- Core messages as defined in the requirements document: `ParseRequest`, `ParsedTransaction`, `ParseResult` (oneof),
  `ReportQuery`, and the `TransactionType` and `ParserSource` enums.
- `buf.yaml` + `buf.gen.yaml`; generated Go code in `gen/`.
- CI runs `buf lint` and `buf breaking` against `main`.

## Acceptance criteria
- [ ] `buf generate` produces code that compiles.
- [ ] Changing an existing field number makes CI fail.
- [ ] Amounts use `int64`; timestamps use `google.protobuf.Timestamp`.

## Dependencies
The Phase 1 gate (RUPI-104) should already be met.
