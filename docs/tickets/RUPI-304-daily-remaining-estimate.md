# RUPI-304 · Daily remaining money estimate

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 3 · Events, budgets, reports | P1 | S | FR-32 |

## Scope
- Compute the remaining budget divided by the remaining days in the period, for example
  "sisa Rp340rb untuk 12 hari, aman ~Rp28rb/hari".
- Show it in `/budget`, in `/ringkas`, and in budget alert replies.

## Acceptance criteria
- [ ] Unit tests for month end, leap years, and budgets that are already exceeded.
- [ ] Rounded down to the nearest thousand.

## Dependencies
RUPI-302.
