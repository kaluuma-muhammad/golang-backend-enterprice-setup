# Backend Platform

A production-ready, enterprise-grade backend platform written in Go, following **Clean Architecture**, **Domain-Driven Design (DDD)** principles, and a modular design that supports scalable business applications.

The project is designed to serve as a reusable foundation for enterprise APIs with authentication, authorization, session management, email verification, and future modules such as RBAC, OAuth, notifications, and background workers.

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
│
internal/
│
├── domain/
├── application/
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
* Constants
* Utilities
* Shared Errors

---

# Development Workflow

## Install Dependencies

```bash
go mod tidy
```

---

## Run the API

```bash
task run
```

---

## Build

```bash
task build
```

---

## Run Tests

```bash
task test
```

---

## Format Code

```bash
task fmt
```

---

## Run Static Analysis

```bash
task vet
```

---

## Generate SQL Code

Whenever SQL queries change:

```bash
task sqlc
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
APP_NAME=Backend Platform
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=backend_platform
DB_SSLMODE=disable

JWT_SECRET=your-secret-key

SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=user@example.com
SMTP_PASSWORD=password
SMTP_FROM=no-reply@example.com
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

---

# Planned Features

* OAuth (Google, GitHub, Microsoft)
* Role-Based Access Control (RBAC)
* Permission-Based Authorization
* Background Workers
* Redis Caching
* Rate Limiting
* Multi-Factor Authentication (MFA)
* Kubernetes Deployment
* OpenTelemetry Metrics & Tracing

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
