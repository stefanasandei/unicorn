# Unicorn API

A comprehensive RESTful API for Unicorn services providing Identity and Access Management (IAM) functionality, storage, compute, and secrets management.

## Features

- **IAM (Identity and Access Management)**: User authentication, role-based access control, organization management
- **Storage**: File storage with bucket management
- **Compute**: Docker container management with resource limits
- **Secrets Manager**: Encrypted secret storage per user
- **Lambda**: Function execution service

## Quick Start

### Prerequisites

- Go 1.21+
- Docker (for compute functionality)
- SQLite (embedded)

### Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd unicorn-api
```

2. Install dependencies:

```bash
go mod download
```

3. Create environment file:

```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run the application:

```bash
go run cmd/main.go
```

The API will be available at `http://localhost:8080`

## Configuration

### Environment Variables

| Variable           | Description                          | Default                       |
| ------------------ | ------------------------------------ | ----------------------------- |
| `DB_PATH`          | SQLite database path                 | `unicorn.db`                  |
| `STORAGE_PATH`     | File storage directory               | `./storage`                   |
| `PORT`             | Server port                          | `8080`                        |
| `ENVIRONMENT`      | Environment (development/production) | `development`                 |
| `JWT_SECRET`       | JWT signing secret                   | Random string                 |
| `JWT_EXPIRY_HOURS` | JWT token expiry                     | `24`                          |
| `DOCKER_HOST`      | Docker daemon socket                 | `unix:///var/run/docker.sock` |
| `LAMBDA_URL`       | Lambda API URL                       | `http://localhost:8081`       |

### Example .env file

```env
# Database Configuration
DB_PATH=unicorn.db

# Storage Configuration
STORAGE_PATH=./storage

# Server Configuration
PORT=8080
ENVIRONMENT=development

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRY_HOURS=24

# Docker Configuration (for compute service)
DOCKER_HOST=unix:///var/run/docker.sock

# Lambda API Configuration
LAMBDA_URL=http://localhost:8081

# Logging Configuration
LOG_LEVEL=info
```

## API Documentation

Once the server is running, you can access the Swagger UI at:
`http://localhost:8080/swagger/index.html`

## Authentication

The API uses JWT tokens for authentication. To get a token:

1. Create an organization and user (or use the default admin)
2. Login with your credentials
3. Use the returned token in the `Authorization` header: `Bearer <token>`

### Default Admin Account

The system creates a default admin account:

- Email: `admin@unicorn.local`
- Password: `admin123`

## Troubleshooting

### Common Issues

#### 1. Secret Manager JSON Errors

**Problem**: Getting "invalid character 'y' looking for beginning of value" error

**Solution**:

- Ensure metadata is valid JSON format
- Check that the secret value doesn't contain invalid characters
- Verify the request body is properly formatted

#### 2. Docker Image Pull Failures

**Problem**: Compute service can't pull Docker images

**Solution**:

- Ensure Docker daemon is running: `sudo systemctl start docker`
- Check Docker permissions: `sudo usermod -aG docker $USER`
- Verify network connectivity for image registry
- Check Docker socket permissions

#### 3. Database Errors

**Problem**: Database initialization fails

**Solution**:

- Ensure the database directory is writable
- Check disk space
- Verify SQLite is available

#### 4. Permission Denied Errors

**Problem**: Getting permission errors for various operations

**Solution**:

- Ensure user has the correct role assigned
- Check role permissions (0=Read, 1=Write, 2=Delete)
- Verify JWT token is valid and not expired

### Debug Mode

To enable debug logging, set the environment variable:

```bash
export LOG_LEVEL=debug
```

### Health Check

Check if the service is running:

```bash
curl http://localhost:8080/health
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific test
go test ./internal/handlers

# Run with coverage
go test -cover ./...
```

### Code Structure

```
unicorn-api/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── auth/                # Authentication utilities
│   ├── common/
│   │   ├── errors/          # Error handling
│   │   └── validation/      # Input validation
│   ├── config/              # Configuration management
│   ├── handlers/            # HTTP request handlers
│   ├── middleware/          # HTTP middleware
│   ├── models/              # Data models
│   ├── routes/              # Route definitions
│   ├── services/            # Business logic layer
│   └── stores/              # Data access layer
├── docs/                    # Generated API documentation
└── scripts/                 # Utility scripts
```

### Architecture

The application follows a clean architecture pattern:

1. **Handlers**: HTTP request/response handling
2. **Services**: Business logic and orchestration
3. **Stores**: Data access and persistence
4. **Models**: Data structures and validation

## Security

### Best Practices

1. **Change Default Credentials**: Update the default admin password
2. **Use Strong JWT Secrets**: Generate a strong random secret
3. **Enable HTTPS**: Use TLS in production
4. **Regular Updates**: Keep dependencies updated
5. **Input Validation**: All inputs are validated
6. **Encryption**: Secrets are encrypted at rest

### Production Deployment

1. Set `ENVIRONMENT=production`
2. Use a strong `JWT_SECRET`
3. Enable HTTPS/TLS
4. Use a proper database (PostgreSQL/MySQL)
5. Set up monitoring and logging
6. Configure proper firewall rules

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

MIT License - see LICENSE file for details.
