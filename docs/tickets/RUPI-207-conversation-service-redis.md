# RUPI-207 · Conversation Service with context in Redis

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 2 · gRPC microservices | P0 | L | `ConversationService`, Redis data model |

## Scope
- `cmd/conversation-service` with `HandleMessage` (unary) and
  `SubscribeOutbound` (server streaming) for proactive messages.
- Move `sessionStore` to Redis: the key `conv:<user_id>` holds the last 10
  transactions and any pending clarification question, with a 30-minute TTL.
- Replace direct package calls with gRPC clients for Parser, Transaction,
  Report, and User.
- Move per-user rate limiting to Redis so it applies across replicas.

## Acceptance criteria
- [ ] The "makan siang" → "15rb" flow still works even when the two messages are handled by different replicas.
- [ ] The service is stateless; restarting a Pod does not erase conversation context.
- [ ] `SubscribeOutbound` clients receive messages sent by other services.

## Dependencies
RUPI-203, RUPI-204, RUPI-205, RUPI-206.
