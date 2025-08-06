# User Service Makefile

# Variables
BINARY_NAME=user-service
MAIN_PATH=cmd/server/main.go
BUILD_DIR=bin
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Go related variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Build info
VERSION?=1.0.0
BUILD_TIME=$(shell date +%Y-%m-%d\ %H:%M:%S)
GIT_COMMIT=$(shell git rev-parse --short HEAD)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

.PHONY: help build clean test test-unit test-integration coverage run deps fmt vet lint docker-build docker-run

# Default target
all: clean deps fmt vet test build

# Help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Dependencies
deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) tidy

# Build
build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

build-linux: ## Build for Linux
	@echo "Building $(BINARY_NAME) for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_PATH)

build-windows: ## Build for Windows
	@echo "Building $(BINARY_NAME) for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows.exe $(MAIN_PATH)

build-darwin: ## Build for macOS
	@echo "Building $(BINARY_NAME) for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin $(MAIN_PATH)

build-all: build-linux build-windows build-darwin ## Build for all platforms

# Clean
clean: ## Remove build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)

# Testing
test: ## Run all tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

test-unit: ## Run unit tests
	@echo "Running unit tests..."
	$(GOTEST) -v ./tests/unit/...

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	$(GOTEST) -v ./tests/integration/...

test-short: ## Run tests with short flag
	@echo "Running short tests..."
	$(GOTEST) -short -v ./...

# Coverage
coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=$(COVERAGE_FILE) ./...
	$(GOCMD) tool cover -func=$(COVERAGE_FILE)

coverage-html: coverage ## Generate HTML coverage report
	@echo "Generating HTML coverage report..."
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

# Code quality
fmt: ## Format code
	@echo "Formatting code..."
	$(GOFMT) ./...

vet: ## Vet code
	@echo "Vetting code..."
	$(GOVET) ./...

lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin v1.50.1"; \
	fi

# Run
run: ## Run the application
	@echo "Running $(BINARY_NAME)..."
	$(GOCMD) run $(MAIN_PATH)

run-build: build ## Build and run the application
	@echo "Running built binary..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Development
dev: ## Run in development mode with auto-reload (requires air)
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Running without auto-reload..."; \
		$(MAKE) run; \
	fi

# Database
migrate-up: ## Run database migrations
	@echo "Running database migrations..."
	# Add your migration command here
	@echo "Migration completed"

migrate-down: ## Rollback database migrations
	@echo "Rolling back database migrations..."
	# Add your rollback command here
	@echo "Rollback completed"

# Docker
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) .
	docker tag $(BINARY_NAME):$(VERSION) $(BINARY_NAME):latest

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8081:8081 --name $(BINARY_NAME) $(BINARY_NAME):latest

docker-stop: ## Stop Docker container
	@echo "Stopping Docker container..."
	docker stop $(BINARY_NAME) || true
	docker rm $(BINARY_NAME) || true

docker-clean: docker-stop ## Clean Docker images
	@echo "Cleaning Docker images..."
	docker rmi $(BINARY_NAME):$(VERSION) || true
	docker rmi $(BINARY_NAME):latest || true

# Documentation
docs: ## Generate swagger documentation
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g $(MAIN_PATH) -o docs/; \
		echo "Swagger documentation generated"; \
	else \
		echo "swag not installed. Install with: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# Installation
install: build ## Install the binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

# Security
security: ## Run security checks
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Performance
bench: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Tools installation
install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.50.1

# Environment setup
setup: deps install-tools ## Setup development environment
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo ".env file created from .env.example"; \
	fi

# Quick commands
quick-test: fmt vet test-short ## Quick test (format, vet, short tests)
quick-build: clean deps quick-test build ## Quick build (clean, deps, quick test, build)

# Release
release: clean deps fmt vet test build-all ## Prepare release (clean, deps, format, vet, test, build all platforms)
