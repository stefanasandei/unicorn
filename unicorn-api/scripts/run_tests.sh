#!/bin/bash

# Test runner script for the Unicorn API IAM service
# This script runs all tests with coverage and proper output formatting

set -e

echo "🧪 Running Unicorn API IAM Service Tests"
echo "========================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    print_error "go.mod not found. Please run this script from the unicorn-api directory."
    exit 1
fi

# Clean up any existing test databases
print_status "Cleaning up existing test databases..."
find . -name "test_*.db" -delete 2>/dev/null || true
find . -name "integration_test_*.db" -delete 2>/dev/null || true

# Run unit tests for stores
print_status "Running store unit tests..."
if go test -v ./internal/stores/... -coverprofile=coverage_stores.out; then
    print_success "Store tests passed"
else
    print_error "Store tests failed"
    exit 1
fi

# Run integration tests
print_status "Running integration tests..."
if go test -v ./internal/integration/... -coverprofile=coverage_integration.out; then
    print_success "Integration tests passed"
else
    print_error "Integration tests failed"
    exit 1
fi

# Generate coverage report
print_status "Generating coverage report..."
go tool cover -html=coverage.out -o coverage.html

# Show coverage summary
print_status "Coverage summary:"
go tool cover -func=coverage.out | tail -1

# Clean up coverage files
print_status "Cleaning up coverage files..."
rm -f coverage_stores.out coverage_handlers.out coverage_integration.out

print_success "Test run completed successfully!"
print_status "Coverage report saved to: coverage.html"

# Optional: Open coverage report in browser (macOS)
if command -v open >/dev/null 2>&1; then
    read -p "Open coverage report in browser? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        open coverage.html
    fi
fi 