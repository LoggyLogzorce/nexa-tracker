[English](README.md) | [Русский](README.ru.md)

# Nexa Task Tracker

Task tracking API built with Go, Gin, PostgreSQL, and GORM.

## Features

### Implemented
- **JWT Authentication** — registration, login, logout, token refresh with rotation and reuse attack detection
- **OAuth 2.0 Login** — Google and Yandex OAuth integration
- **Session Management** — list active sessions, revoke specific sessions
- **User Management** — profile retrieval/update, avatar upload, password change, account deletion with data anonymization, user search
- **Project Management** — full CRUD, scoped listing (owned/participated), search
- **Task Management** — full CRUD with validation (assignee/status/priority scoped to project), archive/unarchive, global search, per-user listing (assigned/reported), update history tracking via JSONB diffs
- **Custom Statuses & Priorities** — per-project, with drag-reorder support (`order_index`), default values (To Do / In Progress / Done; Low / Medium / High)
- **Project Participants** — role-based access (owner / member / read_only)
- **Task Comments** — CRUD with ownership verification, user enrichment
- **File Attachments** — upload, download, delete with file metadata tracking
- **Event Bus** — synchronous pub/sub for cross-module communication (e.g. project creation triggers default statuses/priorities, user deletion cascades)
- **Rate Limiting** — IP-based, configurable (auth endpoints: 5 req/min)
- **Standardized JSON Responses** — consistent `{"success": bool, "data": ..., "error": ...}` envelope
- **Version & Changelog** — build info and release notes endpoints

### In Progress / Stubs
- **2FA (TOTP)** — endpoints wired, handlers return placeholder responses
- **Notifications Module** — scaffolded, `Init()` commented out in `main.go`

### Planned
- **Chat Module** — placeholder file only

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.26.1 |
| HTTP Framework | Gin v1.12.0 |
| ORM | GORM v1.30.0 (with datatypes for JSONB) |
| Database | PostgreSQL 16 |
| JWT | golang-jwt v5.2.1 |
| OAuth 2.0 | golang.org/x/oauth2 |
| Password Hashing | bcrypt (x/crypto v0.48.0) |
| UUID | google/uuid v1.6.0 |
| Rate Limiting | x/time v0.5.0 |
| Env Loading | godotenv v1.5.1 |
| Validation | go-playground/validator v10.30.1 |
| CORS | gin-contrib/cors |
| Containerization | Docker + Docker Compose |
| Frontend | Vite + React (separate SPA) |

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
- **Middleware stack** — rate limiter → body size limit → JWT auth → RBAC (project-level)

## Project Structure

```
├── cmd/app/main.go                 — Application entry point
├── docker/
│   ├── Dockerfile                  — Multi-stage Go build
│   ├── Dockerfile.frontend         — Node/Vite build
│   └── docker-compose.yml          — App + PostgreSQL 16 + Frontend
├── internal/
│   ├── api/router.go               — Route definitions and wiring
│   ├── config/
│   │   ├── config.go               — Environment-based configuration
│   │   ├── google_conf.go          — Google OAuth config
│   │   └── yandex_conf.go          — Yandex OAuth config
│   ├── ctxkeys/ctxkeys.go          — Context key constants
│   ├── db/
│   │   ├── db.go                   — GORM connection + auto-migration
│   │   └── schema.sql              — Reference SQL schema
│   ├── middleware/
│   │   ├── auth.go                 — JWT authentication
│   │   ├── rbac.go                 — Project-level RBAC (owner/member/read_only)
│   │   ├── ratelimit.go            — IP-based rate limiter
│   │   └── requestbody.go          — Request body size limit
│   ├── models/                     — GORM model definitions
│   │   ├── user.go
│   │   ├── user_provider.go        — OAuth provider links
│   │   ├── auth.go                 — RefreshToken model
│   │   ├── project.go
│   │   ├── participant.go
│   │   ├── status.go
│   │   ├── priority.go
│   │   ├── task.go                 — Task + UpdateHistory
│   │   ├── comment.go
│   │   └── attachment.go
│   ├── core/
│   │   ├── auth/
│   │   │   ├── handler.go          — Register, Login, Refresh, Logout, Sessions, ChangePassword, 2FA stubs
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── model.go
│   │   │   ├── errors.go
│   │   │   └── oauth/
│   │   │       ├── handler.go      — Google/Yandex login & callback
│   │   │       ├── google_service.go
│   │   │       ├── yandex_service.go
│   │   │       └── ctxkeys.go
│   │   ├── user/                   — User CRUD, avatar, search, delete
│   │   ├── project/                — Project CRUD, search
│   │   ├── task/                   — Task CRUD, archive, history, search
│   │   ├── status/                 — Per-project statuses with reorder
│   │   ├── priority/               — Per-project priorities
│   │   ├── participant/            — Project participants & roles
│   │   ├── comment/                — Task comments
│   │   └── attachment/             — File uploads & downloads
│   ├── modules/
│   │   ├── notify/                 — Notifications (stub, commented out)
│   │   └── chat/chat.go            — Chat (placeholder)
│   └── version/
│       ├── version.go              — Version constant + changelog embed
│       └── changelog.md            — Release notes
├── pkg/
│   ├── cookie/                     — Cookie set/delete helpers
│   ├── events/                     — Event bus + event structs
│   ├── hash/                       — bcrypt + SHA-256 token hashing
│   ├── jwt/                        — Access/refresh JWT tokens
│   ├── nullable/                   — Nullable types for PATCH updates
│   ├── response/                   — Standardized JSON responses
│   └── validation/                 — Request validation helpers
├── frontend/                       — Vite + React SPA
├── uploads/                        — File upload directory (gitignored)
│   ├── avatars/
│   └── {task_id}/
├── .env.example                    — Environment template
└── .gitignore
```

## Getting Started

### Prerequisites

- Go 1.26+ (or Docker)
- PostgreSQL 16

### Local Development

```bash
# Configure environment
cp .env.example .env
# Edit .env with your database credentials, JWT secret, and OAuth keys

# Run
go mod tidy
go run ./cmd/app
```

### With Docker

```bash
docker compose -f docker/docker-compose.yml up --build
```

Server starts at `http://localhost:8080`. Frontend (when built) runs on port `3000`.

## Configuration

All configuration is via environment variables (see `.env.example`):

| Variable | Default | Description |
|---|---|---|
| `ENV` | `development` | Environment mode |
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
| `UPLOAD_PATH` | `./uploads` | File upload directory |
| `COOKIE_DOMAIN` | — | Cookie domain |
| `COOKIE_SECURE` | `false` | Cookie Secure flag |
| `COOKIE_SAMESITE` | `1` | SameSite (1=Default, 2=Lax, 3=Strict, 4=None) |
| `CORS_ORIGINS` | `http://localhost:5173` | Comma-separated CORS origins |
| `FRONTEND_URL` | `http://localhost:5173` | Frontend URL for OAuth redirects |
| `NOTIFY_MODULE` | `false` | Enable notifications module |
| `GOOGLE_CLIENT_ID` | — | Google OAuth client ID |
| `GOOGLE_CLIENT_SECRET` | — | Google OAuth client secret |
| `GOOGLE_REDIRECT_URL` | — | Google OAuth redirect URL |
| `YANDEX_CLIENT_ID` | — | Yandex OAuth client ID |
| `YANDEX_CLIENT_SECRET` | — | Yandex OAuth client secret |
| `YANDEX_REDIRECT_URL` | — | Yandex OAuth redirect URL |

## API Endpoints

### Health & Meta

| Method | Path | Auth |
|---|---|---|
| `GET` | `/health` | — |
| `GET` | `/api/v1/version` | — |
| `GET` | `/api/v1/changelog` | — |

### Auth (`/api/v1/auth`)

Rate-limited (5 req/min). Cookies used for token storage.

| Method | Path | Status |
|---|---|---|
| `POST` | `/register` | ✅ |
| `POST` | `/login` | ✅ |
| `POST` | `/refresh` | ✅ |
| `POST` | `/logout` | ✅ |
| `GET` | `/google/login` | ✅ |
| `GET` | `/google/callback` | ✅ |
| `GET` | `/yandex/login` | ✅ |
| `GET` | `/yandex/callback` | ✅ |
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
| `PUT` | `/me/avatar` | ✅ |
| `PUT` | `/me/change-password` | ✅ |
| `DELETE` | `/me` | ✅ |
| `GET` | `/me/sessions` | ✅ |
| `DELETE` | `/me/sessions/:id` | ✅ |
| `GET` | `/search?q=` | ✅ |

### Projects (`/api/v1/projects`)

Requires JWT auth. Participant roles: `owner` (read+write+delete), `member` (read+write), `read_only` (read).

| Method | Path | Access | Status |
|---|---|---|---|
| `GET` | `/` | authenticated | ✅ |
| `GET` | `/owned` | authenticated | ✅ |
| `GET` | `/search?q=` | authenticated | ✅ |
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
| `GET` | `/:id/attachments` | read_only+ | ✅ |

### Tasks (`/api/v1/projects/:id/tasks`)

Requires JWT auth + project access. Supports `?archived=true` filter. History tracked as JSONB diffs.

| Method | Path | Access | Status |
|---|---|---|---|
| `GET` | `/?archived=` | read_only+ | ✅ |
| `POST` | `/` | member+ | ✅ |
| `GET` | `/:task_id?archived=` | read_only+ | ✅ |
| `PUT` | `/:task_id?archived=` | member+ | ✅ |
| `DELETE` | `/:task_id` | owner | ✅ |
| `GET` | `/:task_id/history` | read_only+ | ✅ |
| `GET` | `/:task_id/comments` | member+ | ✅ |
| `POST` | `/:task_id/comments` | member+ | ✅ |
| `PUT` | `/:task_id/comments/:comment_id` | member+ | ✅ |
| `DELETE` | `/:task_id/comments/:comment_id` | member+ | ✅ |
| `GET` | `/:task_id/attachments` | read_only+ | ✅ |
| `GET` | `/:task_id/attachments/:attachment_id` | read_only+ | ✅ |
| `POST` | `/:task_id/attachments` | member+ | ✅ |
| `DELETE` | `/:task_id/attachments/:attachment_id` | member+ | ✅ |

### Global Task Search (`/api/v1/tasks`)

Requires JWT auth.

| Method | Path | Status |
|---|---|---|
| `GET` | `/me?type=assigned\|reported` | ✅ |
| `GET` | `/search?q=` | ✅ |

### Notifications (`/api/v1/notifications`)

Commented out in `main.go`. Not available by default.

## Database Schema

### Tables

| Table | Description |
|---|---|
| `users` | Core user accounts (uuid PK, email, password_hash, name, role, 2fa secret) |
| `user_providers` | OAuth provider links (Google, Yandex) |
| `refresh_tokens` | JWT refresh token storage with revocation tracking and reuse detection |
| `projects` | Project entities owned by a user |
| `project_participants` | Many-to-many with roles (owner/member/read_only) |
| `statuses` | Per-project task statuses with `order_index` for drag-reorder |
| `priorities` | Per-project task priorities |
| `tasks` | Tasks with references to project, status, priority, assignee, reporter |
| `update_history` | JSONB-based field-level change tracking |
| `comments` | Task-scoped comments |
| `attachments` | File metadata (filename, path, size, mime type) |

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

## License

MIT
