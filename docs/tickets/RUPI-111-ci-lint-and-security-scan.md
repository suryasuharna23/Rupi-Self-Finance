# RUPI-111 · CI: golangci-lint, coverage threshold, and Trivy scanning

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 1 · Conversation core | P0 | S | NFR-12, SEC-12 |

## Scope
- Add `golangci-lint` with its configuration in the repo.
- Fail CI when coverage of `internal/parser` or `internal/transaction` is < 70%.
- Scan dependencies and images with Trivy.
- Build images in CI tagged with the commit SHA (not `latest`).

## Acceptance criteria
- [ ] A PR with lint errors or coverage < 70% fails CI.
- [ ] Trivy findings rated HIGH/CRITICAL fail CI.
- [ ] Images are pushed with the SHA tag.

## Dependencies
None.
