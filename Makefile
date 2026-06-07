# wait4it Makefile

# Go parameters
BINARY_NAME  := wait4it
GO          := go
MAIN_PKG    := .
BUILD_FLAGS := -ldflags="-s -w"
GOPATH      := $(shell $(GO) env GOPATH)
GOLINT      := $(GOPATH)/bin/golangci-lint

# Test parameters
INTEGRATION_TAG := integration
COVERAGE_DIR    := coverage
COVERAGE_OUT    := $(COVERAGE_DIR)/coverage.out

# Docker parameters
DOCKER_IMAGE := wait4it

.PHONY: all build run clean test test-unit test-integration test-all lint coverage docker-build docker-build-alpine help

## all: Build the binary (default target)
all: build

## build: Build the wait4it binary
build:
	$(GO) build $(BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_PKG)

## run: Run the binary with default flags
run: build
	./$(BINARY_NAME)

## clean: Remove built binary and coverage artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf $(COVERAGE_DIR)

## test: Run unit tests only (alias for test-unit)
test: test-unit

## test-unit: Run unit tests (no integration tag, no external services)
test-unit:
	$(GO) test -v -race -count=1 ./...

## test-integration: Run integration tests (requires Docker for testcontainers)
test-integration:
	$(GO) test -v -race -count=1 -tags=$(INTEGRATION_TAG) ./...

## test-all: Run both unit and integration tests
test-all: test-unit test-integration

## lint: Run golangci-lint
lint:
	@command -v $(GOLINT) >/dev/null 2>&1 || { echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1; }
	$(GOLINT) run ./...

## coverage: Generate test coverage report (unit tests only)
coverage:
	mkdir -p $(COVERAGE_DIR)
	$(GO) test -race -coverprofile=$(COVERAGE_OUT) ./...
	$(GO) tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report: $(COVERAGE_DIR)/coverage.html"

## docker-build: Build Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE) .

## docker-build-alpine: Build Alpine Docker image
docker-build-alpine:
	docker build -f Dockerfile.alpine -t $(DOCKER_IMAGE):alpine .

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ": "}; {printf "  %-22s %s\n", $$1, $$2}' | sed 's/^## //'
