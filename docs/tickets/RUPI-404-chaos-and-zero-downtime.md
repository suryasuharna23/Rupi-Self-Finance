# RUPI-404 · Basic chaos testing and zero-downtime rolling updates

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 4 · Production readiness | P1 | M | NFR-07, NFR-08, NFR-10 |

## Scope
- Graceful shutdown in every service: stop accepting requests, finish in-flight
  ones, within at most 30 seconds.
- A script that runs random `kubectl delete pod` on a schedule during the load test.
- A rolling update of the Parser Service in the middle of the load test.

## Acceptance criteria
- [ ] No lost or duplicate transactions during chaos.
- [ ] A rolling update of the Parser Service during the load test causes no lost messages.

## Dependencies
RUPI-403.
