# RUPI-303 · 80% and 100% budget alerts

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 3 · Events, budgets, reports | P1 | M | FR-31, Phase 3 gate |

## Scope
- The Budget Service consumes transaction events and recalculates usage.
- Alerts are sent via `SubscribeOutbound` when usage crosses 80% and 100%.
- A `budget_alerts` table (budget_id, period_start, level) with a unique constraint
  prevents duplicate alerts.
- Deleted or corrected transactions also reduce usage.

## Acceptance criteria
- [ ] The 80% alert is sent exactly once per period, even if events are reprocessed.
- [ ] A correction that brings usage below a threshold does not trigger a new alert.
- [ ] A new period resets the alert state.

## Dependencies
RUPI-302.
