# RUPI-107 · parser_logs table with redacted text

| Phase | Priority | Size | References |
| --- | --- | --- | --- |
| 1 · Conversation core | P0 | M | SEC-08, `parser_logs` data model |

## Scope
- Migration for the `parser_svc` schema with a `parser_logs` table (user_id hash, text_redacted,
  result JSON, source, latency_ms, tokens, created_at).
- Redaction function: long numbers (bank accounts, phone numbers), emails, and people's names.
- Delete rows older than 30 days on a schedule.
- Ensure INFO-level application logs do not contain raw message content.

## Acceptance criteria
- [ ] Every parser call writes one log row.
- [ ] Unit tests for redacting phone numbers, bank account numbers, and emails.
- [ ] No INFO-level log contains message text.

## Dependencies
None.
