# RUPI-205 · Split out the User Service

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 2 · gRPC microservices | P0 | S | `UserService`, FR-01 – FR-03 |

## Scope
- `cmd/user-service` with the RPCs `GetOrCreateUser`, `UpdateProfile`, `DeleteUser`.
- Only this service accesses the `user_svc` schema.
- Account deletion coordinates the deletion of data in the Transaction Service.

## Acceptance criteria
- [ ] `GetOrCreateUser` is safe to call concurrently for the same Telegram ID.
- [ ] `/hapusakun` still deletes all of the user's data after the split.

## Dependencies
RUPI-201, RUPI-202.
