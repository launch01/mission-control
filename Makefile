.PHONY: build test lint clean install help

# Binary name
BINARY_NAME=mission-control
BINARY_PATH=./cmd/mission-control

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	$(GOBUILD) -o $(BINARY_NAME) $(BINARY_PATH)

test: ## Run tests
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	$(GOCMD) tool cover -html=coverage.out

lint: ## Run linter (requires golangci-lint)
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Run: brew install golangci-lint" && exit 1)
	golangci-lint run ./...

clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out

deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) tidy

install: build ## Install the binary
	cp $(BINARY_NAME) $(GOPATH)/bin/

run: build ## Build and run
	./$(BINARY_NAME)

.DEFAULT_GOAL := help
