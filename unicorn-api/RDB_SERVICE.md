# RDB Service Documentation

The RDB (Relational Database) service provides managed database instances for PostgreSQL and MySQL. Users can create, list, and delete database instances with proper volume management and permission controls.

## Features

- **Database Types**: Support for PostgreSQL and MySQL
- **Resource Presets**: Micro, Small, and Medium presets with different CPU and memory allocations
- **Volume Management**: Persistent storage with configurable volumes
- **Permission Controls**: Read, Write, and Delete permissions for different user roles
- **Docker-based**: Uses Docker containers for database instances

## API Endpoints

### Create RDB Instance

```
POST /api/v1/rdb/create
```

**Request Body:**

```json
{
  "name": "my-postgres-db",
  "type": "postgresql",
  "preset": "micro",
  "database": "mydb",
  "username": "user",
  "password": "password123",
  "volumes": [
    {
      "name": "data",
      "size": 1,
      "mount_path": "/var/lib/postgresql/data"
    }
  ],
  "environment": {
    "POSTGRES_INITDB_ARGS": "--encoding=UTF-8"
  }
}
```

**Response:**

```json
{
  "id": "container_id",
  "name": "my-postgres-db",
  "type": "postgresql",
  "status": "running",
  "port": "12345",
  "host": "localhost",
  "database": "mydb",
  "username": "user",
  "volumes": [
    {
      "name": "data",
      "size": 1,
      "mount_path": "/var/lib/postgresql/data"
    }
  ],
  "environment": {
    "POSTGRES_INITDB_ARGS": "--encoding=UTF-8"
  },
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

### List RDB Instances

```
GET /api/v1/rdb/list
```

**Response:**

```json
[
  {
    "id": "container_id",
    "name": "my-postgres-db",
    "type": "postgresql",
    "status": "running",
    "port": "12345",
    "host": "localhost",
    "database": "mydb",
    "username": "user",
    "volumes": [...],
    "environment": {...},
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
]
```

### Delete RDB Instance

```
DELETE /api/v1/rdb/{id}
```

## Database Types

### PostgreSQL

- **Image**: `postgres:15-alpine`
- **Default Port**: 5432
- **Environment Variables**:
  - `POSTGRES_DB`: Database name
  - `POSTGRES_USER`: Username
  - `POSTGRES_PASSWORD`: Password
  - `POSTGRES_HOST_AUTH_METHOD`: Set to "trust" for development

### MySQL

- **Image**: `mysql:8.0`
- **Default Port**: 3306
- **Environment Variables**:
  - `MYSQL_ROOT_PASSWORD`: Root password
  - `MYSQL_DATABASE`: Database name
  - `MYSQL_USER`: Username
  - `MYSQL_PASSWORD`: Password

## Resource Presets

### Micro

- **CPU**: 0.5 cores (500,000,000 nano CPUs)
- **Memory**: 512MB
- **Use Case**: Development, testing, small applications

### Small

- **CPU**: 1 core (1,000,000,000 nano CPUs)
- **Memory**: 1GB
- **Use Case**: Small production applications

### Medium

- **CPU**: 2 cores (2,000,000,000 nano CPUs)
- **Memory**: 2GB
- **Use Case**: Medium production applications

## Volume Management

Volumes provide persistent storage for database data. Each volume has:

- **Name**: Unique identifier for the volume
- **Size**: Size in MB (minimum 1MB, maximum 100GB)

Mount paths are automatically determined based on the database type:

- **PostgreSQL**: `/var/lib/postgresql/data`
- **MySQL**: `/var/lib/mysql`

Example volume configuration:

```json
{
  "volumes": [
    {
      "name": "data",
      "size": 5120
    },
    {
      "name": "backups",
      "size": 10240
    }
  ]
}
```

## Permission System

The RDB service integrates with the IAM system for access control:

- **Read Permission (0)**: Can list RDB instances
- **Write Permission (1)**: Can create RDB instances
- **Delete Permission (2)**: Can delete RDB instances

Users must have the appropriate permissions to perform operations.

## Usage Examples

### Creating a PostgreSQL Instance

```bash
curl -X POST http://localhost:8080/api/v1/rdb/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "my-app-db",
    "type": "postgresql",
    "preset": "small",
    "database": "myapp",
    "username": "appuser",
    "password": "securepass123",
    "volumes": [
      {
        "name": "data",
        "size": 5120
      }
    ]
  }'
```

### Creating a MySQL Instance

```bash
curl -X POST http://localhost:8080/api/v1/rdb/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "my-mysql-db",
    "type": "mysql",
    "preset": "micro",
    "database": "testdb",
    "username": "testuser",
    "password": "testpass123"
  }'
```

### Listing All Instances

```bash
curl -X GET http://localhost:8080/api/v1/rdb/list \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Deleting an Instance

```bash
curl -X DELETE http://localhost:8080/api/v1/rdb/CONTAINER_ID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Connection Information

Once an RDB instance is created, you can connect to it using:

- **Host**: `localhost`
- **Port**: The port returned in the response
- **Database**: The database name specified in the request
- **Username**: The username specified in the request
- **Password**: The password specified in the request

Example connection string for PostgreSQL:

```
postgresql://user:password123@localhost:12345/mydb
```

Example connection string for MySQL:

```
mysql://testuser:testpass123@localhost:12345/testdb
```

## Security Considerations

1. **Password Generation**: If no password is provided, a random 16-character password is generated
2. **Volume Isolation**: Each user's volumes are isolated and prefixed with their user ID
3. **Container Labels**: Containers are labeled with owner and type for proper isolation
4. **Permission Checks**: All operations require appropriate permissions
5. **Resource Limits**: Containers have CPU and memory limits based on the selected preset
6. **Automatic Mount Paths**: Mount paths are automatically set to secure database directories

## Error Handling

The service returns appropriate HTTP status codes:

- `200`: Success
- `400`: Bad Request (invalid parameters)
- `401`: Unauthorized (missing or invalid token)
- `403`: Forbidden (insufficient permissions)
- `404`: Not Found (instance not found)
- `500`: Internal Server Error (Docker or system errors)

## Docker Requirements

The RDB service requires Docker to be running and accessible. On macOS, it uses the Docker Desktop socket path by default.
