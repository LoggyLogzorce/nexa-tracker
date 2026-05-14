# Nexa Task Tracker

Task tracking API built with Go, Gin, PostgreSQL, and GORM.

## Features

### Implemented
- **JWT Authentication** — registration, login, logout, token refresh with rotation and reuse attack detection
- **User Management** — profile retrieval, update, account deletion with data anonymization
- **Project Management** — full CRUD, scoped listing for owners/participants
- **Task Management** — full CRUD with validation (assignee/status/priority scoped to project), archive support, update history tracking via JSONB diffs
- **Custom Statuses & Priorities** — per-project, with drag-reorder support (order_index), default values (To Do / In Progress / Done; Low / Medium / High)
- **Project Participants** — role-based access (owner / member / read_only)
- **Task Comments** — CRUD with ownership verification, user enrichment
- **Event Bus** — synchronous pub/sub for cross-module communication (e.g. project creation triggers default statuses/priorities, user deletion cascades to notifications)
- **Rate Limiting** — IP-based, configurable (auth endpoints: 5 req/min)
- **Standardized JSON Responses** — consistent `{"success": bool, "data": ..., "error": ...}` envelope

### In Progress / Stubs
- **2FA (TOTP)** — endpoints wired, handlers return placeholder responses
- **File Attachments** — model + repository scaffolded, handlers return placeholder responses
- **Notifications Module** — scaffolded, `Init()` commented out in `main.go`

### Planned
- **Chat Module** — placeholder file only

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.26.1 |
| HTTP Framework | Gin v1.10.0 |
| ORM | GORM v1.30.0 |
| Database | PostgreSQL 16 |
| JWT | golang-jwt v5.2.1 |
| Password Hashing | bcrypt (x/crypto v0.23.0) |
| UUID | google/uuid v1.6.0 |
| Rate Limiting | x/time v0.5.0 |
| Env Loading | godotenv v1.5.1 |
| Validation | go-playground/validator v10.20.0 |
| Containerization | Docker + Docker Compose |

## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌────────────────┐
│   Handler   │────▶│   Service    │────▶│  Repository    │
│ (HTTP only) │     │ (business    │     │ (GORM queries) │
│             │     │  logic)      │     │                │
└─────────────┘     └──────┬───────┘     └────────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │  Event Bus   │
                    │ (pub/sub)    │
                    └──────┬───────┘
                           │
              ┌────────────┼────────────┐
              ▼            ▼            ▼
         status svc   priority svc   notify svc
```

- **Handler layer** — binds requests, extracts context, validates input, returns responses
- **Service layer** — business logic, timeout management (5-10s context deadlines)
- **Repository layer** — GORM queries, transactional writes
- **Event Bus** — synchronous publish/subscribe for decoupled module communication
- **Middleware stack** — rate limiter → JWT auth → RBAC (project-level)

## Project Structure

```
├── cmd/app/main.go              — Application entry point
├── docker/
│   ├── Dockerfile               — Multi-stage build
│   └── docker-compose.yml       — App + PostgreSQL 16
├── internal/
│   ├── api/router.go            — Route definitions and wiring
│   ├── config/config.go         — Environment-based configuration
│   ├── ctxkeys/ctxkeys.go       — Context key constants
│   ├── db/
│   │   ├── db.go                — GORM connection + auto-migration
│   │   └── schema.sql           — Reference SQL schema
│   ├── middleware/
│   │   ├── auth.go              — JWT authentication
│   │   ├── rbac.go              — Project-level RBAC (owner/member/read_only)
│   │   └── ratelimit.go         — IP-based rate limiter
│   ├── core/
│   │   ├── auth/                — Auth (register, login, refresh, logout, 2FA stub)
│   │   ├── user/                — User management
│   │   ├── project/             — Project CRUD
│   │   ├── task/                — Task CRUD + update history
│   │   ├── status/              — Per-project statuses
│   │   ├── priority/            — Per-project priorities
│   │   ├── participant/         — Project participants
│   │   ├── comment/             — Task comments
│   │   └── attachment/          — Attachments (stub)
│   ├── modules/
│   │   ├── notify/              — Notifications (stub, commented out)
│   │   └── chat/chat.go         — Chat (placeholder)
│   └── pkg/
│       ├── events/              — Event bus + event structs
│       ├── hash/                — bcrypt + SHA-256 token hashing
│       ├── jwt/                 — Access/refresh JWT tokens
│       ├── nullable/            — Nullable types for PATCH updates
│       ├── response/            — Standardized JSON responses
│       └── validation/          — Request validation helpers
└── .env.example                 — Environment template
```

## Getting Started

### Prerequisites

- Go 1.26+ (or Docker)
- PostgreSQL 16

### Local Development

```bash
# Clone and navigate
git clone <repo-url> && cd nexa-task-tracker

# Configure environment
cp .env.example .env
# Edit .env with your database credentials and JWT secret

# Run
go mod tidy
go run ./cmd/app
```

### With Docker

```bash
docker compose -f docker/docker-compose.yml up --build
```

Server starts at `http://localhost:8080`.

## Configuration

All configuration is via environment variables (see `.env.example`):

| Variable | Default | Description |
|---|---|---|
| `SERVER_HOST` | `0.0.0.0` | Bind address |
| `SERVER_PORT` | `8080` | HTTP port |
| `DB_HOST` | `127.0.0.1` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | — | Database password |
| `DB_NAME` | `nexa_tracker` | Database name |
| `DB_SSLMODE` | `disable` | PostgreSQL SSL mode |
| `JWT_SECRET` | — | HMAC signing key |
| `JWT_ACCESS_EXPIRY` | `15m` | Access token TTL |
| `JWT_REFRESH_EXPIRY` | `168h` (7d) | Refresh token TTL |
| `COOKIE_DOMAIN` | — | Cookie domain |
| `COOKIE_SECURE` | `false` | Cookie Secure flag |
| `COOKIE_SAMESITE` | `1` | SameSite (1=Default, 2=Lax, 3=Strict, 4=None) |
| `NOTIFY_MODULE` | `false` | Enable notifications module |

## API Endpoints

### Health

| Method | Path | Auth |
|---|---|---|
| `GET` | `/health` | — |

### Auth (`/api/v1/auth`)

Rate-limited (5 req/min). Cookies used for token storage.

| Method | Path | Status |
|---|---|---|
| `POST` | `/register` | ✅ |
| `POST` | `/login` | ✅ |
| `POST` | `/refresh` | ✅ |
| `POST` | `/logout` | ✅ |
| `POST` | `/2fa/setup` | 🚧 stub |
| `POST` | `/2fa/verify` | 🚧 stub |
| `POST` | `/2fa/enable` | 🚧 stub |
| `POST` | `/2fa/disable` | 🚧 stub |

### Users (`/api/v1/users`)

Requires JWT auth.

| Method | Path | Status |
|---|---|---|
| `GET` | `/me` | ✅ |
| `PUT` | `/me` | ✅ |
| `DELETE` | `/me` | ✅ |

### Projects (`/api/v1/projects`)

Requires JWT auth. Participants roles: `owner` (read+write+delete), `member` (read+write), `read_only` (read).

| Method | Path | Access | Status |
|---|---|---|---|
| `GET` | `/` | authenticated | ✅ |
| `POST` | `/` | authenticated | ✅ |
| `GET` | `/:id` | read_only+ | ✅ |
| `PUT` | `/:id` | owner | ✅ |
| `DELETE` | `/:id` | owner | ✅ |
| `GET` | `/:id/participants` | read_only+ | ✅ |
| `POST` | `/:id/participants` | owner | ✅ |
| `PUT` | `/:id/participants/:user_id` | owner | ✅ |
| `DELETE` | `/:id/participants/:user_id` | owner | ✅ |
| `GET` | `/:id/statuses` | read_only+ | ✅ |
| `POST` | `/:id/statuses` | member+ | ✅ |
| `PUT` | `/:id/statuses/:status_id` | member+ | ✅ |
| `DELETE` | `/:id/statuses/:status_id` | owner | ✅ |
| `GET` | `/:id/priorities` | read_only+ | ✅ |
| `POST` | `/:id/priorities` | member+ | ✅ |
| `PUT` | `/:id/priorities/:priority_id` | member+ | ✅ |
| `DELETE` | `/:id/priorities/:priority_id` | owner | ✅ |

### Tasks (`/api/v1/projects/:id/tasks`)

Requires JWT auth + project access. History tracked as JSONB diffs in `update_history` table.

| Method | Path | Access | Status |
|---|---|---|---|
| `GET` | `/tasks` | read_only+ | ✅ |
| `POST` | `/` | member+ | ✅ |
| `GET` | `/:task_id` | read_only+ | ✅ |
| `PUT` | `/:task_id` | member+ | ✅ |
| `DELETE` | `/:task_id` | owner | ✅ |
| `GET` | `/:task_id/history` | read_only+ | ✅ |
| `GET` | `/:task_id/comments` | member+ | ✅ |
| `POST` | `/:task_id/comments` | member+ | ✅ |
| `PUT` | `/:task_id/comments/:comment_id` | member+ | ✅ |
| `DELETE` | `/:task_id/comments/:comment_id` | member+ | ✅ |

### Attachments (`/api/v1/projects/:id/tasks/:task_id/attachments`)

| Method | Path | Status |
|---|---|---|
| `GET` | `/` | 🚧 stub |
| `POST` | `/` | 🚧 stub |
| `GET` | `/:attachment_id` | 🚧 stub |
| `DELETE` | `/:attachment_id` | 🚧 stub |

### Notifications (`/api/v1/notifications`)

Commented out in `main.go`. Not available by default.

## Database Schema

### Key Tables

| Table | Description |
|---|---|
| `users` | Core user accounts (uuid PK, email, password_hash, name, role, 2fa secret) |
| `refresh_tokens` | JWT refresh token storage with revocation tracking and reuse detection |
| `projects` | Project entities owned by a user |
| `project_participants` | Many-to-many with roles (owner/member/read_only) |
| `statuses` | Per-project task statuses with order_index for drag-reorder |
| `priorities` | Per-project task priorities |
| `tasks` | Tasks with references to project, status, priority, assignee, reporter |
| `update_history` | JSONB-based field-level change tracking |
| `comments` | Task-scoped comments |
| `attachments` | File metadata (stub implementation) |

See `internal/db/schema.sql` for the complete schema.

## Development

### Adding a New Core Module

1. Create `internal/core/<module>/` with `model.go`, `handler.go`, `service.go`, `repository.go`, `errors.go`
2. Implement the `Repository` interface with GORM
3. Implement business logic in the `Service`
4. Wire HTTP handlers in the `Handler`
5. Register routes in `internal/api/router.go`
6. Initialize in `cmd/app/main.go`

### Code Style

- Three-layer architecture: Handler → Service → Repository
- Services use `context.WithTimeout` for all database operations
- Use `pkg/nullable` types for PATCH endpoints to distinguish "not provided" from "null"
- Use `pkg/response` helpers for consistent JSON formatting
- Events for cross-module communication go through the synchronous EventBus

### Known Issues

- Dockerfile uses Go 1.22; `go.mod` requires Go 1.26.1 — update the Dockerfile base image if building with Docker
- `.env` DB name typo: `DBHOST` (missing underscore) should be `DB_HOST`

## License

MIT
