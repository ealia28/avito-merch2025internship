.PHONY: build run lint test cover

build:
	go build -o bin/shop ./cmd/shop

run:
	go run ./cmd/shop

lint:
	golangci-lint run ./...

test:
	go test -race -count=1 ./...

cover:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -1
