# Rupi Tickets

These tickets are derived from the [system requirements document (in Indonesian)](../Rupi%20-%20Dokumen%20Kebutuhan%20Sistem.pdf)
and follow a five-phase roadmap. Each phase starts only after the previous phase's
gate has been met. Sizes: **S** ≈ 1–2 days, **M** ≈ 3–5 days, **L** ≈ 1–2 weeks.

**Ready to start now** (no dependencies): RUPI-101, RUPI-103, RUPI-104, RUPI-107, RUPI-111.
RUPI-105 can start once the LLM provider has been chosen.

## Phase 1 · Conversation core (monolith)
Gate: parser accuracy ≥ 95% on the 300-sentence evaluation dataset (NFR-03).

| Ticket | Title | Priority | Size | Dependencies |
| --- | --- | --- | --- | --- |
| [RUPI-101](RUPI-101-local-end-to-end-verification.md) | Verify the bot runs end-to-end locally | P0 | S | — |
| [RUPI-102](RUPI-102-integration-test-testcontainers.md) | Store integration tests with Testcontainers | P0 | M | 101 |
| [RUPI-103](RUPI-103-unit-test-conversation.md) | Conversation Service unit tests with fakes | P0 | M | — |
| [RUPI-104](RUPI-104-evaluation-dataset-300.md) | 300-sentence parser evaluation dataset | P0 | L | — |
| [RUPI-105](RUPI-105-llm-fallback-implementation.md) | Implement the LLM as a parser fallback | P0 | L | LLM decision |
| [RUPI-106](RUPI-106-llm-limits-and-cost.md) | Daily LLM limit and cost tracking | P0 | S | 105 |
| [RUPI-107](RUPI-107-parser-logs.md) | parser_logs table with redacted text | P0 | M | — |
| [RUPI-108](RUPI-108-complete-onboarding.md) | Complete onboarding: time zone and opening balance | P0 | S | 101 |
| [RUPI-109](RUPI-109-account-purge.md) | Permanently purge accounts after 30 days | P0 | S | 102 |
| [RUPI-110](RUPI-110-category-description-correction.md) | Correct category, description, and date | P0 | M | 101 |
| [RUPI-111](RUPI-111-ci-lint-and-security-scan.md) | CI: golangci-lint, coverage threshold, Trivy | P0 | S | — |

## Phase 2 · Split into gRPC microservices
MVP gate: all P0 acceptance criteria pass end-to-end on kind.

| Ticket | Title | Priority | Size | Dependencies |
| --- | --- | --- | --- | --- |
| [RUPI-201](RUPI-201-proto-contracts-and-buf.md) | Protobuf contracts and buf | P0 | M | Phase 1 gate |
| [RUPI-202](RUPI-202-grpc-interceptors-and-conventions.md) | gRPC interceptors and error conventions | P0 | M | 201 |
| [RUPI-203](RUPI-203-parser-service.md) | Split out the Parser Service | P0 | M | 201, 202 |
| [RUPI-204](RUPI-204-transaction-service.md) | Split out the Transaction Service | P0 | M | 201, 202, 102 |
| [RUPI-205](RUPI-205-user-service.md) | Split out the User Service | P0 | S | 201, 202 |
| [RUPI-206](RUPI-206-report-service.md) | Split out the Report Service | P0 | M | 204 |
| [RUPI-207](RUPI-207-conversation-service-redis.md) | Conversation Service with context in Redis | P0 | L | 203–206 |
| [RUPI-208](RUPI-208-channel-adapter.md) | Channel Adapter as a separate service | P0 | M | 207 |
| [RUPI-209](RUPI-209-helm-chart-kind.md) | Helm chart and deployment to kind | P0 | L | 203–208 |
| [RUPI-210](RUPI-210-inter-service-mtls.md) | Inter-service mTLS | P0 | M | 209 |
| [RUPI-211](RUPI-211-e2e-test-kind.md) | End-to-end tests on kind (MVP gate) | P0 | M | 209, 210 |

## Phase 3 · Events, budgets, and reports
Gate: budget alerts are sent exactly once per period.

| Ticket | Title | Priority | Size | Dependencies |
| --- | --- | --- | --- | --- |
| [RUPI-301](RUPI-301-nats-jetstream-events.md) | Transaction events via NATS JetStream | P1 | M | 211 |
| [RUPI-302](RUPI-302-budget-service.md) | Budget Service and setting budgets via chat | P1 | M | 301 |
| [RUPI-303](RUPI-303-budget-alerts.md) | 80% and 100% budget alerts | P1 | M | 302 |
| [RUPI-304](RUPI-304-daily-remaining-estimate.md) | Daily remaining money estimate | P1 | S | 302 |
| [RUPI-305](RUPI-305-daily-reminder.md) | Optional daily reminder | P1 | S | 306 |
| [RUPI-306](RUPI-306-periodic-reports-cronjob.md) | Automatic weekly and monthly reports | P1 | M | 206, 207 |
| [RUPI-307](RUPI-307-period-comparison.md) | Comparison with the previous period | P1 | S | 306 |
| [RUPI-308](RUPI-308-custom-categories.md) | Custom categories and keyword rules | P1 | M | 203, 204 |
| [RUPI-309](RUPI-309-advanced-queries.md) | Advanced queries: count, largest, breakdown | P0 | M | 206 |

## Phase 4 · Production readiness
Gate: passes the load test with no lost or duplicate messages (NFR-06, NFR-07).

| Ticket | Title | Priority | Size | Dependencies |
| --- | --- | --- | --- | --- |
| [RUPI-401](RUPI-401-observability.md) | Observability: metrics, logs, and tracing | P1 | L | 209 |
| [RUPI-402](RUPI-402-hpa-pdb-networkpolicy.md) | HPA, PodDisruptionBudget, NetworkPolicy | P1 | M | 209 |
| [RUPI-403](RUPI-403-load-test.md) | Load test at 50 messages/second | P1 | M | 401, 402 |
| [RUPI-404](RUPI-404-chaos-and-zero-downtime.md) | Basic chaos testing and rolling updates | P1 | M | 403 |
| [RUPI-405](RUPI-405-backup-and-encryption.md) | Daily backups and encryption at rest | P1 | M | 209 |

## Phase 5 · Advanced features
Order is chosen based on feedback from test users.

| Ticket | Title | Priority | Size | Dependencies |
| --- | --- | --- | --- | --- |
| [RUPI-501](RUPI-501-web-chat-and-dashboard.md) | Web chat and chart dashboard | P2 | L | Phase 4 gate |
| [RUPI-502](RUPI-502-receipt-photo-ocr.md) | Receipt photo → OCR → per-item transactions | P2 | L | Phase 4 gate |
| [RUPI-503](RUPI-503-whatsapp-channel.md) | WhatsApp channel | P2 | L | Phase 4 gate |
| [RUPI-504](RUPI-504-shared-wallet.md) | Shared wallet | P2 | L | Phase 4 gate |
| [RUPI-505](RUPI-505-csv-export.md) | Export transactions to CSV | P2 | S | 206 |
