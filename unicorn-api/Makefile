# Makefile for Unicorn API

# Variables
BINARY_NAME=unicorn-api
BUILD_DIR=build
MAIN_FILE=cmd/main.go

# Go related variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_UNIX=$(BINARY_NAME)_unix

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(shell git describe --tags --always --dirty) -X main.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S')"

# Default target
all: clean build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	swag init
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)

# Build for Linux
build-linux:
	@echo "Building $(BINARY_NAME) for Linux..."
	swag init
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_UNIX) $(MAIN_FILE)

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	@swag init -g cmd/main.go
	$(GOCMD) run $(MAIN_FILE)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Generate API documentation
docs:
	@echo "Generating API documentation..."
		swag init -g $(MAIN_FILE); \
	else \
		echo "swag not found. Installing..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		swag init -g $(MAIN_FILE); \
	fi
# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest

# Create build directory
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)
