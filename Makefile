.PHONY: build run lint tidy

build:
	go build ./...

run:
	go run ./cmd/main.go

lint:
	golangci-lint run ./...

tidy:
	go mod tidy
