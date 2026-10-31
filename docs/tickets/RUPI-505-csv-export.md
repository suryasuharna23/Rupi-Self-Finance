# RUPI-505 · Export transactions to CSV

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 5 · Advanced features | P2 | S | FR-24 |

## Scope
- An `/ekspor` command with a date range ("/ekspor bulan ini", "/ekspor 1-15 okt").
- Send the CSV file via Telegram (`sendDocument`).
- Columns: date, type, category, description, amount.

## Acceptance criteria
- [ ] The file opens correctly in Excel and Google Sheets (UTF-8 with BOM).
- [ ] It contains only the requesting user's transactions.

## Dependencies
RUPI-206.
