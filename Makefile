.PHONY: help build test test-coverage lint fmt vet tidy clean install run-current run-forecast run-daily

# Binary name
BINARY_NAME=sky
BUILD_DIR=.

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOFMT=gofmt
GOMOD=$(GOCMD) mod

# Build parameters
MAIN_PATH=./cmd/sky
BUILD_FLAGS=-v
LDFLAGS=-s -w

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	$(GOBUILD) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

build-all: tidy ## Build with dependency cleanup
	$(GOBUILD) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

test: ## Run tests
	$(GOTEST) -v ./...

test-coverage: ## Run tests with coverage
	$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-short: ## Run tests without verbose output
	$(GOTEST) ./...

lint: fmt vet ## Run all linting checks

fmt: ## Format code
	$(GOFMT) -w .
	@echo "Code formatted"

fmt-check: ## Check if code is formatted
	@if [ "$$($(GOFMT) -s -l . | wc -l)" -gt 0 ]; then \
		echo "The following files are not formatted:"; \
		$(GOFMT) -s -l .; \
		exit 1; \
	fi

vet: ## Run go vet
	$(GOVET) ./...

tidy: ## Tidy and verify dependencies
	$(GOMOD) tidy
	$(GOMOD) verify
	@echo "Dependencies tidied and verified"

clean: ## Remove build artifacts
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	rm -f coverage.txt coverage.html
	@echo "Build artifacts removed"

install: build ## Install binary to GOPATH/bin
	$(GOCMD) install $(MAIN_PATH)

run-current: build ## Build and run current weather command
	./$(BINARY_NAME) current

run-forecast: build ## Build and run forecast command
	./$(BINARY_NAME) forecast --hours 24

run-daily: build ## Build and run daily forecast command
	./$(BINARY_NAME) daily --days 7

dev: ## Run without building (useful for development)
	$(GOCMD) run $(MAIN_PATH) current

release-snapshot: ## Create a snapshot release with GoReleaser
	goreleaser release --snapshot --clean

release-check: ## Check GoReleaser configuration
	goreleaser check

all: tidy lint test build ## Run tidy, lint, test, and build
