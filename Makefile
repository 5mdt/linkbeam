.PHONY: build test run clean help

# Build the linkbeam binary
build:
	@mkdir -p dist
	go build -o dist/linkbeam cmd/linkbeam/main.go
	@echo "Binary built at: dist/linkbeam"

# Run all tests
test:
	go test -v ./...

# Run with default config
run:
	go run cmd/linkbeam/main.go

# Run with custom config
run-example:
	go run cmd/linkbeam/main.go --config internal/config/testdata/config.yaml

# Clean build artifacts
clean:
	rm -rf dist/
	rm -f linkbeam

# Run tests with coverage
test-coverage:
	go test -v -cover ./...

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Run linter
lint:
	golangci-lint run

# Install pre-commit hooks
install-hooks:
	pre-commit install

# Run pre-commit checks on all files
pre-commit:
	pre-commit run --all-files

# Show help
help:
	@echo "LinkBeam Makefile Commands:"
	@echo "  make build          - Build the linkbeam binary"
	@echo "  make test           - Run all tests"
	@echo "  make run            - Run with default config.yaml"
	@echo "  make run-example    - Run with example config"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make fmt            - Format code"
	@echo "  make vet            - Run go vet"
	@echo "  make lint           - Run golangci-lint"
	@echo "  make install-hooks  - Install pre-commit hooks"
	@echo "  make pre-commit     - Run pre-commit on all files"
	@echo "  make help           - Show this help message"
