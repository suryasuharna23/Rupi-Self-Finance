# RUPI-403 · Load test at 50 messages/second

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 4 · Production readiness | P1 | M | NFR-01, NFR-06, Phase 4 gate |

## Scope
- A k6 (or ghz) script that sends fake webhooks from 1,000 simulated users.
- 50 messages/second for 10 minutes, a mix of rule-path and (mocked) LLM sentences.
- The results report is stored in `docs/`.

## Acceptance criteria
- [ ] Rule-path p95 < 500 ms.
- [ ] No lost messages and no duplicate transactions (reconciled against the database).
- [ ] The Parser Service HPA scales up and down during the test.

## Dependencies
RUPI-401, RUPI-402.
