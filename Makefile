# GoPray Makefile
.PHONY: build test clean install run-tests coverage fmt vet lint help

# Build configuration
APP_NAME = gopray

# Default target
all: build

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	go build  -o $(APP_NAME)

# Install the application
install:
	@echo "Installing $(APP_NAME)..."
	go install $(GOFLAGS) $(LDFLAGS)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	go test -race -v ./...

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Generate test coverage
coverage:
	@echo "Generating test coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Run staticcheck (requires: go install honnef.co/go/tools/cmd/staticcheck@latest)
lint:
	@echo "Running staticcheck..."
	@if command -v staticcheck > /dev/null; then \
		staticcheck ./...; \
	else \
		echo "staticcheck not found. Install with: go install honnef.co/go/tools/cmd/staticcheck@latest"; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(APP_NAME)
	rm -f coverage.out coverage.html
	go clean -testcache
	go clean -modcache

# Update dependencies
deps:
	@echo "Updating dependencies..."
	go mod tidy
	go mod download

# Check for security vulnerabilities (requires: go install golang.org/x/vuln/cmd/govulncheck@latest)
security:
	@echo "Checking for security vulnerabilities..."
	@if command -v govulncheck > /dev/null; then \
		govulncheck ./...; \
	else \
		echo "govulncheck not found. Install with: go install golang.org/x/vuln/cmd/govulncheck@latest"; \
	fi

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o build/$(APP_NAME)-linux-amd64
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o build/$(APP_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) $(LDFLAGS) -o build/$(APP_NAME)-darwin-arm64
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o build/$(APP_NAME)-windows-amd64.exe

# Run the application (development)
run: build
	./$(APP_NAME)

# Run with arguments
run-list: build
	./$(APP_NAME) list

run-next: build
	./$(APP_NAME) next

run-date: build
	./$(APP_NAME) date

run-config: build
	./$(APP_NAME) config

# Development workflow
dev: fmt vet test build

# Release workflow
release: clean fmt vet test build-all

# Help
help:
	@echo "GoPray Build Commands:"
	@echo ""
	@echo "  build      - Build the application"
	@echo "  install    - Install the application"
	@echo "  test       - Run tests"
	@echo "  test-race  - Run tests with race detection"
	@echo "  bench      - Run benchmarks"
	@echo "  coverage   - Generate test coverage report"
	@echo "  fmt        - Format code"
	@echo "  vet        - Run go vet"
	@echo "  lint       - Run staticcheck linter"
	@echo "  clean      - Clean build artifacts"
	@echo "  deps       - Update dependencies"
	@echo "  security   - Check for security vulnerabilities"
	@echo "  build-all  - Build for multiple platforms"
	@echo "  run        - Build and run the application"
	@echo "  run-*      - Build and run with specific commands"
	@echo "  dev        - Development workflow (format, vet, test, build)"
	@echo "  release    - Release workflow (clean, format, vet, test, build-all)"
	@echo "  help       - Show this help message"
