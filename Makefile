.PHONY: all build test lint clean install run help

# Binary name
BINARY_NAME=erd-viewer

# Build directory
BUILD_DIR=bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOINSTALL=$(GOCMD) install

# Build flags
LDFLAGS=-ldflags "-s -w"

all: test build

## build: Build the binary
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/erd-viewer
	@echo "Binary built at $(BUILD_DIR)/$(BINARY_NAME)"

## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race ./...

## test-integration: Run integration tests (requires TEST_DATABASE_URL)
test-integration:
	@echo "Running integration tests..."
	@if [ -z "$(TEST_DATABASE_URL)" ]; then \
		echo "Error: TEST_DATABASE_URL environment variable not set"; \
		echo "Example: export TEST_DATABASE_URL=postgres://user:pass@localhost:5432/testdb"; \
		exit 1; \
	fi
	$(GOTEST) -v -race ./...

## lint: Run linters
lint:
	@echo "Running linters..."
	@which golangci-lint > /dev/null 2>&1 || (echo "golangci-lint not installed. Run: brew install golangci-lint" && exit 1)
	golangci-lint run --timeout=5m

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.txt

## install: Install the binary to GOPATH/bin
install:
	@echo "Installing..."
	$(GOINSTALL) ./cmd/erd-viewer

## run: Build and run the binary
run: build
	@echo "Running..."
	$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

## coverage: Generate coverage report
coverage:
	@echo "Generating coverage report..."
	$(GOTEST) -coverprofile=coverage.txt ./...
	$(GOCMD) tool cover -html=coverage.txt

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
