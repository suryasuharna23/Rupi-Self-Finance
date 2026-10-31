# RUPI-204 · Split out the Transaction Service

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 2 · gRPC microservices | P0 | M | `TransactionService`, SEC-02, NFR-08 |

## Scope
- `cmd/transaction-service` with the RPCs `CreateTransactions`, `UpdateTransaction`,
  `DeleteTransaction`, `ListTransactions`, `FindRecent`.
- Only this service accesses the `transaction_svc` schema.
- `user_id` is always taken from the `x-user-id` metadata, never from the request body.
- `CreateTransactions` is idempotent based on `source_message_id` (unique index).

## Acceptance criteria
- [ ] Calling `CreateTransactions` twice with the same `source_message_id`
      returns `ALREADY_EXISTS` instead of creating duplicate transactions.
- [ ] User A cannot access user B's transactions through any RPC.
- [ ] The integration tests from RUPI-102 are moved over and still pass.

## Dependencies
RUPI-201, RUPI-202, RUPI-102.
