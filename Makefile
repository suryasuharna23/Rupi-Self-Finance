.PHONY: db run test e2e cover eval lint build

db:
	docker compose up -d postgres

run:
	set -a && . ./.env && set +a && go run ./cmd/rupi

test:
	go test ./...

# Needs a running Postgres (make db). Telegram is replaced by a fake server.
e2e:
	RUPI_TEST_DATABASE_URL=postgres://rupi:rupi@localhost:5432/rupi?sslmode=disable go test -tags e2e -count=1 -v ./internal/e2e/

cover:
	go test -coverprofile=coverage.out ./internal/parser/... && go tool cover -func=coverage.out

eval:
	go run ./cmd/eval -v

lint:
	gofmt -l . && go vet ./...

build:
	go build -o bin/rupi ./cmd/rupi
