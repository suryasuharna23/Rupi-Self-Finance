# RUPI-502 · Receipt photo → OCR → per-item transactions

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 5 · Advanced features | P2 | L | FR-18 |

## Scope
- The Channel Adapter accepts photos from Telegram.
- OCR (an external service or a local model), followed by extraction of items and amounts.
- The bot shows proposed per-item transactions for confirmation before saving.
- Photos are not stored after processing.

## Acceptance criteria
- [ ] A typical minimarket receipt produces proposed items whose total matches.
- [ ] No transaction is saved without the user's confirmation.

## Dependencies
Phase 4 gate.
