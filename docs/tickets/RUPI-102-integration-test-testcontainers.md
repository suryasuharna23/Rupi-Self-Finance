# RUPI-102 · Store integration tests with Testcontainers

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 1 · Conversation core | P0 | M | NFR-07, NFR-08, SEC-02, Testing section |

## Background
None of the SQL queries have automated tests yet. The requirements document specifies
that integration tests use a real Postgres instance via Testcontainers.

## Scope
- A test helper that starts `postgres:16-alpine` and runs `storage.Migrate`.
- Tests for `transaction.Store`, `user.Store`, `report.Store`, `conversation.IdempotencyStore`.
- Separate them with the `integration` build tag so `go test ./...` stays fast.
- A separate CI job for integration tests.

## Acceptance criteria
- [x] `Create` saves multiple transactions atomically and maps category names to IDs.
- [x] `UpdateAmount` writes an audit trail to `transaction_edits`.
- [x] `DeleteLast` does not delete transactions created more than 24 hours ago.
- [x] User A cannot read, modify, or delete user B's transactions (SEC-02).
- [x] `MarkProcessed` returns `false` the second time it sees the same message (NFR-07).
- [x] `GetOrCreate` creates a new user after the old account has been soft-deleted.
- [x] `Summary` is correct for a week that spans the start of a month.

## Dependencies
RUPI-101.
