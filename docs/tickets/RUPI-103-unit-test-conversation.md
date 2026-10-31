# RUPI-103 · Conversation Service unit tests with fakes

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 1 · Conversation core | P0 | M | FR-01 – FR-16, NFR-12 |

## Background
`conversation.Service` already uses interfaces for all of its dependencies, but
its flows (onboarding, clarification, correction, callbacks) are untested. Package
coverage is currently 19.7%.

## Scope
- In-memory fakes for `Users`, `Transactions`, `Reports`, `Idempotency`.
- Table-driven tests for every branch of `HandleMessage` and `handleCallback`.

## Acceptance criteria
- [ ] A duplicate message produces neither a reply nor a second transaction.
- [ ] A store error un-marks the message so a retry can succeed.
- [ ] Transactions with confidence < 0.7 are not saved and the bot asks a follow-up question (FR-12).
- [ ] An amount reply after "Berapa nominalnya?" is merged with the previous message.
- [ ] The 31st message within one minute is rejected (SEC-07).
- [ ] Coverage of the `conversation` package is ≥ 70%.

## Dependencies
None.
