# RUPI-401 · Observability: metrics, logs, and tracing

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 4 · Production readiness | P1 | L | NFR-09, NFR-11 |

## Scope
- A `monitoring` namespace with Prometheus, Grafana, Loki, and the OpenTelemetry Collector.
- OTel instrumentation in the gRPC interceptors; the trace ID equals `x-request-id`.
- Metrics: per-RPC latency, error rate, share of rule path vs LLM, LLM tokens and cost.
- Grafana dashboards: service health, LLM cost, parser accuracy (from `parser_logs`).
- Basic alerts: high error rate, p95 latency exceeding NFR-01/NFR-02.

## Acceptance criteria
- [ ] A single message can be traced from the Channel Adapter to PostgreSQL in one trace.
- [ ] The dashboard shows p95 latency for the rule path and the LLM path.
- [ ] The dashboard shows LLM cost per active user.

## Dependencies
RUPI-209.
