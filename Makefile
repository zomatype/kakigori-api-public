APP_NAME=kakigori-api
PORT?=8080

.PHONY: build run test lint fmt

build:
	go build -o $(APP_NAME) ./cmd/server/main.go

run:
	KAKIGORI_PORT=$(PORT) go run ./cmd/server/main.go

test:
	go test ./...

lint:
	golangci-lint run || true

fmt:
	gofmt -w .
