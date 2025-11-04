# BrickBang Server

BrickBang is a modern, modular backend server built with **Go 1.25**, using **Fiber v2** as the web framework and **PostgreSQL (pgx/sqlc)** for database operations.  
It’s designed for clarity, observability, and maintainability — providing consistent API responses, embedded build metadata, and modular architecture.

---

## 🚀 Features

- **Fiber v2** — high-performance HTTP framework.
- **Unified JSON response middleware** for consistent API output (success & error).
- **Metadata embedding** (version, githash, staging, compile time, etc.) injected at build time.
- **Health, Version, and Uptime endpoints** via `/api/v1/aux/*` routes.
- **Service-oriented structure** separating controllers, services, and middleware.
- **Graceful logging and error handling**.
- **PostgreSQL-ready structure** (via pgx + sqlc).

---

## 🧩 Project Structure

```
brickbang/
│
├── cmd/
│   └── bootstrap.go           # Application entry point (main)
│
├── internal/
│   ├── controller/
│   │   └── aux_controller.go  # Health/version/uptime/metadata endpoints
│   │
│   ├── middleware/
│   │   └── unified_response.go # Unified API response middleware
│   │
│   ├── service/
│   │   └── aux_service.go     # Business logic for health/version/uptime
│   │
│   └── meta/
│       └── metadata.go        # Build metadata (version, githash, compile time...)
│
└── go.mod / go.sum
```

---

## 🧱 Middleware: `UnifiedResponse`

All responses (both success and error) are wrapped in a unified JSON structure with useful metadata.

### Example Output

**Success Response:**
```json
{
  "content": {
    "compile": "2025-11-04T09:38:46.586363000Z",
    "githash": "a94e2451",
    "gobuild": "go1.25.3-darwin/arm64",
    "staging": "dev",
    "version": "0.1.0"
  },
  "metadata": {
    "origin": "/api/v1/aux/metadata",
    "request": "77a38713-09f9-43c9-b224-a5884eae55d3",
    "status": 200,
    "timeout": 0,
    "timestamp": "2025-11-04T07:39:16Z"
  }
}
```

---

## ⚙️ Build Metadata

Metadata is compiled into the binary at build time and can be accessed both in runtime and via `/api/v1/aux/metadata`.

Defined in `internal/meta/metadata.go`:

```go
var (
    Version = "none"
    Staging = "none"
    Githash = "none"
    Gobuild = "none"
    Compile = "none"
)

func Metadata() map[string]string {
    return map[string]string{
        "version": Version,
        "staging": Staging,
        "githash": Githash,
        "gobuild": Gobuild,
        "compile": Compile,
    }
}
```

Injected via `go build`:

```bash
go build -ldflags "-X brickbang/internal/meta.Version=0.1.0                    -X brickbang/internal/meta.Staging=dev                    -X brickbang/internal/meta.Githash=$(git rev-parse --short HEAD)                    -X brickbang/internal/meta.Gobuild=$(go version | awk '{print $3"-"$4}')                    -X brickbang/internal/meta.Compile=$(date -u +'%Y-%m-%dT%H:%M:%S.%NZ')"     -o ./bin/brickbang ./cmd/bootstrap.go
```

---

## 🧠 Controller: `AuxController`

Responsible for system info routes:

| Endpoint | Description | Example |
|-----------|--------------|----------|
| `/api/v1/aux/health` | Returns app health info | ✅ |
| `/api/v1/aux/version` | Returns current build version | ✅ |
| `/api/v1/aux/uptime` | Returns uptime since start | ✅ |
| `/api/v1/aux/metadata` | Returns full build metadata | ✅ |

---

## 🛠️ Example Startup

```bash
go run ./cmd/bootstrap.go
```

Example log:
```
Server started on port :8081
Environment: dev
Git Commit: a94e2451
Version: 0.1.0
```

Then open:
```
GET http://localhost:8081/api/v1/aux/metadata
```

---

## 🧩 Example UnifiedResponse Middleware (summary)

- Measures request duration
- Generates request UUID
- Wraps all responses into:
  - `content`: actual response (or error)
  - `metadata`: technical details about the request

---

## 🧪 Example Error Response

```json
{
  "content": {
    "error": {
      "message": "resource not found",
      "code": 404
    }
  },
  "metadata": {
    "timestamp": "2025-11-04T07:39:16Z",
    "request_id": "82e2a1c7-9a10-4e88-95c3-c0e9b66d9aa8",
    "path": "/api/v1/unknown",
    "status": 404,
    "processing_time_ms": 12
  }
}
```

---

## 🧭 Future Enhancements

- Integrate PostgreSQL with sqlc-based repository layer
- Add structured logging (zerolog or slog)
- Add configuration loader (YAML-based)
- Implement graceful shutdown
- Add Prometheus metrics for health/uptime

---

## 📜 License

MIT © 2025 BrickBang Authors
