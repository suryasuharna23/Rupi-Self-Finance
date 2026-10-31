# RUPI-108 · Complete onboarding: time zone and opening balance

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 1 · Conversation core | P0 | S | FR-02 |

## Background
Onboarding currently only asks for the user's name; the time zone is always WIB and
there is no opening balance yet.

## Scope
- A time zone selection step with WIB / WITA / WIT buttons (default WIB).
- An optional opening balance step with a "Lewati" (skip) button; it is saved as income
  in the Lainnya category with the description "Saldo awal".
- A `/zonawaktu` command to change the time zone after onboarding.

## Acceptance criteria
- [ ] A WITA user who logs "kemarin" gets yesterday's date according to WITA.
- [ ] The opening balance can be skipped without creating a transaction.
- [ ] `/zonawaktu` updates `user_svc.users.timezone`.

## Dependencies
RUPI-101.
