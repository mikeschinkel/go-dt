.PHONY: help test test-unit test-corpus test-all lint build clean fmt vet tidy examples

LINTER = "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.2"

# Default target
help:
	@echo "Available targets:"
	@echo "  make test         - Run unit tests (main + all sub-modules)"
	@echo "  make test-corpus  - Run fuzz corpus regression tests (main + all sub-modules)"
	@echo "  make test-all     - Run all tests (unit + corpus)"
	@echo "  make lint         - Run golangci-lint"
	@echo "  make fmt          - Format code with gofmt"
	@echo "  make vet          - Run go vet"
	@echo "  make tidy         - Run go mod tidy (main + all sub-modules)"
	@echo "  make build        - Build the package"
	@echo "  make examples     - Build all examples to ./bin/"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make ci           - Run all CI checks (fmt, vet, lint, test-all)"

# Go environment (requires jsonv2 experiment)
GOEXPERIMENT ?= jsonv2
GO := GOEXPERIMENT=$(GOEXPERIMENT) go

ensure-valid: tidy test lint vet examples

# Run unit tests
test: test-unit

test-unit:
	@echo "Running tests for main package..."
	@cd test && $(GO) test -v -race -coverprofile=./coverage.txt -covermode=atomic ./... || exit 1
	@echo ""
	@set -e; for dir in */test; do \
		if [ -d "$$dir" ] && [ -f "$$dir/go.mod" ]; then \
			module=$$(dirname $$dir); \
			echo "Running tests for $$module sub-module..."; \
			cd $$dir && $(GO) test -v -race ./... || exit 1; \
			cd ../..; \
			echo ""; \
		fi \
	done

# Run fuzz corpus regression tests
test-corpus:
	@echo "Running fuzz corpus tests for main package..."
	@cd test && $(GO) test -v -run=TestFuzzCorpus || exit 1
	@echo ""
	@set -e; for dir in */test; do \
		if [ -d "$$dir" ] && [ -f "$$dir/go.mod" ]; then \
			module=$$(dirname $$dir); \
			echo "Running fuzz corpus tests for $$module sub-module..."; \
			cd $$dir && $(GO) test -v -run=TestFuzzCorpus || exit 1; \
			cd ../..; \
			echo ""; \
		fi \
	done

# Run all tests
test-all: test-unit test-corpus

# Run linter
lint:
	@echo "Running linter for main package..."
	@go run $(LINTER) run ./... --timeout=5m || exit 1
	@set -e; for dir in */; do \
		module=$${dir%/}; \
		if [ -d "$$module" ] && [ -f "$$module/go.mod" ]; then \
			echo "Running linter for $$module/..."; \
			cd $$module && go run $(LINTER) run ./... --timeout=5m || exit 1; \
			cd ..; \
		fi \
	done

# Format code
fmt:
	gofmt -s -w .

# Run go vet
vet:
	$(GO) vet ./...

# Run go mod tidy
tidy:
	@echo "Running go mod tidy for main package..."
	@$(GO) mod tidy || exit 1
	@echo "Running go mod tidy for main test..."
	@cd test && $(GO) mod tidy || exit 1
	@set -e; for dir in */; do \
		module=$${dir%/}; \
		if [ -d "$$module" ] && [ -f "$$module/go.mod" ]; then \
			echo "Running go mod tidy for $$module/..."; \
			cd $$module && $(GO) mod tidy || exit 1; \
			cd ..; \
			if [ -d "$$module/test" ] && [ -f "$$module/test/go.mod" ]; then \
				echo "Running go mod tidy for $$module/test/..."; \
				cd $$module/test && $(GO) mod tidy || exit 1; \
				cd ../..; \
			fi \
		fi \
	done

# Build the package
build:
	$(GO) build ./...

# Build all examples to ./bin/
examples:
	@mkdir -p bin
	@for example in examples/*/; do \
		name=$$(basename $$example); \
		echo "Building $$name..."; \
		cd $$example && \
		go mod init example 2>/dev/null || true && \
		go mod edit -replace github.com/mikeschinkel/go-dt=../.. && \
		go mod tidy && \
		go build -o ../../bin/$$name . && \
		cd ../..; \
	done
	@echo "Examples built to ./bin/"

# Clean build artifacts
clean:
	$(GO) clean
	rm -f coverage.txt
	rm -rf bin
	cd test && $(GO) clean

# Run all CI checks locally
ci: fmt vet lint test-all
	@echo "All CI checks passed!"
