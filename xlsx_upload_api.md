# Go XLSX Upload API — Full Project Specification

## Goal
Build a **high-performance HTTP API in Go (1.21+)** that allows users to upload `.xlsx` files, parse their rows into structured records, and query them.  
The system must handle **~100 concurrent users** gracefully while maintaining memory safety, clear error handling, and clean architecture.

---

## ✅ Core Requirements

### 1. Language & Framework
- Go **1.21+ (preferably 1.22)**
- Use **any lightweight router** (`chi`, `gin`, or `fiber`) — or `net/http` directly.

---

### 2. Upload Endpoint

**`POST /v1/uploads`**

- Content type: `multipart/form-data` (field name = `file`)
- Accept **only `.xlsx`** files; reject other formats with descriptive JSON errors.
- Max upload size: **10MB** (configurable via environment variable).
- Use `excelize` to parse the XLSX file.
- Validate **required headers** (assume headers exist).  
  - If missing → return HTTP 400 with a descriptive message.
- Parse rows into structured records.
- Return JSON response:
  ```json
  {
    "uploadId": "<uuid>",
    "rowsAccepted": <n>,
    "rowsRejected": <m>
  }
  ```
- Can be synchronous or asynchronous (if async, return a status reference).

---

### 3. Listing Endpoint

**`GET /v1/records?limit=&offset=`**

- Return paginated list of parsed records in JSON.
- Include total record count.

---

### 4. Health Endpoint

**`GET /healthz`**

- Returns HTTP 200 with payload:
  ```json
  { "status": "ok" }
  ```

---

### 5. Error Handling

- Consistent JSON error structure, e.g.:
  ```json
  { "code": "bad_request", "message": "Missing required headers: Name, Email" }
  ```
- Handle common cases:
  - Invalid file type
  - Missing/invalid headers
  - Oversized files
  - Timeout or cancellation
  - Internal errors

---

## ⚙️ Concurrency & Performance

- Design for **~100 concurrent uploads/reads**.
- Use **context** for request timeouts and cancellation.
- Avoid loading entire files into memory; **stream parse rows** if possible.
- Use **worker pools** or bounded concurrency for row processing.
- Implement a **basic rate limiter** (per IP or API key).
- Write notes in the README on how to scale horizontally (e.g. using load balancers, shared DB, or Redis).

---

## 💾 Storage

- Either:
  - In-memory storage (map + mutex + pagination support), or
  - Lightweight DB (e.g. SQLite or PostgreSQL).
- If DB used:
  - Provide simple schema and migration SQL file.
  - Include in Docker setup.

---

## 🧰 Quality & Operational Requirements

- **Graceful shutdown** on SIGTERM (time-bounded server stop).
- **Structured logging** (e.g. `zerolog`, `logrus`, or `zap`).
- Minimal **metrics** (e.g. `/metrics` endpoint or internal counters).
- **Configuration** via environment variables (`.env` or flags):
  - `PORT`
  - `MAX_UPLOAD_SIZE_MB`
  - `RATE_LIMIT`
  - `API_KEY` (optional)
- Include **unit tests** for:
  - XLSX parsing logic
  - HTTP handlers
  - Validation and error responses
  - (Use **table-driven tests**)

---

## 🔒 Security

- Validate `Content-Type` and file extension.
- Enforce upload size limits and request timeouts.
- Reject files with unexpected sheet names or no data.
- Optional: Require API key in header:  
  `X-API-Key: <key>`  
  - Compare with value from env var.

---

## 🐳 Dockerization

- Provide a **Dockerfile** and (if DB used) a **docker-compose.yml**.
- Build commands:
  ```bash
  docker build -t go-xlsx-api .
  docker run -p 8080:8080 go-xlsx-api
  ```
- Expose port 8080 by default.

---

## 🧾 Deliverables

- Full Go project (main.go, routes, handlers, storage, models, config, middleware, tests).
- **README.md** containing:
  - Setup & run instructions (with Docker)
  - API endpoints & sample `curl` or Postman examples
  - Assumptions & limitations
  - Notes on concurrency design & memory usage
- **Tests** runnable via:
  ```bash
  go test ./... -v
  ```
- Must pass without errors.

---

## 🧪 Example API Flow

### Upload a file
```bash
curl -X POST http://localhost:8080/v1/uploads   -H "X-API-Key: secret123"   -F "file=@data.xlsx"
```

### List records
```bash
curl http://localhost:8080/v1/records?limit=10&offset=0
```

### Health check
```bash
curl http://localhost:8080/healthz
```

---

## 🧱 Suggested Folder Structure

```
/go-xlsx-api
├── cmd/
│   └── server/main.go
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── router.go
│   ├── config/
│   ├── storage/
│   ├── models/
│   └── xlsx/
├── pkg/
│   └── utils/
├── tests/
│   ├── handler_test.go
│   └── parser_test.go
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

---

## 🧠 Copilot Task

> **Your goal:** Generate the entire project codebase based on the above specification.
> 
> Include:
> - Fully functional Go HTTP API (Go 1.24)
> - XLSX parsing via `excelize`
> - Middleware for logging, rate limiting, API key validation
> - Graceful shutdown logic
> - Config via env vars
> - Unit tests and Docker setup
> - README with setup instructions

Be idiomatic, robust, and well-commented throughout.
