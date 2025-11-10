# Quick Start Guide

Get the Go XLSX Upload API up and running in minutes!

## Prerequisites

- Go 1.22+ (for local development)
- Docker & Docker Compose (for containerized deployment)

## Option 1: Docker (Recommended)

### 1. Start the Application

```bash
docker-compose up --build
```

This will:
- Build the Docker image
- Start the API server on port 8080
- Apply the default configuration from docker-compose.yml

### 2. Test the API

```bash
# Health check
curl http://localhost:8080/healthz

# Expected response:
# {"status":"ok"}
```

### 3. Upload an XLSX File

Create a sample Excel file with this structure:

| Name | Email | Age | City |
|------|-------|-----|------|
| John Doe | john@example.com | 30 | New York |
| Jane Smith | jane@example.com | 25 | San Francisco |

Then upload it:

```bash
curl -X POST http://localhost:8080/v1/uploads \
  -H "X-API-Key: secret123" \
  -F "file=@sample.xlsx"
```

### 4. List Records

```bash
curl http://localhost:8080/v1/records?limit=10&offset=0 \
  -H "X-API-Key: secret123"
```

### 5. Stop the Application

```bash
docker-compose down
```

## Option 2: Local Development

### 1. Install Dependencies

```bash
go mod download
```

### 2. Configure Environment (Optional)

```bash
cp .env.example .env
# Edit .env if you want to change defaults
```

### 3. Build the Application

```bash
go build -o server ./cmd/server
```

Or use the Makefile:

```bash
make build
```

### 4. Run the Application

```bash
./server
```

Or use the Makefile:

```bash
make run
```

### 5. Test the API

Same as Docker option - see step 2-4 above.

## Quick Test Script

Run the automated test script:

```bash
./test_api.sh
```

This will test:
- Health check
- File upload (if sample.xlsx exists)
- Record listing
- Pagination
- API key validation

## Common Commands

### Using Makefile

```bash
make help           # Show all available commands
make build          # Build the application
make run            # Build and run
make test           # Run all tests
make test-coverage  # Run tests with coverage report
make docker-build   # Build Docker image
make docker-run     # Run in Docker
make docker-stop    # Stop Docker containers
make clean          # Clean build artifacts
```

### Manual Commands

```bash
# Run tests
go test ./... -v

# Format code
go fmt ./...

# Build for production
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Run with custom configuration
PORT=3000 API_KEY=mykey ./server
```

## Configuration

Configure via environment variables:

```bash
# Server
export PORT=8080

# Upload limits
export MAX_UPLOAD_SIZE_MB=10

# Rate limiting (requests per minute)
export RATE_LIMIT=100

# Security (leave empty to disable)
export API_KEY=secret123

# Logging
export LOG_LEVEL=info

# Timeouts
export SHUTDOWN_TIMEOUT=30s
export REQUEST_TIMEOUT=60s

# Performance
export WORKER_POOL_SIZE=10
```

## Sample cURL Commands

### Upload with Different Scenarios

```bash
# Successful upload
curl -X POST http://localhost:8080/v1/uploads \
  -H "X-API-Key: secret123" \
  -F "file=@data.xlsx"

# Missing API key (will fail)
curl -X POST http://localhost:8080/v1/uploads \
  -F "file=@data.xlsx"

# Wrong file type (will fail)
curl -X POST http://localhost:8080/v1/uploads \
  -H "X-API-Key: secret123" \
  -F "file=@document.pdf"
```

### Listing with Pagination

```bash
# First page (10 records)
curl "http://localhost:8080/v1/records?limit=10&offset=0" \
  -H "X-API-Key: secret123"

# Second page (10 records)
curl "http://localhost:8080/v1/records?limit=10&offset=10" \
  -H "X-API-Key: secret123"

# Get 50 records starting from record 100
curl "http://localhost:8080/v1/records?limit=50&offset=100" \
  -H "X-API-Key: secret123"
```

## Troubleshooting

### Port Already in Use

```bash
# Check what's using port 8080
lsof -i :8080

# Change the port
export PORT=3000
./server
```

### Docker Container Won't Start

```bash
# Check logs
docker-compose logs

# Rebuild from scratch
docker-compose down
docker-compose up --build --force-recreate
```

### Tests Failing

```bash
# Clean and rebuild
make clean
go mod tidy
go test ./... -v
```

### File Upload Fails

Check:
1. File is .xlsx format
2. File has header row
3. File is under 10MB (or configured MAX_UPLOAD_SIZE_MB)
4. API key is correct

## Next Steps

1. Read the full [README.md](README.md) for detailed documentation
2. Check [STRUCTURE.md](STRUCTURE.md) to understand the architecture
3. Review test files in `tests/` for usage examples
4. Customize configuration for your use case
5. Consider implementing database persistence for production

## Production Deployment

For production use:

1. **Use HTTPS**: Add TLS/SSL certificates
2. **Use a real database**: Replace in-memory storage with PostgreSQL/MongoDB
3. **Enable monitoring**: Add Prometheus metrics
4. **Set up logging**: Use centralized logging (ELK, Datadog, etc.)
5. **Configure secrets**: Use vault solutions (not environment variables)
6. **Scale horizontally**: Deploy multiple instances behind a load balancer
7. **Add health checks**: Configure deep health checks for orchestration

## Support

For issues or questions:
- Check the README.md
- Review the code in `internal/`
- Run the test suite for examples
- Check error responses for debugging info

Happy coding! 🚀
