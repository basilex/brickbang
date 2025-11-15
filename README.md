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
- **RBAC and JWT-based auth** (register, login, logout, refresh, me, block).
- **Redis cache & blacklist service integration**.

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
│   ├── repository/            # SQLC generated queries & database access
│   │
│   └── meta/
│       └── metadata.go        # Build metadata (version, githash, compile time...)
│
└── go.mod / go.sum
```

---

## ⚙️ Build Metadata

Metadata is compiled into the binary at build time and can be accessed both in runtime and via `/api/v1/aux/metadata`.

```go
var (
    Version = "none"
    Staging = "none"
    Githash = "none"
    Gobuild = "none"
    Compile = "none"
)
```

Injected via `go build`:

```bash
go build -ldflags "-X brickbang/internal/meta.Version=0.1.0     -X brickbang/internal/meta.Staging=dev     -X brickbang/internal/meta.Githash=$(git rev-parse --short HEAD)     -X brickbang/internal/meta.Gobuild=$(go version | awk '{print $3"-"$4}')     -X brickbang/internal/meta.Compile=$(date -u +'%Y-%m-%dT%H:%M:%S.%NZ')"     -o ./dist/brb ./cmd/bootstrap.go
```

---

## 🧠 Controller: `AuxController`

Responsible for system info routes:

| Endpoint | Description |
|-----------|--------------|
| `/api/v1/aux/health` | Returns app health info |
| `/api/v1/aux/version` | Returns current build version |
| `/api/v1/aux/uptime` | Returns uptime since start |
| `/api/v1/aux/metadata` | Returns full build metadata |

---

## 🧪 UnifiedResponse Middleware

- Measures request duration
- Generates request UUID
- Wraps all responses into:
  - `content`: actual response (or error)
  - `metadata`: technical details about the request

**Example success response:**

```json
{
  "content": { "version": "0.1.0", "staging": "dev", "githash": "a94e2451" },
  "metadata": { "request": "uuid", "status": 200, "timestamp": "2025-11-04T07:39:16Z" }
}
```

**Example error response:**

```json
{
  "content": { "error": { "message": "resource not found", "code": 404 } },
  "metadata": { "request": "uuid", "status": 404, "path": "/api/v1/unknown" }
}
```

---

## 🧭 Development Flow

BrickBang uses **Makefile** as a central entry point for dev tasks, app builds, DB migrations, Docker, and certificates.

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `BRB_SVC` | `brb` | Service name |
| `BRB_ENV` | `dev` | Environment (dev, staging, prod) |
| `BRB_VER` | `0.1.0` | Version |
| `BRB_DIST` | `dist` | Output folder |
| `BRB_CERT` | `./resource/cert` | TLS certificates folder |

### Makefile Targets

#### Database

```bash
make dbs-gen      # Generate SQLC db layer
make dbs-up       # Install db schema + default data
make dbs-up1      # Migrate up one level
make dbs-down     # Uninstall db schema
make dbs-down1    # Migrate down one level
make dbs-drop     # Drop all schema + data
make dbs-version  # Show current migration version
```

#### App

```bash
make app-tidy     # Go mod tidy
make app-build    # Build binary with embedded metadata
```

#### Docker / Compose

```bash
make app-up       # Run service via docker compose
make app-down     # Stop service
make app-clean    # Remove exited containers & images
make app-prune    # Docker system prune
make app-cert     # Generate self-signed TLS cert
```

#### Swagger

```bash
make docs         # Generate Swagger docs
```

### Example Local Development

```bash
# Build and run locally
make app-build
SERVER_CHILD_PROCESSES=2 BRB_ENV=dev ./dist/brb

# Docker
make app-up
make app-down
```

---

## 🧩 Call Flow Diagram

```
        +------------+
        |  HTTP Req  |
        +-----+------+
              |
              v
       +------+------
       |  Fiber App  |
       +------+------
              |
       +------+------
       | Middleware  |  --> UnifiedResponse, RBAC, Auth
       +------+------
              |
       +------+------
       | Controllers | --> AuxController, AuthController, RoleController
       +------+------
              |
       +------+------
       |  Services   | --> Business logic, Redis, DB operations
       +------+------
              |
       +------+------
       | Repositories| --> SQLC/Postgres or Redis clients
       +-------------+
```

---

## 📜 License

MIT © 2025 BrickBang Authors
