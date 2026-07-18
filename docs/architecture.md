# Backend Architecture

## Overview

This project is a REST API written in Go following **Clean Architecture**, **Domain-Driven Design (DDD)** principles, and **Dependency Injection**.

The primary goals are:

- Separation of concerns
- Testability
- Maintainability
- Scalability
- Framework independence

The application is organized so that business logic never depends on infrastructure such as PostgreSQL, Gin, JWT, or SMTP.

---

# High-Level Architecture

```
                   Client
                      │
                      ▼
              HTTP (Gin Router)
                      │
                      ▼
                HTTP Middleware
      (Logging, Recovery, Auth, RBAC)
                      │
                      ▼
                 HTTP Handlers
                      │
                      ▼
            Application Services
             (Business Logic)
                      │
                      ▼
          Domain Repository Interfaces
                      │
                      ▼
     PostgreSQL Repository Implementations
                      │
                      ▼
                 PostgreSQL
```

Dependencies always point inward.

```
Infrastructure
        │
        ▼
Application
        │
        ▼
Domain
```

The Domain layer has no knowledge of HTTP, PostgreSQL, JWT, or external services.

---

# Project Structure

```
cmd/
    api/
    seed/

docs/

internal/

    application/

    bootstrap/

    domain/

    infrastructure/

    interfaces/

    shared/
```

---

# cmd/

Contains executable applications.

```
cmd/api
```

Starts the HTTP server.

```
cmd/seed
```

Seeds the database with default data.

Additional executables can be added later, such as:

```
cmd/worker

cmd/migrate

cmd/scheduler
```

---

# internal/application

Contains all business logic.

Examples

```
application/auth

application/user

application/security

application/authorization
```

Responsibilities

- Validation
- Business rules
- Transactions
- Authorization
- Authentication
- Password management

Application services communicate only through repository interfaces.

---

# internal/domain

Contains the core domain.

It defines

```
Entities

Repository interfaces

Value objects

Enums

Constants
```

No framework-specific code exists here.

The Domain layer does not know

- Gin
- PostgreSQL
- sqlc
- JWT
- SMTP

---

# internal/infrastructure

Contains external implementations.

Examples

```
postgres

jwt

email

storage

logger
```

Responsibilities

- PostgreSQL repositories
- SQLC generated code
- JWT generation
- Email delivery
- Local storage
- Logging

Infrastructure implements interfaces defined by the Domain layer.

---

# internal/interfaces

Contains delivery mechanisms.

Current implementation

```
HTTP
```

Inside

```
handlers

middleware

routes

request

response
```

Future additions could include

```
gRPC

GraphQL

CLI

WebSockets
```

without changing the business logic.

---

# internal/bootstrap

Responsible for application startup.

Responsibilities include

- Dependency injection
- Repository construction
- Service construction
- Handler registration
- Middleware registration

The bootstrap package wires the entire application together.

---

# Dependency Injection

Dependencies are created once during startup.

Example

```
Repository

↓

Service

↓

Handler

↓

Router
```

This avoids global variables and keeps components loosely coupled.

---

# Request Lifecycle

A request flows through the application as follows.

```
Client

↓

Gin Router

↓

Request Middleware

↓

Authentication

↓

Authorization

↓

Handler

↓

Application Service

↓

Repository

↓

Database

↓

Response
```

Each layer has a single responsibility.

---

# Authentication Flow

Authentication is JWT-based.

Flow

```
Login

↓

Verify credentials

↓

Generate Access Token

↓

Generate Refresh Token

↓

Create Session

↓

Return Tokens
```

Every authenticated request

```
Bearer Token

↓

JWT Validation

↓

Session Validation

↓

Authenticated User

↓

Request Context
```

The authenticated user and session are stored inside the Gin context.

---

# Authorization (RBAC)

Authorization is role-based.

```
User

↓

Role

↓

Permission
```

Users inherit permissions from assigned roles.

Protected routes use middleware.

```
RequirePermission("roles.read")
```

The middleware verifies that the authenticated user possesses the required permission before allowing the request to continue.

---

# Middleware Stack

The application currently uses middleware for

```
Request ID

↓

Logging

↓

Recovery

↓

CORS

↓

Rate Limiting

↓

Authentication

↓

Email Verification

↓

Authorization
```

Middleware order is important because each middleware builds upon the previous one.

---

# Database Layer

The application uses

```
PostgreSQL

↓

sqlc

↓

Repositories
```

Repository implementations wrap sqlc-generated queries and expose domain-friendly methods to the application layer.

---

# Transactions

Complex operations execute inside transactions.

Example

```
Register User

↓

Create User

↓

Generate Verification Token

↓

Create Session

↓

Commit
```

If any step fails, the transaction is rolled back.

---

# Security

Security features include

- Argon2 password hashing
- JWT access tokens
- Refresh tokens
- Session tracking
- Audit logs
- Login history
- Rate limiting
- Account verification
- Password reset
- RBAC authorization

Passwords are never stored in plain text.

---

# Session Management

Every successful login creates a session.

Each session stores

- Device
- Browser
- Platform
- IP Address
- Refresh Token
- Expiration

This enables

- Logout
- Logout everywhere
- Active session management

---

# Audit Logging

Security-sensitive actions generate audit logs.

Examples

```
USER_LOGIN

USER_LOGOUT

PASSWORD_CHANGED

EMAIL_VERIFIED

ROLE_ASSIGNED

PERMISSION_ASSIGNED
```

Audit logs provide traceability for administrative actions.

---

# Login History

Login history records

- Successful logins
- Failed logins
- Device information
- Login time
- Logout time
- IP address

This improves security monitoring and user visibility into account activity.

---

# Validation

Incoming requests are validated before reaching the application layer.

Validation includes

- Required fields
- UUID parsing
- Email validation
- Password rules
- Pagination parameters

Handlers remain thin because validation is centralized.

---

# Error Handling

Errors originate in the application layer.

Handlers translate them into consistent HTTP responses.

Example

```
Application Error

↓

Response Helper

↓

JSON Response
```

This keeps error formatting consistent across the API.

---

# Rate Limiting

The API includes configurable rate limiting.

Different route groups may have different limits.

Examples

```
Public

Authenticated

Verified
```

Rate limiting protects against abuse while allowing flexibility for authenticated users.

---

# Seeding

The project provides database seeders.

```
Permissions

↓

Roles

↓

Role Assignments

↓

Users
```

Seeders are idempotent, meaning they can be executed multiple times without creating duplicate records.

Run using

```bash
task seed
```

or

```bash
go run ./cmd/seed
```

---

# Development Workflow

Typical development process

```
Create Migration

↓

Run Migration

↓

Update SQL Queries

↓

Generate sqlc Code

↓

Implement Repository

↓

Implement Service

↓

Implement Handler

↓

Register Routes

↓

Write Tests

↓

Seed Data
```

Useful commands

```bash
task migrate:create
task migrate:up
task sqlc
task seed
task run
task test
```

---

# Design Principles

This project follows several important principles.

## Single Responsibility Principle

Each package has one responsibility.

---

## Dependency Inversion

Business logic depends on interfaces instead of implementations.

---

## Clean Architecture

Outer layers depend on inner layers.

Never the reverse.

---

## Explicit Dependency Injection

No global state.

Dependencies are passed through constructors.

---

## Thin Handlers

Handlers coordinate requests.

Business logic belongs inside services.

---

## Repository Pattern

Repositories abstract data persistence.

Changing PostgreSQL to another database requires only infrastructure changes.

---

## Idempotent Seeders

Seeders can be executed repeatedly without duplicate data.

---

# Future Improvements

Potential enhancements include

- Redis caching
- Background workers
- Event bus
- WebSockets
- Multi-tenancy
- OpenTelemetry tracing
- Prometheus metrics
- Docker Compose
- Kubernetes deployment
- CI/CD pipelines
- Integration tests
- API versioning
- Feature flags

---

# Summary

The project is organized around Clean Architecture and Domain-Driven Design principles.

Business logic remains independent of frameworks and infrastructure, making the application easier to maintain, extend, and test.

The architecture emphasizes:

- Clear separation of concerns
- Dependency injection
- Testability
- Security
- Scalability
- Maintainability

Each feature follows a consistent implementation flow:

```
Migration
    ↓
SQLC Queries
    ↓
Domain
    ↓
Repository
    ↓
Application Service
    ↓
Handler
    ↓
Middleware
    ↓
Routes
    ↓
Tests
    ↓
Seed Data
```

Following this workflow ensures every module is implemented consistently across the project.