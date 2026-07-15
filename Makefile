.PHONY: help lint test build run clean

help:
	@echo "CYTAXI - Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make lint       Run linters"
	@echo "  make test       Run tests"
	@echo "  make build      Build backend"
	@echo "  make run        Run backend"
	@echo "  make clean      Clean build artifacts"

lint:
	cd backend && golangci-lint run ./...

test:
	cd backend && go test ./... -v -race -cover

build:
	cd backend && go build -o ../bin/cytaxi ./cmd

run:
	cd backend && go run ./cmd

clean:
	rm -rf bin/
