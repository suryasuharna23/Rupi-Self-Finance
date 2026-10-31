# RUPI-203 · Split out the Parser Service

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 2 · gRPC microservices | P0 | M | `ParserService.Parse`, NFR-01, NFR-02 |

## Scope
- `cmd/parser-service`, which wraps `internal/parser` behind `ParserService.Parse`.
- Conversion between `parser.Result` ↔ the protobuf `ParseResult`, including `source` for metrics.
- A 5-second deadline set by the caller; only this service holds the LLM API key.

## Acceptance criteria
- [ ] The evaluation dataset gives the same results via gRPC as via direct calls.
- [ ] Calls that exceed the deadline return `DEADLINE_EXCEEDED`.
- [ ] The service is stateless and can run with multiple replicas.

## Dependencies
RUPI-201, RUPI-202.
