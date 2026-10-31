# RUPI-101 · Verify the bot runs end-to-end locally

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 1 · Conversation core | P0 | S | FR-01 – FR-16 |

## Background
The Phase 1 setup has only been verified up to compilation and unit tests. The database
code (`internal/transaction`, `user`, `report`, `conversation/store.go`) and the
Telegram flow have never been run against a real PostgreSQL instance and a real bot.

## Scope
- Install Docker Desktop, run `make db`, fill in `.env`, run `make run`.
- Manually test every flow in Telegram and fix any bugs that come up.
- Document Windows/PowerShell-specific steps in the README.

## Acceptance criteria
- [ ] Migrations run cleanly on an empty database and are safe to re-run.
- [ ] `/start` → privacy policy → ask for name → onboarding complete.
- [ ] "beli bakso 15rb" is saved and the reply includes Ubah / Batalkan buttons.
- [ ] "bakso 15rb sama es teh goceng" produces exactly 2 transactions.
- [ ] "makan siang" → "Berapa nominalnya?" → "15rb" saves a lunch transaction.
- [ ] "yang bakso tadi 18rb" changes the amount and writes a `transaction_edits` row.
- [ ] The Batalkan button, `/undo`, `/ringkas`, and `/hapusakun` work.
- [ ] "bulan ini habis berapa buat makan?" is answered with the correct total.
- [ ] `GET /healthz` returns 200 when the database is up and 503 when it is down.

## Dependencies
None. Do this first.
