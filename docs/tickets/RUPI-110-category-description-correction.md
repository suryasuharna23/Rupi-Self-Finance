# RUPI-110 · Correct category, description, and date

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 1 · Conversation core | P0 | M | FR-13, FR-14 |

## Background
Corrections can currently only change the amount, and the "Ubah" button only shows
a text hint.

## Scope
- Support sentences such as "yang bensin tadi masuk transportasi" and "yang bakso tadi kemarin".
- The "Ubah" button shows Nominal / Kategori / Tanggal options; the category is chosen via buttons.
- Every change is recorded in `transaction_edits` along with the field name.

## Acceptance criteria
- [ ] Category, description, and date can be corrected both via sentences and via buttons.
- [ ] There is a per-field audit trail.
- [ ] New example correction sentences are added to the evaluation dataset.

## Dependencies
RUPI-101.
