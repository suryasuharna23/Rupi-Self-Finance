# Rupi

A personal finance assistant for university students, used entirely through
Telegram chat. Type "beli bakso 15rb" and Rupi records, categorises, and
reports it.

The full requirements are in [`docs/Rupi - Dokumen Kebutuhan Sistem.pdf`](docs/)
(Indonesian). Work is tracked as tickets in [`docs/tickets/`](docs/tickets/).

## Status: Phase 1, conversation core (monolith)

Following the roadmap, Phase 1 is deliberately a single Go service so the
parser can be refined before the system is split up. The packages in
`internal/` already follow the Phase 2 service boundaries, so each package
later becomes one gRPC service.

| Package                  | Future service (Phase 2) | Contents                                                     |
| ------------------------ | ------------------------ | ------------------------------------------------------------ |
| `internal/telegram`      | Channel Adapter          | Long polling / webhook with `secret_token`                   |
| `internal/conversation`  | Conversation Service     | Onboarding, clarification, corrections, idempotency, rate limit |
| `internal/parser`        | Parser Service           | Rules (amounts, time, categories) + LLM interface            |
| `internal/transaction`   | Transaction Service      | Store, correct + audit trail, undo                           |
| `internal/user`          | User Service             | Profile, time zone, account deletion                         |
| `internal/report`        | Report Service           | Totals per period/category, `/ringkas`                       |
| `internal/storage`       | —                        | PostgreSQL connection and migrations (one schema per service) |

The bot talks to users in Indonesian; code, docs, and tickets are in English.

## Running locally

Prerequisites: Go 1.27+, Docker, and a bot token from [@BotFather](https://t.me/BotFather).

```sh
cp .env.example .env        # fill in TELEGRAM_BOT_TOKEN
make db                     # PostgreSQL 16 via docker compose
make run                    # runs migrations, then starts long polling
```

Without `make`, in PowerShell:

```powershell
docker compose up -d postgres
Get-Content .env | Where-Object { $_ -match "^[A-Z_]+=" } | ForEach-Object {
  $k, $v = $_ -split "=", 2; Set-Item "env:$k" $v
}
go run ./cmd/rupi
```

## Testing

```sh
make test         # unit tests
make integration  # store tests against a throwaway Postgres container (needs Docker)
make e2e          # full chat flow: real Postgres + fake Telegram
make eval         # parser accuracy against testdata/eval/dataset.jsonl
```

On Windows with Docker Desktop, Testcontainers sometimes fails to detect Docker
("rootless Docker is not supported on Windows"). Pointing it at the engine fixes it:
`export DOCKER_HOST=npipe:////./pipe/dockerDesktopLinuxEngine`.

The Phase 1 gate is ≥ 95% extraction accuracy on a 300-sentence evaluation
dataset (NFR-03). The dataset currently has 20 examples taken from the
requirements document; add real sentences to `testdata/eval/dataset.jsonl`.

## Bot commands

| Message                              | What it does                             |
| ------------------------------------ | ---------------------------------------- |
| `/start`                             | Sign up and onboarding                   |
| `beli bakso 15rb`                    | Record a transaction                     |
| `yang bakso tadi 18rb`               | Correct an amount                        |
| `/undo`, `batalin yang terakhir`     | Undo the latest transaction (≤ 24 hours) |
| `bulan ini habis berapa buat makan?` | Ask for a total                          |
| `/ringkas`                           | Summary for today, this week, this month |
| `/hapusakun`                         | Delete all data                          |

## Not done yet

- LLM implementation: the provider is still an open question in the
  requirements. The `parser.LLM` interface and output validation (SEC-05)
  exist; for now, messages the rules don't understand get a clarifying question.
- Opening balance during onboarding (optional in FR-02) and changing time zone.
- Purging soft-deleted accounts after 30 days.
