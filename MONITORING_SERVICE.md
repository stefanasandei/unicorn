# Monitoring Service Documentation

## Overview

The Monitoring Service is a comprehensive solution for tracking resource usage, billing, and performance metrics across all Unicorn API services. It provides real-time monitoring, historical data analysis, and automated billing generation.

## Features

### 🔍 Resource Tracking

- **Real-time Monitoring**: Track CPU, memory, storage, and network usage
- **Resource Lifecycle**: Monitor creation, updates, and deletion of resources
- **Status Tracking**: Active, inactive, and deleted resource states
- **Cost Calculation**: Automatic cost calculation based on usage duration and resource type

### 💰 Billing & Cost Management

- **Hourly Pricing**: Different rates for compute, storage, lambda, RDB, and secrets
- **Monthly Billing**: Automated monthly billing period generation
- **Cost Breakdown**: Detailed cost analysis by resource type
- **Payment Tracking**: Track payment status and transaction history

### 📊 Analytics & Reporting

- **Usage Trends**: Monthly usage patterns and cost trends
- **Resource Summary**: Overview of total resources and active resources
- **Performance Metrics**: CPU, memory, and storage utilization
- **Historical Data**: Long-term usage history for analysis

### 🎯 Dashboard Integration

- **Real-time Dashboard**: Live monitoring interface
- **Interactive Charts**: Visual representation of usage data
- **Resource Management**: View and manage active resources
- **Billing History**: Complete billing period overview

## Architecture

### Data Models

#### ResourceUsage

Tracks current resource usage and metrics:

```go
type ResourceUsage struct {
    ID                uuid.UUID
    AccountID         uuid.UUID
    OrganizationID    uuid.UUID
    ResourceType      ResourceType
    ResourceID        string
    ResourceName      string
    Status            ResourceStatus
    CPUUsage          float64
    MemoryUsage       float64
    StorageUsage      float64
    NetworkUsage      float64
    RequestCount      int64
    ExecutionTime     float64
    CostPerHour       float64
    TotalCost         float64
    Currency          string
    ResourceCreatedAt time.Time
    LastActiveAt      *time.Time
    LastUpdatedAt     time.Time
    Configuration     string
    Tags              string
    Notes             string
}
```

#### ResourceUsageHistory

Historical usage data for billing and analytics:

```go
type ResourceUsageHistory struct {
    ID                 uuid.UUID
    AccountID          uuid.UUID
    OrganizationID     uuid.UUID
    ResourceType       ResourceType
    ResourceID         string
    ResourceName       string
    AvgCPUUsage        float64
    PeakCPUUsage       float64
    AvgMemoryUsage     float64
    PeakMemoryUsage    float64
    TotalStorageUsage  float64
    TotalNetworkUsage  float64
    TotalRequests      int64
    TotalExecutionTime float64
    TotalCost          float64
    Currency           string
    PeriodStart        time.Time
    PeriodEnd          time.Time
    DurationHours      float64
}
```

#### BillingPeriod

Monthly billing periods with cost breakdown:

```go
type BillingPeriod struct {
    ID              uuid.UUID
    OrganizationID  uuid.UUID
    PeriodStart     time.Time
    PeriodEnd       time.Time
    TotalCost       float64
    Currency        string
    IsPaid          bool
    PaidAt          *time.Time
    PaymentMethod   string
    TransactionID   string
    ComputeCost     float64
    LambdaCost      float64
    StorageCost     float64
    RDBCost         float64
    SecretCost      float64
}
```

#### MonitoringMetrics

Real-time monitoring metrics:

```go
type MonitoringMetrics struct {
    ID             uuid.UUID
    OrganizationID uuid.UUID
    ResourceType   ResourceType
    ResourceID     string
    CPUUsage       float64
    MemoryUsage    float64
    StorageUsage   float64
    NetworkUsage   float64
    RequestRate    float64
    ErrorRate      float64
    ResponseTime   float64
    Status         ResourceStatus
    HealthStatus   string
    LastHealthCheck *time.Time
}
```

### Pricing Structure

| Resource Type   | Cost per Hour | Description                    |
| --------------- | ------------- | ------------------------------ |
| Compute (Micro) | $0.05         | Small compute instances        |
| Compute (Small) | $0.10         | Medium compute instances       |
| Lambda          | $0.02         | Serverless function executions |
| Storage         | $0.01/GB      | Per GB per hour                |
| RDB             | $0.15         | Database instances             |
| Secrets         | $0.01         | Secret management              |

## API Endpoints

### Resource Usage

#### GET /api/v1/monitoring/usage

Get resource usage summary for the current organization.

**Query Parameters:**

- `start` (optional): Start date in YYYY-MM-DD format
- `end` (optional): End date in YYYY-MM-DD format

**Response:**

```json
{
  "current_usage": {
    "id": "uuid",
    "resource_type": "compute",
    "resource_name": "web-server",
    "status": "active",
    "cpu_usage": 45.2,
    "memory_usage": 1024,
    "storage_usage": 512,
    "network_usage": 256,
    "cost_per_hour": 0.05,
    "total_cost": 12.5,
    "currency": "USD"
  },
  "summary": {
    "total_resources": 8,
    "active_resources": 6,
    "total_cost": 45.25,
    "currency": "USD",
    "usage_by_type": {
      "compute": 2,
      "storage": 3,
      "lambda": 1,
      "rdb": 1,
      "secret": 1
    },
    "cost_by_type": {
      "compute": 25.5,
      "storage": 8.75,
      "lambda": 5.25,
      "rdb": 3.75,
      "secret": 2.0
    }
  }
}
```

#### GET /api/v1/monitoring/resources/active

Get all active resources for the current organization.

**Response:**

```json
[
  {
    "id": "uuid",
    "resource_type": "compute",
    "resource_name": "web-server-1",
    "status": "active",
    "cpu_usage": 45.2,
    "memory_usage": 1024,
    "storage_usage": 512,
    "network_usage": 256,
    "cost_per_hour": 0.05,
    "total_cost": 12.5,
    "currency": "USD"
  }
]
```

### Monitoring Metrics

#### GET /api/v1/monitoring/metrics/{resource_type}/{resource_id}

Get real-time monitoring metrics for a specific resource.

**Response:**

```json
{
  "id": "uuid",
  "organization_id": "uuid",
  "resource_type": "compute",
  "resource_id": "container-123",
  "cpu_usage": 45.2,
  "memory_usage": 1024,
  "storage_usage": 512,
  "network_usage": 256,
  "request_rate": 150.5,
  "error_rate": 0.1,
  "response_time": 125.3,
  "status": "active",
  "health_status": "healthy",
  "last_health_check": "2024-01-15T10:30:00Z"
}
```

#### PUT /api/v1/monitoring/metrics/{resource_type}/{resource_id}

Update monitoring metrics for a specific resource.

**Request Body:**

```json
{
  "cpu_usage": 45.2,
  "memory_usage": 1024,
  "storage_usage": 512,
  "network_usage": 256,
  "request_rate": 150.5,
  "error_rate": 0.1,
  "response_time": 125.3,
  "status": "active",
  "health_status": "healthy"
}
```

### Billing

#### GET /api/v1/monitoring/billing

Get billing history for the current organization.

**Response:**

```json
[
  {
    "id": "uuid",
    "organization_id": "uuid",
    "period_start": "2024-01-01T00:00:00Z",
    "period_end": "2024-01-31T23:59:59Z",
    "total_cost": 45.25,
    "currency": "USD",
    "is_paid": true,
    "paid_at": "2024-02-01T10:00:00Z",
    "payment_method": "credit_card",
    "transaction_id": "txn_123456",
    "compute_cost": 25.5,
    "lambda_cost": 5.25,
    "storage_cost": 8.75,
    "rdb_cost": 3.75,
    "secret_cost": 2.0
  }
]
```

#### POST /api/v1/monitoring/billing/generate

Generate monthly billing for a specific month.

**Query Parameters:**

- `year` (required): Year (e.g., 2024)
- `month` (required): Month (1-12)

**Response:**

```json
{
  "id": "uuid",
  "organization_id": "uuid",
  "period_start": "2024-01-01T00:00:00Z",
  "period_end": "2024-01-31T23:59:59Z",
  "total_cost": 45.25,
  "currency": "USD",
  "is_paid": false,
  "compute_cost": 25.5,
  "lambda_cost": 5.25,
  "storage_cost": 8.75,
  "rdb_cost": 3.75,
  "secret_cost": 2.0
}
```

### Usage Trends

#### GET /api/v1/monitoring/trends

Get monthly usage trends for the current organization.

**Query Parameters:**

- `months` (optional): Number of months to include (default: 6, max: 24)

**Response:**

```json
[
  {
    "month": "2024-01",
    "total_cost": 45.25,
    "resources": 8
  },
  {
    "month": "2023-12",
    "total_cost": 38.5,
    "resources": 7
  }
]
```

### Resource Tracking

#### POST /api/v1/monitoring/track/create

Track resource creation.

**Request Body:**

```json
{
  "resource_type": "compute",
  "resource_id": "container-123",
  "resource_name": "web-server",
  "status": "active",
  "cpu_usage": 0,
  "memory_usage": 0,
  "storage_usage": 0,
  "network_usage": 0,
  "cost_per_hour": 0.05,
  "currency": "USD",
  "configuration": "{\"image\":\"nginx\",\"preset\":\"micro\"}",
  "tags": "{\"environment\":\"production\"}",
  "notes": "Web server for production"
}
```

#### PUT /api/v1/monitoring/track/{resource_type}/{resource_id}

Track resource updates.

**Request Body:**

```json
{
  "status": "active",
  "cpu_usage": 45.2,
  "memory_usage": 1024,
  "storage_usage": 512,
  "network_usage": 256,
  "request_count": 1500,
  "execution_time": 125.3
}
```

#### DELETE /api/v1/monitoring/track/{resource_type}/{resource_id}

Track resource deletion.

## Integration with Existing Services

### Compute Service Integration

The monitoring service automatically tracks compute resource creation and deletion:

```go
// In compute handler
containerInfo, err := h.service.CreateContainer(userID, req)
if err != nil {
    return err
}

// Track resource creation
if h.monitoringService != nil {
    account, err := h.iamStore.GetAccountByID(claims.AccountID)
    if err == nil {
        configuration := fmt.Sprintf(`{"image":"%s","preset":"%s"}`, req.Image, req.Preset)
        err = h.monitoringService.TrackResourceCreation(
            userID,
            account.OrganizationID,
            models.ResourceTypeCompute,
            containerInfo.ID,
            req.Name,
            configuration,
        )
    }
}
```

### Storage Service Integration

Similar integration can be added to storage, lambda, RDB, and secrets services.

## Dashboard Features

### Overview Tab

- **Resource Summary**: Total and active resource counts
- **Cost Overview**: Total monthly cost and cost breakdown
- **Performance Metrics**: Average CPU usage and storage utilization
- **Cost Breakdown**: Visual representation of costs by resource type

### Active Resources Tab

- **Resource List**: All active resources with real-time metrics
- **Status Indicators**: Visual status badges (active, inactive, deleted)
- **Performance Metrics**: CPU, memory, storage, and network usage
- **Cost Information**: Hourly rate and total cost for each resource

### Billing History Tab

- **Monthly Bills**: Complete billing period overview
- **Payment Status**: Paid/pending status with payment details
- **Cost Breakdown**: Detailed cost analysis by resource type
- **Transaction History**: Payment method and transaction IDs

### Usage Trends Tab

- **Monthly Trends**: Cost and resource usage over time
- **Visual Charts**: Graphical representation of trends
- **Historical Data**: Long-term usage patterns

## Database Schema

The monitoring service uses SQLite with the following tables:

### resource_usages

```sql
CREATE TABLE resource_usages (
    id TEXT PRIMARY KEY,
    created_at DATETIME,
    updated_at DATETIME,
    account_id TEXT NOT NULL,
    organization_id TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id TEXT NOT NULL,
    resource_name TEXT NOT NULL,
    status TEXT NOT NULL,
    cpu_usage REAL DEFAULT 0,
    memory_usage REAL DEFAULT 0,
    storage_usage REAL DEFAULT 0,
    network_usage REAL DEFAULT 0,
    request_count INTEGER DEFAULT 0,
    execution_time REAL DEFAULT 0,
    cost_per_hour REAL DEFAULT 0,
    total_cost REAL DEFAULT 0,
    currency TEXT DEFAULT 'USD',
    resource_created_at DATETIME,
    last_active_at DATETIME,
    last_updated_at DATETIME,
    configuration TEXT,
    tags TEXT,
    notes TEXT
);
```

### resource_usage_histories

```sql
CREATE TABLE resource_usage_histories (
    id TEXT PRIMARY KEY,
    created_at DATETIME,
    account_id TEXT NOT NULL,
    organization_id TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id TEXT NOT NULL,
    resource_name TEXT NOT NULL,
    avg_cpu_usage REAL,
    peak_cpu_usage REAL,
    avg_memory_usage REAL,
    peak_memory_usage REAL,
    total_storage_usage REAL,
    total_network_usage REAL,
    total_requests INTEGER,
    total_execution_time REAL,
    total_cost REAL,
    currency TEXT DEFAULT 'USD',
    period_start DATETIME NOT NULL,
    period_end DATETIME NOT NULL,
    duration_hours REAL
);
```

### billing_periods

```sql
CREATE TABLE billing_periods (
    id TEXT PRIMARY KEY,
    created_at DATETIME,
    updated_at DATETIME,
    organization_id TEXT NOT NULL,
    period_start DATETIME NOT NULL,
    period_end DATETIME NOT NULL,
    total_cost REAL DEFAULT 0,
    currency TEXT DEFAULT 'USD',
    is_paid BOOLEAN DEFAULT FALSE,
    paid_at DATETIME,
    payment_method TEXT,
    transaction_id TEXT,
    compute_cost REAL DEFAULT 0,
    lambda_cost REAL DEFAULT 0,
    storage_cost REAL DEFAULT 0,
    rdb_cost REAL DEFAULT 0,
    secret_cost REAL DEFAULT 0
);
```

### monitoring_metrics

```sql
CREATE TABLE monitoring_metrics (
    id TEXT PRIMARY KEY,
    created_at DATETIME,
    organization_id TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id TEXT NOT NULL,
    cpu_usage REAL,
    memory_usage REAL,
    storage_usage REAL,
    network_usage REAL,
    request_rate REAL,
    error_rate REAL,
    response_time REAL,
    status TEXT,
    health_status TEXT,
    last_health_check DATETIME
);
```

## Security & Permissions

### Authentication

All monitoring endpoints require JWT authentication with valid user tokens.

### Authorization

- **Read Access**: Users can view monitoring data for their organization
- **Write Access**: Users can update metrics and track resource changes
- **Admin Access**: Organization admins can generate billing and manage monitoring

### Data Privacy

- All monitoring data is scoped to the user's organization
- No cross-organization data access
- Historical data is retained for billing and analytics purposes

## Performance Considerations

### Database Optimization

- Indexed queries on organization_id, resource_type, and resource_id
- Efficient aggregation queries for billing calculations
- Regular cleanup of old monitoring metrics

### Caching Strategy

- Cache frequently accessed metrics (5-minute TTL)
- Cache billing summaries (1-hour TTL)
- Cache usage trends (1-day TTL)

### Scalability

- Horizontal scaling support for high-volume monitoring
- Batch processing for metric updates
- Asynchronous billing generation

## Monitoring & Alerting

### Health Checks

- Database connectivity monitoring
- API endpoint availability
- Metric collection accuracy

### Alerts

- High resource usage alerts (>80% CPU/memory)
- Cost threshold alerts (>$100/day)
- Failed metric collection alerts
- Database performance alerts

## Future Enhancements

### Planned Features

- **Real-time Notifications**: WebSocket-based live updates
- **Advanced Analytics**: Machine learning for usage prediction
- **Cost Optimization**: Automated resource scaling recommendations
- **Multi-currency Support**: Support for different currencies
- **Custom Metrics**: User-defined monitoring metrics
- **Integration APIs**: Third-party monitoring tool integration

### Performance Improvements

- **Time-series Database**: Migration to InfluxDB or TimescaleDB
- **Streaming Analytics**: Real-time data processing
- **Distributed Monitoring**: Multi-region monitoring support
- **Advanced Caching**: Redis-based caching layer

## Troubleshooting

### Common Issues

#### Missing Monitoring Data

1. Check if the monitoring service is running
2. Verify database connectivity
3. Check resource tracking integration
4. Review authentication and permissions

#### Incorrect Billing Calculations

1. Verify pricing configuration
2. Check resource usage duration calculations
3. Review billing period generation
4. Validate cost breakdown logic

#### Dashboard Not Loading

1. Check API endpoint availability
2. Verify authentication tokens
3. Review browser console for errors
4. Check network connectivity

### Debug Commands

```bash
# Check monitoring service status
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/v1/monitoring/usage

# Generate test billing
curl -X POST "http://localhost:8080/api/v1/monitoring/billing/generate?year=2024&month=1" \
  -H "Authorization: Bearer <token>"

# Get active resources
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/v1/monitoring/resources/active
```

## Support

For issues and questions related to the monitoring service:

1. **Documentation**: Check this documentation first
2. **API Reference**: Review the Swagger documentation
3. **GitHub Issues**: Report bugs and feature requests
4. **Community**: Join the Unicorn community discussions

---

_This monitoring service provides comprehensive resource tracking, billing, and analytics capabilities for the Unicorn API platform._
