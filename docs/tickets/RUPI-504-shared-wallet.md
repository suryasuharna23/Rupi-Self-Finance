# RUPI-504 · Shared wallet

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 5 · Advanced features | P2 | L | FR-40, FR-41, FR-42 |

## Scope
- Create a shared wallet and invite members via a link (`wallets`, `wallet_members` tables).
- Shared wallet transactions record who paid and how it is split.
- Calculate who owes whom (debt simplification) and send reminders.

## Acceptance criteria
- [ ] Members only see transactions in the shared wallets they belong to.
- [ ] The debts of all members always sum to zero.
- [ ] Bill reminders can be turned off per member.

## Dependencies
Phase 4 gate.
