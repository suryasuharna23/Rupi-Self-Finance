.PHONY: db run test cover eval lint build

db:
	docker compose up -d postgres

run:
	set -a && . ./.env && set +a && go run ./cmd/rupi

test:
	go test ./...

cover:
	go test -coverprofile=coverage.out ./internal/parser/... && go tool cover -func=coverage.out

eval:
	go run ./cmd/eval -v

lint:
	gofmt -l . && go vet ./...

build:
	go build -o bin/rupi ./cmd/rupi
