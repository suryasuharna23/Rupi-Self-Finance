# RUPI-209 · Helm chart and deployment to kind

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 2 · gRPC microservices | P0 | L | Kubernetes infrastructure section |

## Scope
- A local kind cluster with the `rupi-dev` and `rupi-staging` namespaces.
- A Helm chart for all services, with replicas and requests/limits as listed in the table in the requirements document.
- PostgreSQL and Redis as StatefulSet + PVC (or CloudNativePG).
- Schema migrations as a Helm pre-upgrade hook Job.
- `readinessProbe` and `livenessProbe` use `grpc.health.v1`.
- Configuration via ConfigMap; bot token, LLM API key, and DB credentials via Secret.
- gRPC load balancing via a headless Service + client-side round robin.
- One multi-target Dockerfile or one Dockerfile per service; non-root images.

## Acceptance criteria
- [ ] `helm install` on kind brings every service to Ready.
- [ ] No secrets in values files or images.
- [ ] gRPC load is spread across both replicas (visible in the logs).

## Dependencies
RUPI-203 – RUPI-208.
