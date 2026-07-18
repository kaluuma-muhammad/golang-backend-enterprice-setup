# Database Design

## Overview

The application uses **PostgreSQL** as its primary database.

Database access is implemented using:

- PostgreSQL
- pgx
- sqlc
- Repository Pattern
- Database Migrations (Goose)

Every database interaction passes through repository implementations, ensuring that the application layer never depends directly on SQL.

---

# Architecture

```
Application Service
        │
        ▼
Repository Interface
        │
        ▼
Repository Implementation
        │
        ▼
sqlc Generated Queries
        │
        ▼
PostgreSQL
```

Business logic never executes SQL directly.

---

# Database Migrations

Schema changes are managed using Goose migrations.

Create a migration

```bash
task migrate:create create_users_table
```

Apply migrations

```bash
task migrate:up
```

Rollback

```bash
task migrate:down
```

Migration status

```bash
task migrate:status
```

Reset database

```bash
task migrate:reset
```

---

# SQLC

SQL queries are written manually inside

```
internal/infrastructure/postgres/queries/
```

Example

```
users.sql

sessions.sql

roles.sql

permissions.sql
```

Generate Go code

```bash
task sqlc
```

sqlc generates

```
internal/infrastructure/postgres/sqlc/
```

Generated files should never be edited manually.

---

# Repository Pattern

Repositories isolate database access.

Example

```
Domain

UserRepository

↓

Infrastructure

PostgresUserRepository
```

Application services only know about interfaces.

---

# Main Tables

## users

Stores application users.

Important columns

- id
- email
- password
- first_name
- last_name
- avatar
- is_verified
- failed_login_attempts
- locked_until
- last_login
- created_at
- updated_at

---

## sessions

Stores login sessions.

Contains

- refresh token hash
- browser
- platform
- device
- IP address
- expiration
- last used timestamp

Supports

- logout
- logout all devices
- session management

---

## tokens

Stores temporary verification tokens.

Examples

- email verification
- password reset
- account activation

---

## audit_logs

Tracks security-sensitive actions.

Examples

- login
- logout
- password change
- role assignment
- permission assignment

---

## login_history

Stores login attempts.

Contains

- success/failure
- login time
- logout time
- IP
- device
- browser

---

## roles

RBAC roles.

Examples

- Super Admin
- Admin
- Manager
- User

---

## permissions

Individual permissions.

Examples

```
roles.read

roles.create

users.update

permissions.delete
```

---

## user_roles

Many-to-many relationship.

```
User

↓

Role
```

A user can have multiple roles.

---

## role_permissions

Many-to-many relationship.

```
Role

↓

Permission
```

A role can have many permissions.

---

# Relationships

```
User
 │
 ├─────────────┐
 │             │
 ▼             ▼
Sessions   User Roles
               │
               ▼
             Roles
               │
               ▼
        Role Permissions
               │
               ▼
          Permissions
```

---

# Transactions

Complex operations execute inside transactions.

Examples

Registration

```
Create User

↓

Create Session

↓

Generate Verification Token

↓

Commit
```

If one step fails, the transaction rolls back.

---

# UUIDs

Every primary key uses UUID.

Advantages

- globally unique
- safer public identifiers
- easier distributed systems

---

# Soft Deletes

Current implementation uses hard deletes.

Future improvement

```
deleted_at TIMESTAMP NULL
```

for soft deletion.

---

# Seeding

Database seeders populate default data.

Seed order

```
Permissions

↓

Roles

↓

Assignments

↓

Users
```

Run

```bash
task seed
```

Seeders are idempotent and safe to run multiple times.

---

# Naming Conventions

Tables

```
snake_case
plural
```

Columns

```
snake_case
```

Primary Keys

```
id UUID
```

Foreign Keys

```
user_id

role_id

permission_id
```

Indexes should exist on

- email
- refresh_token
- user_id
- role_id
- permission_id

---

# Development Workflow

Typical workflow

```
Create Migration

↓

Write SQL

↓

Generate sqlc

↓

Implement Repository

↓

Implement Service

↓

Implement Handler

↓

Register Route

↓

Seed Database

↓

Write Tests
```

---

# Summary

The database layer follows a clean separation of concerns.

```
Migration
        ↓
SQL
        ↓
sqlc
        ↓
Repository
        ↓
Application
```

This architecture provides:

- type-safe SQL
- maintainable repositories
- clean business logic
- easy testing
- scalable schema evolution