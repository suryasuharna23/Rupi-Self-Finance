# RUPI-211 · End-to-end tests on kind (MVP gate)

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 2 · gRPC microservices | P0 | M | Testing section: End-to-end, MVP gate |

## Scope
- A fake Telegram script that sends webhooks to the Channel Adapter and captures the replies.
- The LLM is mocked with fixed responses.
- Run all P0 acceptance criteria from the requirements document.

## Acceptance criteria
- [ ] "bakso 15rb sama es teh goceng" produces exactly 2 transactions of Rp15.000 and Rp5.000.
- [ ] The same webhook delivered twice produces only one set of transactions.
- [ ] "yang bakso tadi 18rb" changes the most recent bakso transaction and records an audit trail in `transaction_edits`.
- [ ] When the LLM does not respond within 5 seconds, the user receives a clarification question, not an error.
- [ ] User A cannot view or modify user B's transactions through any RPC.
- [ ] The tests run in CI.

## Dependencies
RUPI-209, RUPI-210.
