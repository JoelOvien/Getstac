# Go XLSX Upload API

A high-performance HTTP API built in Go 1.22 that allows users to upload `.xlsx` files, parse their rows into structured records, and query them with pagination support.

## Features

- **XLSX File Upload**: Parse and validate Excel files with automatic header detection
- **Concurrent Processing**: Handle ~100 concurrent users with worker pools and bounded concurrency
- **Rate Limiting**: Per-IP rate limiting to prevent abuse
- **API Key Authentication**: Optional API key authentication
- **Structured Logging**: Request logging with zerolog
- **Graceful Shutdown**: Proper cleanup on SIGTERM/SIGINT
- **Health Check**: Built-in health endpoint for monitoring
- **Pagination**: Efficient record listing with offset/limit support
- **Docker Support**: Full containerization with Docker and docker-compose

## Architecture

The project follows clean architecture principles with clear separation of concerns:

```
├── cmd/server/          # Application entry point
├── internal/
│   ├── api/             # HTTP routing and handlers
│   │   ├── handlers/    # Request handlers
│   │   ├── middleware/  # Custom middleware
│   │   └── router.go    # Route configuration
│   ├── config/          # Configuration management
│   ├── models/          # Data models
│   ├── storage/         # In-memory storage implementation
│   └── xlsx/            # XLSX parsing logic
└── tests/               # Unit tests
```

## Requirements

- Go 1.22+
- Docker (optional, for containerized deployment)

## Quick Start

### Using Docker (Recommended)

1. **Build and run with docker-compose:**
   ```bash
   docker-compose up --build
   ```

2. **The API will be available at:** `http://localhost:8080`

### Local Development

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Configure environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your preferred settings
   ```

3. **Run the application:**
   ```bash
   go run cmd/server/main.go
   ```

4. **Run tests:**
   ```bash
   go test ./... -v
   ```

## API Endpoints

### Health Check
```bash
GET /healthz
```

**Response:**
```json
{
  "status": "ok"
}
```

### Upload XLSX File
```bash
POST /v1/uploads
Content-Type: multipart/form-data
X-API-Key: secret123

Form field: file (must be .xlsx)
```

**Response:**
```json
{
  "uploadId": "550e8400-e29b-41d4-a716-446655440000",
  "rowsAccepted": 150,
  "rowsRejected": 5
}
```

**Example using curl:**
```bash
curl -X POST http://localhost:8080/v1/uploads \
  -H "X-API-Key: secret123" \
  -F "file=@sample.xlsx"
```

### List Records
```bash
GET /v1/records?limit=10&offset=0
X-API-Key: secret123
```

**Query Parameters:**
- `limit` (optional): Number of records to return (default: 10, max: 1000)
- `offset` (optional): Number of records to skip (default: 0)

**Response:**
```json
{
  "records": [
    {
      "id": "uuid-1",
      "uploadId": "upload-uuid",
      "data": {
        "Name": "John Doe",
        "Email": "john@example.com",
        "Age": "30"
      },
      "createdAt": "2025-11-09T10:30:00Z"
    }
  ],
  "total": 150,
  "limit": 10,
  "offset": 0
}
```

**Example using curl:**
```bash
curl "http://localhost:8080/v1/records?limit=20&offset=0" \
  -H "X-API-Key: secret123"
```

## Configuration

Configuration is managed through environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `MAX_UPLOAD_SIZE_MB` | Maximum file upload size in MB | `10` |
| `RATE_LIMIT` | Requests per minute per IP | `100` |
| `API_KEY` | API key for authentication (empty = disabled) | `""` |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |
| `SHUTDOWN_TIMEOUT` | Graceful shutdown timeout | `30s` |
| `REQUEST_TIMEOUT` | Maximum request processing time | `60s` |
| `WORKER_POOL_SIZE` | Number of workers for row processing | `10` |

## Error Handling

The API returns consistent JSON error responses:

```json
{
  "code": "error_code",
  "message": "Human-readable error message"
}
```

### Common Error Codes

- `bad_request`: Invalid request parameters or malformed data
- `invalid_file_type`: Non-.xlsx file uploaded
- `invalid_content_type`: Incorrect content type header
- `invalid_headers`: Missing or invalid XLSX headers
- `parse_error`: Failed to parse XLSX file
- `rate_limit_exceeded`: Too many requests
- `missing_api_key`: API key not provided
- `invalid_api_key`: Incorrect API key
- `internal_error`: Server-side error

## XLSX File Requirements

- File must have `.xlsx` extension
- Must contain at least one sheet
- First row is treated as headers (all headers must be non-empty)
- Minimum of 2 rows (1 header + 1 data row)
- Maximum file size: 10MB (configurable)

### Example XLSX Structure

| Name | Email | Age | City |
|------|-------|-----|------|
| John Doe | john@example.com | 30 | New York |
| Jane Smith | jane@example.com | 25 | San Francisco |

## Concurrency & Performance

### Design Decisions

1. **Worker Pool Pattern**:
   - Configurable number of workers process rows concurrently
   - Prevents memory exhaustion with bounded concurrency
   - Default pool size: 10 workers

2. **Context-Based Cancellation**:
   - All requests have timeouts (default: 60s)
   - Graceful cancellation propagates through worker pools
   - Protects against hung requests

3. **Rate Limiting**:
   - Token bucket algorithm per IP address
   - Automatic cleanup of stale buckets
   - Prevents DoS attacks

4. **Memory Safety**:
   - Thread-safe in-memory storage with RWMutex
   - Pagination prevents loading entire dataset
   - File size limits prevent memory exhaustion

### Performance Characteristics

- **Concurrent Users**: Designed for ~100 concurrent users
- **Upload Processing**: O(n) where n = number of rows
- **Record Listing**: O(1) for pagination (in-memory slice access)
- **Memory Usage**: ~1KB per record (approximate)

## Horizontal Scaling

For production deployments requiring higher throughput:

### 1. Load Balancer
Deploy multiple API instances behind a load balancer (nginx, HAProxy, AWS ALB):

```
[Load Balancer]
    ├── API Instance 1
    ├── API Instance 2
    └── API Instance 3
```

### 2. Shared Storage Layer
Replace in-memory storage with a shared database:

**Option A: PostgreSQL**
- Persistent storage
- ACID compliance
- Full-text search capabilities

**Option B: Redis**
- High-performance caching
- Built-in data structures
- Pub/sub for real-time updates

**Option C: MongoDB**
- Document-based storage
- Flexible schema
- Horizontal sharding

### 3. Rate Limiting
Use Redis-backed rate limiting for distributed rate limiting:
- Shared rate limit counters across instances
- More accurate limit enforcement
- Examples: `redis-rate-limiter`, `go-redis/rate`

### 4. File Processing Queue
For heavy workloads, use asynchronous processing:
- Upload endpoint returns immediately with job ID
- Background workers process files from queue (RabbitMQ, AWS SQS)
- Status endpoint to check processing progress

## Testing

### Run All Tests
```bash
go test ./... -v
```

### Run Specific Tests
```bash
go test ./tests -v -run TestHealthHandler
```

### Test Coverage
```bash
go test ./... -cover
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Load Testing
Use tools like `hey`, `wrk`, or `ab` for load testing:

```bash
# Install hey
go install github.com/rakyll/hey@latest

# Test upload endpoint (requires sample.xlsx)
hey -n 1000 -c 50 -m POST \
  -H "X-API-Key: secret123" \
  -T "multipart/form-data; boundary=----WebKitFormBoundary" \
  http://localhost:8080/v1/uploads

# Test list endpoint
hey -n 10000 -c 100 \
  -H "X-API-Key: secret123" \
  http://localhost:8080/v1/records?limit=10
```

## Security

### Current Implementation

1. **File Validation**: Extension and content type validation
2. **Size Limits**: Enforced maximum upload size
3. **Rate Limiting**: Per-IP request throttling
4. **API Key Authentication**: Optional header-based authentication
5. **Request Timeouts**: Prevents resource exhaustion
6. **Input Validation**: Query parameter validation and sanitization

### Production Recommendations

1. **HTTPS Only**: Use TLS certificates (Let's Encrypt)
2. **API Key Rotation**: Implement key rotation policies
3. **JWT Tokens**: Replace static API keys with JWT
4. **CORS Configuration**: Restrict allowed origins
5. **File Scanning**: Integrate antivirus scanning for uploads
6. **Audit Logging**: Log all authenticated actions
7. **Secret Management**: Use vault solutions (HashiCorp Vault, AWS Secrets Manager)

## Monitoring & Observability

### Current Implementation
- Structured JSON logging with zerolog
- HTTP request logging (method, path, status, duration)
- Upload processing metrics (rows accepted/rejected)

### Production Additions
1. **Metrics Endpoint**: Prometheus-compatible `/metrics` endpoint
2. **Distributed Tracing**: OpenTelemetry integration
3. **APM**: DataDog, New Relic, or similar
4. **Error Tracking**: Sentry integration
5. **Health Checks**: Deep health checks (database connectivity, disk space)

## Troubleshooting

### Common Issues

**Issue: "File size exceeds maximum allowed size"**
- Solution: Increase `MAX_UPLOAD_SIZE_MB` environment variable

**Issue: "Rate limit exceeded"**
- Solution: Increase `RATE_LIMIT` or implement API key-based limits

**Issue: "Failed to parse XLSX file"**
- Solution: Ensure file has proper structure (headers + data rows)

**Issue: Container health check failing**
- Solution: Ensure port 8080 is not already in use

## Development

### Adding New Middleware
```go
// internal/api/middleware/custom.go
func CustomMiddleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Pre-processing
            next.ServeHTTP(w, r)
            // Post-processing
        })
    }
}
```

### Adding New Handlers
```go
// internal/api/handlers/custom.go
type CustomHandler struct {
    storage *storage.MemoryStorage
    logger  *zerolog.Logger
}

func (h *CustomHandler) Handle(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

## License

MIT License - feel free to use this project for any purpose.

## Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## Assumptions & Limitations

### Assumptions
- XLSX files follow standard structure (header row + data rows)
- All required data fits in memory (for in-memory storage)
- Single server deployment (no distributed coordination)
- API keys are pre-shared (no dynamic key generation)

### Current Limitations
1. **Storage**: In-memory only (data lost on restart)
2. **No Persistence**: Records are not saved to disk
3. **No Authentication Management**: Static API key only
4. **Single Sheet**: Only processes first sheet in workbook
5. **File Size**: Limited to configured max upload size
6. **No Filtering**: List endpoint doesn't support filtering by upload ID or fields

### Future Enhancements
- Database persistence (PostgreSQL, MongoDB)
- Multi-sheet processing
- Advanced filtering and search
- Export functionality (CSV, JSON)
- Webhook notifications on upload completion
- Real-time upload progress tracking
- Batch delete operations
- Field-level validation rules

## Support

For issues, questions, or contributions, please open an issue on GitHub.
