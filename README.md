# Nexa Task Tracker

Task tracking system built with Go, Gin, PostgreSQL, and GORM.

## Features

- User authentication with JWT (access + refresh tokens)
- 2FA support (TOTP)
- Project management with participants and roles
- Task management with statuses, priorities, and assignments
- Comments and attachments
- Update history tracking
- Notifications system

## Tech Stack

- **Backend**: Go 1.22, Gin
- **Database**: PostgreSQL 16
- **ORM**: GORM
- **Authentication**: JWT, TOTP (2FA)
- **Containerization**: Docker, Docker Compose

## Project Structure

```
/cmd/app                    - Application entry point
/internal
  /api                      - Router and API setup
  /config                   - Configuration management
  /db                       - Database connection and migrations
  /core                     - Core business logic
    /auth                   - Authentication (login, refresh, 2FA)
    /user                   - User management
    /project                - Project management
    /task                   - Task management
    /status                 - Task statuses
    /priority               - Task priorities
    /comment                - Task comments
    /participant            - Project participants
    /history                - Update history
    /attachment             - File attachments
  /middleware               - Auth, RBAC, rate limiting
  /pkg                      - Shared utilities (JWT, hash, response)
  /modules                  - Optional modules
    /notify                 - Notifications
/docker                     - Docker configuration
```

## Getting Started

### Prerequisites

- Go 1.22+
- PostgreSQL 16+
- Docker & Docker Compose (optional)

### Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and configure
3. Run with Docker Compose:

```bash
docker-compose -f docker/docker-compose.yml up
```

Or run locally:

```bash
go mod tidy
go run cmd/app/main.go
```

### API Endpoints

```
/health                                 - Health check

/api/v1/auth
  POST /register                        - Register new user
  POST /login                           - Login
  POST /refresh                         - Refresh access token
  POST /logout                          - Logout
  POST /2fa/setup                       - Setup 2FA
  POST /2fa/verify                      - Verify 2FA code
  POST /2fa/enable                      - Enable 2FA
  POST /2fa/disable                     - Disable 2FA

/api/v1/users
  GET  /me                              - Get current user
  PUT  /me                              - Update current user
  DELETE /me                            - Delete current user

/api/v1/projects
  GET    /                              - List projects
  POST   /                              - Create project
  GET    /:id                           - Get project
  PUT    /:id                           - Update project
  DELETE /:id                           - Delete project
  GET    /:id/participants              - Get participants
  POST   /:id/participants              - Add participant
  PUT    /:id/participants/:user_id     - Update participant role
  DELETE /:id/participants/:user_id     - Remove participant
  GET    /:id/statuses                  - Get statuses
  POST   /:id/statuses                  - Create status
  PUT    /:id/statuses/:status_id       - Update status
  DELETE /:id/statuses/:status_id       - Delete status
  GET    /:id/priorities                - Get priorities
  POST   /:id/priorities                - Create priority
  PUT    /:id/priorities/:priority_id   - Update priority
  DELETE /:id/priorities/:priority_id   - Delete priority

/api/v1/tasks
  GET    /                              - List tasks (with filters)
  POST   /                              - Create task
  GET    /:id                           - Get task
  PUT    /:id                           - Update task
  DELETE /:id                           - Delete task
  GET    /:id/history                   - Get task history
  GET    /:id/comments                  - Get comments
  POST   /:id/comments                  - Create comment
  PUT    /:id/comments/:comment_id      - Update comment
  DELETE /:id/comments/:comment_id      - Delete comment
  GET    /:id/attachments               - Get attachments
  POST   /:id/attachments               - Upload attachment
  GET    /:id/attachments/:attachment_id - Download attachment
  DELETE /:id/attachments/:attachment_id - Delete attachment

/api/v1/notifications
  GET  /                                - List notifications
  PUT  /:id/read                        - Mark as read
  PUT  /read-all                        - Mark all as read
```

## Database Schema

See `internal/db/schema.sql` for the complete database schema.

## Development

All handlers, services, and repositories are scaffolded with TODO comments for implementation.

Key features to implement:
- JWT token generation and validation
- Password hashing with bcrypt
- TOTP 2FA setup and verification
- GORM repository implementations
- Business logic in services
- Request validation in handlers
- GORM hooks for update history
- File upload/download for attachments

## License

MIT
