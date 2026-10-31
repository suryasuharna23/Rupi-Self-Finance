# RUPI-302 · Budget Service and setting budgets via chat

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 3 · Events, budgets, reports | P1 | M | FR-30, `BudgetService` |

## Scope
- `cmd/budget-service` with `SetBudget` and `GetBudgetStatus`; a `budget_svc` schema
  with `budgets` and `budget_alerts` tables.
- The parser recognizes budget sentences: "budget makan 500rb sebulan", "budget total 1jt per minggu".
- The `/budget` command shows usage for each budget.

## Acceptance criteria
- [ ] Per-category and total budgets, weekly or monthly, can be set via chat.
- [ ] `/budget` shows correct usage and remaining amounts.
- [ ] Example budget sentences are added to the evaluation dataset.

## Dependencies
RUPI-301.
