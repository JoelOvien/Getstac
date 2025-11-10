# Project Structure

This document outlines the complete structure of the project.

```
go-xlsx-api/
│
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
│
├── internal/                       # Private application code
│   ├── api/                        # HTTP layer
│   │   ├── handlers/               # Request handlers
│   │   │   ├── health.go           # Health check handler
│   │   │   ├── list.go             # List records handler
│   │   │   └── upload.go           # Upload XLSX handler
│   │   ├── middleware/             # HTTP middleware
│   │   │   ├── auth.go             # API key authentication
│   │   │   ├── logger.go           # Request logging
│   │   │   ├── ratelimit.go        # Rate limiting
│   │   │   └── timeout.go          # Request timeout
│   │   └── router.go               # Route configuration
│   │
│   ├── config/
│   │   └── config.go               # Configuration management
│   │
│   ├── models/
│   │   └── models.go               # Data structures
│   │
│   ├── storage/
│   │   └── memory.go               # In-memory storage implementation
│   │
│   └── xlsx/
│       └── parser.go               # XLSX parsing logic
│
├── pkg/                            # Public libraries (empty for now)
│   └── utils/
│
├── tests/                          # Unit tests
│   ├── handlers_test.go            # Handler tests
│   ├── middleware_test.go          # Middleware tests
│   ├── parser_test.go              # Parser tests
│   └── storage_test.go             # Storage tests
│
├── .env.example                    # Environment variable template
├── .gitignore                      # Git ignore rules
├── docker-compose.yml              # Docker Compose configuration
├── Dockerfile                      # Docker image definition
├── go.mod                          # Go module definition
├── go.sum                          # Go module checksums
├── Makefile                        # Build automation
├── README.md                       # Project documentation
├── STRUCTURE.md                    # This file
└── test_api.sh                     # API testing script
```

## Component Descriptions

### cmd/server/main.go
The application entry point that:
- Loads configuration from environment variables
- Sets up structured logging
- Creates and configures the HTTP router
- Starts the HTTP server
- Handles graceful shutdown on SIGTERM/SIGINT

### internal/api/
HTTP layer components:

**handlers/**
- `health.go`: Returns service health status
- `list.go`: Lists records with pagination
- `upload.go`: Processes XLSX file uploads

**middleware/**
- `auth.go`: Validates API keys
- `logger.go`: Logs HTTP requests
- `ratelimit.go`: Implements per-IP rate limiting
- `timeout.go`: Adds request timeouts

**router.go**
- Configures routes and middleware chain
- Creates handler instances
- Sets up global middleware

### internal/config/
Configuration management with environment variable support:
- PORT
- MAX_UPLOAD_SIZE_MB
- RATE_LIMIT
- API_KEY
- LOG_LEVEL
- Timeout settings
- Worker pool size

### internal/models/
Data structures:
- `Record`: Parsed XLSX row
- `UploadResponse`: Upload result
- `ListRecordsResponse`: Paginated list response
- `HealthResponse`: Health check response
- `ErrorResponse`: Standardized error format

### internal/storage/
In-memory storage with thread-safe operations:
- Store records
- List records with pagination
- Get records by upload ID
- Thread-safe with RWMutex

### internal/xlsx/
XLSX parsing:
- Stream processing with worker pools
- Header validation
- Row-by-row parsing
- Context-aware cancellation
- Bounded concurrency

### tests/
Comprehensive unit tests covering:
- HTTP handlers
- Middleware (auth, rate limiting, timeout)
- Storage operations
- XLSX parsing
- Error handling

## Data Flow

### Upload Flow
```
Client → Router → Middleware Chain → Upload Handler
                                           ↓
                                     XLSX Parser
                                           ↓
                                    Worker Pool
                                           ↓
                                  Memory Storage
                                           ↓
                                    Response ← Client
```

### List Flow
```
Client → Router → Middleware Chain → List Handler
                                           ↓
                                  Memory Storage
                                           ↓
                                  Pagination Logic
                                           ↓
                                    Response ← Client
```

## Middleware Chain
```
Request
  ↓
Recoverer (panic recovery)
  ↓
RequestID (generate request ID)
  ↓
RealIP (extract real client IP)
  ↓
Logger (log request)
  ↓
Timeout (add context timeout)
  ↓
RateLimiter (check rate limits)
  ↓
APIKeyAuth (validate API key)
  ↓
Handler
  ↓
Response
```

## Key Design Patterns

### Clean Architecture
- Clear separation between layers
- Dependency injection
- Interface-based design

### Concurrency
- Worker pools for bounded parallelism
- Context-based cancellation
- Thread-safe data structures

### Middleware Pattern
- Composable request processing
- Cross-cutting concerns isolation
- Reusable components

### Repository Pattern
- Storage abstraction
- Easy to swap implementations
- Testable without I/O

## Testing Strategy

### Unit Tests
- Table-driven tests
- Mock-free where possible
- Test all edge cases
- Validate error handling

### Integration Points
- HTTP handlers tested with httptest
- Middleware tested with real handlers
- Storage tested with actual operations

## Configuration Management

Environment variables with sensible defaults:
```
PORT=8080
MAX_UPLOAD_SIZE_MB=10
RATE_LIMIT=100
API_KEY=secret123
LOG_LEVEL=info
SHUTDOWN_TIMEOUT=30s
REQUEST_TIMEOUT=60s
WORKER_POOL_SIZE=10
```

## Build & Deployment

### Local Development
```bash
make build    # Build binary
make run      # Build and run
make test     # Run tests
```

### Docker
```bash
make docker-build    # Build image
make docker-run      # Start containers
make docker-stop     # Stop containers
```

### Production Considerations
- Use environment-specific configs
- Enable HTTPS/TLS
- Configure proper logging levels
- Set up monitoring and alerting
- Use secret management for API keys
- Consider horizontal scaling
