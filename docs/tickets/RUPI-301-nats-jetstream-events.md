# RUPI-301 · Transaction events via NATS JetStream

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 3 · Events, budgets, reports | P1 | M | Architecture (event path), NATS infrastructure |

## Scope
- NATS JetStream as a StatefulSet (1–3 replicas) in the Helm chart.
- The Transaction Service publishes `transaction.created`, `transaction.updated`,
  and `transaction.deleted` in protobuf format.
- The outbox pattern, so events are not lost if publishing fails after the database commit.

## Acceptance criteria
- [ ] Every transaction change produces exactly one event.
- [ ] Briefly shutting down NATS does not lose events (they are delivered after recovery).
- [ ] Consumers can safely reprocess events (idempotent).

## Dependencies
MVP gate (RUPI-211).
