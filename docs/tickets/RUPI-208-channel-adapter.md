# RUPI-208 · Channel Adapter as a separate service

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 2 · gRPC microservices | P0 | M | SEC-01, NFR-07, single-message flow |

## Scope
- `cmd/channel-adapter` containing `internal/telegram`; the only service
  exposed to the internet (Ingress).
- Check for duplicate message IDs in Redis before calling `ConversationService.HandleMessage`.
- Subscribe to `SubscribeOutbound` to deliver alerts and reports to Telegram.
- Webhook mode with `secret_token` verification as the default in the cluster.

## Acceptance criteria
- [ ] Webhooks without a valid `secret_token` are rejected with 403.
- [ ] The same webhook delivered twice produces only one set of transactions.
- [ ] Proactive messages from the stream are delivered to the correct chat.

## Dependencies
RUPI-207.
