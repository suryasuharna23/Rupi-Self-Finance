# RUPI-109 · Permanently purge accounts after 30 days

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 1 · Conversation core | P0 | S | FR-03, SEC-10 |

## Scope
- A scheduled job (a daily goroutine in Phase 1, a CronJob in Phase 2) that permanently
  deletes `users` rows whose `deleted_at` is more than 30 days old.
- Ensure no user data remains in other schemas (transactions, custom categories,
  parser_logs, budgets).

## Acceptance criteria
- [ ] An account deleted 31 days ago is completely gone; one deleted 29 days ago still exists.
- [ ] An integration test that checks every table owned by the user.

## Dependencies
RUPI-102.
