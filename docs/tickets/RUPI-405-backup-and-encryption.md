# RUPI-405 · Daily backups and encryption at rest

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 4 · Production readiness | P1 | M | SEC-09, SEC-10 |

## Scope
- Daily PostgreSQL backups with 7-day retention (CloudNativePG or a `pg_dump` CronJob).
- Encrypted database and backup volumes.
- A documented restore procedure that has been tested.

## Acceptance criteria
- [ ] Restoring yesterday's backup succeeds in a separate namespace.
- [ ] Backups older than 7 days are deleted automatically, so data from deleted accounts
      disappears from backups within 30 days (SEC-10).

## Dependencies
RUPI-209.
