# Backend Platform

A production-ready, enterprise-grade backend platform written in Go, following **Clean Architecture**, **Domain-Driven Design (DDD)** principles, and a modular design that supports scalable business applications.

The project is designed to serve as a reusable foundation for enterprise APIs with authentication, authorization, session management, email verification, and future modules such as RBAC, OAuth, notifications, and background workers.

---

# Quick Start

## Prerequisites

- Go 1.26+
- PostgreSQL 17+
- Docker & Docker Compose (recommended)
- Task

Clone the repository:

```bash
git clone https://github.com/kaluuma-muhammad/golang-backend-enterprice-setup.git
cd go-api
```

Create your environment file:

```bash
cp .env.example .env
```

### Run with Docker (Recommended)

Development:

```bash
task docker:dev
```

Production:

```bash
task docker:prod
```

### Run Locally

Install dependencies:

```bash
task tidy
```

Run database migrations:

```bash
task migrate:up
```

Seed the database:

```bash
task seed
```

Start the API:

```bash
task run
```

For detailed Docker documentation, see:

```text
docs/docker.md
```

---

# Technology Stack

## Language

* **Go (Golang)** — Primary programming language.

---

## HTTP Framework

* **Gin** — High-performance HTTP router and middleware framework.

Used for:

* Routing
* Middleware
* Request binding
* Response handling

---

## Database

* **PostgreSQL**

Primary relational database used for persistent storage.

Current modules:

* Users
* Sessions
* Verification Tokens

Future modules:

* Roles
* Permissions
* Audit Logs
* Notifications

---

## Database Driver

* **pgx**

PostgreSQL driver and connection pool.

Advantages:

* High performance
* Native PostgreSQL support
* Context-aware queries
* Connection pooling

---

## Caching & Distributed Infrastructure

* **Redis**

Redis is integrated as the application's distributed cache and infrastructure layer.

Current uses:

* Distributed rate limiting
* Cache abstraction

The architecture supports switching between in-memory and Redis-backed implementations without changing application code.

Future modules will reuse the same Redis infrastructure for:

* Permission caching (RBAC)
* User profile caching
* JWT blacklist
* Password reset tokens
* Email verification tokens
* API response caching
* Background jobs
* Distributed locking

---

The application also supports environment-based infrastructure selection.

For example, the rate limiter can be configured to use either an in-memory implementation or Redis without requiring code changes.

```env
RATE_LIMIT_STORE=memory
```

or

```env
RATE_LIMIT_STORE=redis
```

This makes local development lightweight while enabling distributed rate limiting in production.

---

## SQL Query Generation

* **sqlc**

Instead of writing ORM models, SQL queries are written manually and converted into type-safe Go code.

Benefits:

* Compile-time query validation
* Excellent performance
* Full SQL control
* Strong typing

Query files are located in:

```text
internal/infrastructure/postgres/queries/
```

Generated code:

```text
internal/infrastructure/postgres/generated/
```

Whenever SQL changes, regenerate the code:

```bash
task sqlc
```

---

## Database Migrations

* **Goose**

Database schema changes are managed through versioned SQL migrations.

Migration files:

```text
internal/infrastructure/postgres/migrations/
```

Create a migration:

```bash
task migrate:create -- create_users
```

Apply migrations:

```bash
task migrate:up
```

Rollback:

```bash
task migrate:down
```

Migration status:

```bash
task migrate:status
```

Reset database:

```bash
task migrate:reset
```

---

## Authentication

Authentication is implemented using:

* JWT Access Tokens
* Rotating Refresh Tokens
* Session Management
* Email Verification

Passwords are hashed using:

* argon2 package

Refresh tokens and verification tokens are hashed using:

* SHA-256

---

## JWT

Library:

```text
github.com/golang-jwt/jwt/v5
```

Access tokens contain:

* User ID
* Session ID
* Issued At
* Expiration

Authentication is stateless while refresh tokens remain stateful through the database.

---

## Validation

Validation uses:

```text
go-playground/validator
```

Used for validating DTOs before business logic executes.

Example:

```go
Email string `validate:"required,email"`
```

---

## Logging

Structured logging is implemented using:

* Uber Zap

Request logging middleware records:

* HTTP method
* URL
* Response status
* Request duration

---

## Configuration

Configuration is loaded from environment variables.

Environment variables are managed using:

```text
godotenv
```

Configuration is centralized under:

```text
internal/shared/config
```

---

## UUID Generation

Universally Unique Identifiers are generated using:

```text
github.com/google/uuid
```

All primary keys use UUIDs instead of auto-incrementing integers.

---

## Email

Email delivery is implemented through a provider abstraction.

Current provider:

* SMTP

Architecture allows replacing SMTP with:

* Resend
* Mailgun
* Amazon SES
* SendGrid

without changing business logic.

---

## Password Hashing

Passwords are hashed using bcrypt.

Plain-text passwords are never stored.

---

## Project Architecture

The project follows a layered architecture.

```text
cmd/

deployments/
├── docker/
└── kubernetes/

docs/

internal/
├── application/
├── domain/
├── infrastructure/
├── interfaces/
└── shared/
```

### Domain

Contains:

* Entities
* Repository interfaces
* Domain errors

No external dependencies.

---

### Application

Contains business logic.

Examples:

* Register
* Login
* Refresh
* Logout
* Verify Email

Application services coordinate repositories and infrastructure services.

---

### Infrastructure

Contains implementations.

Examples:

* PostgreSQL
* JWT
* Email
* Logging

Nothing inside infrastructure contains business rules.

---

### Interfaces

Contains delivery mechanisms.

Currently:

* HTTP API

Future:

* gRPC
* GraphQL
* Message Queues

---

### Shared

Contains reusable utilities.

Examples:

* Configuration
* Validation
* Cache Abstractions
* Rate Limiting
* Constants
* Utilities
* Shared Errors

---

# Development Workflow

The project supports two development workflows.

## Local Development

### Install Dependencies

```bash
go mod tidy
```
---

### Run the API

```bash
task run
```
---

### Seed the Database

```bash
task seed
```
---

### Build

```bash
task build
```
---

### Run Tests

```bash
task test
```

---

### Format Code

```bash
task fmt
```

---

### Run Static Analysis

```bash
task vet
```

---

### Generate SQL Code

Whenever SQL queries change:

```bash
task sqlc
```

---


## Docker Development Environment

### Start the Docker development environment

```bash
task docker:dev
```
---


### Build and start the Docker development environment

```bash
task docker:dev:build
```
---

## Start Production Environment

### Start the Docker production environment

```bash
task docker:prod
```
---

### Build and start the Docker production environment

```bash
task docker:prod:build
```
---

### Stop all Docker containers

```bash
task docker:down
```
---

### View Docker logs

```bash
task docker:logs
```
---

### Rebuild and restart the development environment

```bash
task docker:dev:restart
```
---

### Rebuild and restart the production environment

```bash
task docker:prod:restart
```
---

### Stop containers and remove volumes

```bash
task docker:clean
```
---

### Show running Docker containers

```bash
task docker:ps
```

---


## Migration Workflow

Create a migration:

```bash
task migrate:create -- migration_name
```

Apply migrations:

```bash
task migrate:up
```

Rollback:

```bash
task migrate:down
```

Check migration status:

```bash
task migrate:status
```

Reset the database:

```bash
task migrate:reset
```

---

# Environment Variables

Example:

```env
APP_NAME=Go API
APP_PORT=8080
APP_BASE_URL=http://localhost:8080

DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_SSLMODE=disable

JWT_SECRET=super-secret-key
JWT_ISSUER=backend-platform
JWT_ACCESS_TOKEN_MINUTES=10080
JWT_REFRESH_TOKEN_DAYS=30

EMAIL_VERIFICATION_MINUTES=15
PASSWORD_RESET_MINUTES=15
MAGIC_LINK_MINUTES=10

MAIL_MAILER=smtp
MAIL_HOST=
MAIL_PORT=
MAIL_USERNAME=
MAIL_PASSWORD=
MAIL_ENCRYPTION=ssl
MAIL_FROM_ADDRESS=
MAIL_FROM_NAME="Go API"

SEED_ADMIN_EMAIL=admin@example.com
SEED_ADMIN_PASSWORD=Admin@12345
SEED_ADMIN_FIRST_NAME=System
SEED_ADMIN_LAST_NAME=Administrator

RATE_LIMIT_STORE=memory  # local developement
# RATE_LIMIT_STORE=redis  # deploying with Docker

REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

---

# Current Features

* User Registration
* User Login
* JWT Authentication
* Refresh Token Rotation
* Session Management
* Logout
* Logout All Devices (planned)
* Email Verification
* Resend Verification
* Forgot Password
* Password Reset
* Authentication Middleware
* Email Verification Middleware
* Device Management
* Login History
* Audit Logs
* Account Lockout
* Structured Logging
* PostgreSQL Integration
* SQLC Query Generation
* Goose Database Migrations
* Rate Limiting
* Role-Based Access Control (RBAC)
* Permission-Based Authorization
* Docker Development Environment
* Multi-stage Docker Builds
* Docker Compose
* Air Hot Reload
* pgAdmin Integration
* Automatic Database Migrations
* Automatic Database Seeding
* Redis Caching Module

---

# Planned Features

* OAuth (Google, GitHub, Microsoft)
* Background Workers
* MinIO Object Storage
* Mailpit Development Mail Server
* Kubernetes Deployment
* OpenTelemetry Metrics & Tracing
* CI/CD Pipelines

---

# Documentation

Project documentation is available under the `docs/` directory.

- `docs/architecture.md`
- `docs/api.md`
- `docs/authentication.md`
- `docs/database.md`
- `docs/docker.md`

---

# Contributing

1. Fork the repository.
2. Create a feature branch.
3. Make your changes.
4. Run formatting and static analysis:

```bash
task fmt
task vet
task test
```

5. If database queries change:

```bash
task sqlc
```

6. If the schema changes:

* Create a Goose migration.
* Apply the migration locally.
* Commit both the migration and regenerated SQLC code.

7. Open a Pull Request with a clear description of the change.

---

# License

This project is intended as a reusable enterprise backend foundation and may be licensed according to the project's chosen open-source or proprietary license.
