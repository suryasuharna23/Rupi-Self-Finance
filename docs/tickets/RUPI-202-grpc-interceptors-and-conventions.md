# RUPI-202 · gRPC interceptors and error conventions

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 2 · gRPC microservices | P0 | M | General rules of the gRPC contracts section, SEC-02 |

## Scope
- A shared `internal/grpcx` package containing server and client interceptors that:
  - propagate the `x-request-id` and `x-user-id` metadata into the context, logs, and traces;
  - emit structured logs per RPC (method, duration, status code) without message content;
  - recover from panics.
- Mapping of domain errors to gRPC statuses: `INVALID_ARGUMENT`, `NOT_FOUND`,
  `ALREADY_EXISTS`, `DEADLINE_EXCEEDED`, `UNAVAILABLE`.
- Retry policy: idempotent RPCs only, at most 3 attempts with exponential backoff.
- A `grpc.health.v1` health check server.

## Acceptance criteria
- [ ] Requests to Transaction/Report without `x-user-id` are rejected with `UNAUTHENTICATED`.
- [ ] The same `x-request-id` appears in the logs of every service for a single message.
- [ ] Unit tests for retries: non-idempotent RPCs are never retried.

## Dependencies
RUPI-201.
