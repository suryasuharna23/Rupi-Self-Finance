# RUPI-402 · HPA, PodDisruptionBudget, and NetworkPolicy

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 4 · Production readiness | P1 | M | NFR-05, NFR-06, Mandatory Deployment requirements |

## Scope
- HPA for the Parser Service: 2–10 Pods, 70% CPU target.
- A `PodDisruptionBudget` with `minAvailable: 1` for every Deployment.
- A default-deny `NetworkPolicy`; only the call paths in the architecture diagram are allowed.
- All containers run as non-root with a read-only filesystem (SEC-12).

## Acceptance criteria
- [ ] Calls not shown in the diagram (for example Channel Adapter → PostgreSQL) are blocked.
- [ ] Running `kubectl drain` on a node does not bring any service down to 0 replicas.
- [ ] The HPA scales up under load and back down afterwards.

## Dependencies
RUPI-209.
