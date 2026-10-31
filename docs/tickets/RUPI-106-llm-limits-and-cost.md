# RUPI-106 · Daily LLM limit and cost tracking

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 1 · Conversation core | P0 | S | SEC-07, NFR-11 |

## Scope
- A limit of 500 LLM calls per user per day; after that, only the rule path is used.
- Store accumulated tokens per user per day to compute the cost per active user
  (target < Rp50 per active user per month).

## Acceptance criteria
- [ ] The 501st call in a day does not call the LLM, and the user gets a follow-up question.
- [ ] There is a way to view total tokens per day (SQL query or structured log).

## Dependencies
RUPI-105.
