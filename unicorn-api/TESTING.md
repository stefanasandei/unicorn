# Testing Guide for Unicorn API IAM Service

This document provides comprehensive information about testing the Unicorn API IAM service, including unit tests, integration tests, and full workflow tests.

## Overview

The IAM service includes three types of tests:

1. **Unit Tests** - Test individual components in isolation
2. **Integration Tests** - Test complete workflows and API endpoints
3. **Full Workflow Tests** - Test end-to-end scenarios

## Test Structure

```
unicorn-api/
├── internal/
│   ├── stores/
│   │   ├── iam_sqlite.go
│   │   └── iam_sqlite_test.go      # Unit tests for database operations
│   ├── handlers/
│   │   ├── iam.go
│   │   └── iam_test.go             # Unit tests for API handlers
│   └── integration/
│       └── iam_workflow_test.go    # Integration and workflow tests
├── scripts/
│   └── run_tests.sh                # Test runner script
└── TESTING.md                      # This file
```

## Running Tests

### Quick Start

```bash
# Run all tests with comprehensive output
make test-all

# Or use the script directly
./scripts/run_tests.sh
```

### Individual Test Commands

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run only unit tests
make test-unit

# Run only integration tests
make test-integration

# Run specific test packages
go test -v ./internal/stores/...
go test -v ./internal/handlers/...
go test -v ./internal/integration/...
```

### Test Coverage

The test suite provides comprehensive coverage of:

- **Database Operations**: CRUD operations for roles, organizations, and accounts
- **API Endpoints**: All HTTP handlers with various input scenarios
- **Authentication**: Login, token validation, and refresh flows
- **Authorization**: Role assignment and permission management
- **Error Handling**: Invalid inputs, missing data, and edge cases

## Test Categories

### 1. Unit Tests (`internal/stores/iam_sqlite_test.go`)

Tests individual database operations:

- **Role Management**

  - Creating roles with different permission combinations
  - Retrieving roles by name
  - Handling duplicate role names

- **Organization Management**

  - Creating organizations
  - Retrieving organizations by name
  - Handling duplicate organization names

- **Account Management**

  - Creating user and bot accounts
  - Retrieving accounts by email
  - Updating account information
  - Assigning roles to accounts

- **Error Scenarios**
  - Invalid UUID formats
  - Non-existent records
  - Database constraint violations

### 2. Handler Unit Tests (`internal/handlers/iam_test.go`)

Tests API endpoints with mock data:

- **Role Endpoints**

  - `POST /api/v1/roles` - Creating roles
  - `POST /api/v1/roles/assign` - Assigning roles

- **Organization Endpoints**

  - `POST /api/v1/organizations` - Creating organizations

- **User Endpoints**

  - `POST /api/v1/organizations/{org_id}/users` - Creating users

- **Authentication Endpoints**

  - `POST /api/v1/login` - User login
  - `GET /api/v1/token/validate` - Token validation
  - `POST /api/v1/token/refresh` - Token refresh

- **Validation Tests**
  - Required field validation
  - Email format validation
  - Password strength requirements
  - Invalid UUID formats

### 3. Integration Tests (`internal/integration/iam_workflow_test.go`)

Tests complete workflows using real database:

#### Full IAM Workflow Test

Tests the complete lifecycle:

1. **Organization Creation**

   - Create an organization
   - Verify organization properties

2. **Role Creation**

   - Create admin role (Read, Write, Delete permissions)
   - Create user role (Read, Write permissions)
   - Create read-only role (Read permission only)

3. **User Creation**

   - Create admin user in the organization
   - Create regular user in the organization
   - Create read-only user in the organization

4. **Role Assignment**

   - Assign roles to users
   - Change user roles dynamically

5. **Authentication Flow**

   - Login with valid credentials
   - Generate JWT tokens
   - Validate tokens
   - Test invalid login scenarios

6. **Data Persistence**
   - Verify data persists across operations
   - Test data integrity

#### Organization Isolation Test

Tests multi-tenant isolation:

- Create multiple organizations
- Test user isolation between organizations
- Verify email uniqueness constraints

#### Role Permission Workflow Test

Tests various permission combinations:

- Read-only permissions
- Read-write permissions
- Full access permissions
- Write-only permissions
- Delete-only permissions

## Test Data

### Sample Organization Structure

```
Acme Corporation
├── Admin User (admin@acme.com)
│   └── Role: admin (Read, Write, Delete)
├── Regular User (user@acme.com)
│   └── Role: user (Read, Write)
└── Read-Only User (readonly@acme.com)
    └── Role: readonly (Read only)
```

### Permission Levels

- **0 (Read)**: Can view data
- **1 (Write)**: Can create and modify data
- **2 (Delete)**: Can delete data

## Test Environment

### Database

- Tests use temporary SQLite databases
- Each test creates its own database file
- Databases are automatically cleaned up after tests

### Configuration

- Test-specific JWT secret keys
- 24-hour token expiration for testing
- Test environment mode

## Coverage Reports

After running tests with coverage, you'll get:

- **HTML Coverage Report**: `coverage.html`
- **Console Coverage Summary**: Shows overall coverage percentage
- **Function-level Coverage**: Detailed breakdown by function

### Viewing Coverage

```bash
# Generate and view coverage report
make test-coverage

# Open HTML report in browser (macOS)
open coverage.html
```

## Best Practices

### Writing New Tests

1. **Use descriptive test names** that explain what is being tested
2. **Test both success and failure scenarios**
3. **Use table-driven tests** for multiple similar test cases
4. **Clean up test data** after each test
5. **Mock external dependencies** in unit tests
6. **Use real database** in integration tests

### Test Organization

```go
func TestFeatureName(t *testing.T) {
    // Setup
    setupTestData(t)
    defer cleanupTestData(t)

    // Test cases
    tests := []struct {
        name           string
        input          string
        expectedResult string
        expectError    bool
    }{
        // Test cases here
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Troubleshooting

### Common Issues

1. **Database Lock Errors**

   - Ensure no other processes are using the test database
   - Check for proper cleanup in test teardown

2. **Permission Errors**

   - Ensure the script has execute permissions: `chmod +x scripts/run_tests.sh`

3. **Import Errors**

   - Run `go mod tidy` to ensure dependencies are up to date
   - Check that all required packages are imported

4. **Test Failures**
   - Check test database cleanup
   - Verify test data isolation
   - Review error messages for specific failure reasons

### Debug Mode

Run tests with verbose output for debugging:

```bash
go test -v -count=1 ./internal/integration/...
```

The `-count=1` flag ensures tests run without caching.

## Continuous Integration

The test suite is designed to run in CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
- name: Run Tests
  run: |
    cd unicorn-api
    make test-all
```

## Performance Considerations

- Unit tests run quickly with mocked dependencies
- Integration tests use real database but are isolated
- Test databases are small and temporary
- Cleanup ensures no resource leaks

## Contributing

When adding new features:

1. **Add unit tests** for new functionality
2. **Add integration tests** for new workflows
3. **Update this documentation** if needed
4. **Ensure all tests pass** before submitting

## Support

For test-related issues:

1. Check the troubleshooting section above
2. Review test logs for specific error messages
3. Ensure your Go environment is properly configured
4. Verify all dependencies are installed
