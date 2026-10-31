# RUPI-206 · Split out the Report Service

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 2 · gRPC microservices | P0 | M | `ReportService.RunQuery`, FR-20, FR-21 |

## Background
The requirements document states that services do not read tables owned by other services
directly, yet the architecture diagram shows the Report Service reading from
PostgreSQL. Decide on the approach first.

## Scope
- Decide between: (a) a read-only role on the transaction schema, (b) calling
  `TransactionService.ListTransactions`, or (c) its own projection built from NATS events (Phase 3).
  Record the decision in `docs/`.
- `cmd/report-service` with `RunQuery(ReportQuery) → ReportResult`.

## Acceptance criteria
- [ ] `/ringkas` and total questions return the same numbers as the monolith version.
- [ ] Every query is filtered by `x-user-id` (SEC-02).

## Dependencies
RUPI-201, RUPI-202, RUPI-204.
