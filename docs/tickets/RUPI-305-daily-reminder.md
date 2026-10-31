# RUPI-305 · Optional daily reminder

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 3 · Events, budgets, reports | P1 | S | FR-33 |

## Scope
- A `/pengingat on|off` command and a reminder time (default 20.00 in the user's time zone).
- The reminder is only sent if the user has not logged anything that day.

## Acceptance criteria
- [ ] Users who have already logged something today do not receive a reminder.
- [ ] The reminder time follows the user's time zone.
- [ ] Off by default; users must turn it on themselves.

## Dependencies
RUPI-306 (uses the same scheduling mechanism).
