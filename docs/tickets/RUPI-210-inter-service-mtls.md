# RUPI-210 · Inter-service mTLS

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 2 · gRPC microservices | P0 | M | SEC-03 |

## Scope
- Choose between cert-manager (mTLS in code) or Linkerd (service mesh). Linkerd also
  solves HTTP/2 load balancing.
- All inter-service communication is encrypted and authenticated.

## Acceptance criteria
- [ ] A Pod without a valid certificate cannot call internal services.
- [ ] Certificate rotation does not interrupt service.

## Dependencies
RUPI-209.
