# RBAC & Authorization Module

## Overview

This project implements a complete **Role-Based Access Control (RBAC)** system following the project's Clean Architecture principles.

The authorization system allows:

- Managing roles
- Managing permissions
- Assigning roles to users
- Assigning permissions to roles
- Performing authorization checks
- Protecting API endpoints through middleware
- Seeding default RBAC data

---

# Architecture

The module follows the same application flow used throughout the project.

```
HTTP Request
      │
      ▼
Authorization Middleware
      │
      ▼
Application Service
      │
      ▼
Repository Interface
      │
      ▼
Repository Implementation
      │
      ▼
Database
```

The module is split into the following layers:

```
Database
    ↓

sqlc

    ↓

Domain
    entities
    repositories

    ↓

Infrastructure
    postgres repositories

    ↓

Application
    authorization service

    ↓

HTTP
    handlers
    middleware

    ↓

Routes
```

---

# Database Schema

The RBAC system introduces four primary concepts.

## Roles

Represents a collection of permissions.

Examples:

- Super Admin
- Admin
- User

---

## Permissions

Represents an action that can be performed.

Examples

```
roles.read
roles.create
roles.update
roles.delete

permissions.read
permissions.create

users.read
users.update

audit-logs.view

sessions.view
```

Permission names follow the convention

```
resource.action
```

Example

```
roles.create
users.update
sessions.view
```

---

## User Roles

Many-to-many relationship.

```
User
    │
    ├──────────────► Role
```

A user can have multiple roles.

---

## Role Permissions

Many-to-many relationship.

```
Role
    │
    ├────────────► Permission
```

A role can have many permissions.

---

# Repository Layer

Repository interfaces are defined in the Domain layer.

Examples

```
RoleRepository

PermissionRepository

AssignmentRepository
```

The Infrastructure layer implements these interfaces using sqlc-generated queries.

This keeps the Application layer independent of PostgreSQL.

---

# Application Layer

The Authorization service contains all RBAC business logic.

Responsibilities include

- Creating roles
- Updating roles
- Deleting roles
- Creating permissions
- Updating permissions
- Deleting permissions
- Assigning roles
- Removing roles
- Assigning permissions
- Removing permissions
- Listing assignments
- Authorization checks

---

# Authorization Methods

The service exposes several helper methods.

## Authorize

Checks for one permission.

```go
Authorize(ctx, userID, "roles.read")
```

---

## AuthorizeAny

User must have at least one permission.

```go
AuthorizeAny(
    ctx,
    userID,
    "roles.read",
    "roles.create",
)
```

---

## AuthorizeAll

User must have every permission.

```go
AuthorizeAll(
    ctx,
    userID,
    "roles.read",
    "roles.create",
)
```

---

## UserHasPermission

Returns a boolean.

```go
allowed, err := service.UserHasPermission(...)
```

---

## UserHasRole

Returns whether a user owns a role.

---

## RoleHasPermission

Checks whether a role contains a permission.

---

# HTTP Handlers

The authorization handler exposes CRUD endpoints for

## Roles

```
GET     /authorization/roles

GET     /authorization/roles/get/:id

POST    /authorization/roles/create

PUT     /authorization/roles/update/:id

DELETE  /authorization/roles/delete/:id
```

---

## Permissions

```
GET

POST

PUT

DELETE
```

---

## Assignments

```
Assign role to user

Remove role

Assign permission to role

Remove permission

List user roles

List user permissions

List role permissions
```

---

## Authorization Checks

```
UserHasRole

UserHasPermission

RoleHasPermission
```

---

# Authentication vs Authorization

Authentication verifies identity.

```
Who are you?
```

Authorization verifies permissions.

```
What are you allowed to do?
```

Authentication always runs before authorization.

```
Request

↓

JWT Authentication

↓

Email Verification

↓

Permission Middleware

↓

Handler
```

---

# Authorization Middleware

The middleware protects endpoints using permission names.

Example

```go
middleware.RequirePermission("roles.read")
```

The middleware

1. Reads the authenticated user from context.

2. Calls

```go
AuthorizationService.Authorize(...)
```

3. Returns

```
403 Forbidden
```

if the user lacks the required permission.

---

# Route Integration

Routes are grouped by authentication level.

```
Public

↓

Authenticated

↓

Verified

↓

Permission Protected
```

Example

```go
roles.GET(
    "/",
    authorization.RequirePermission("roles.read"),
    handler.GetRoles,
)
```

Each endpoint declares the permission it requires.

This keeps authorization close to the route definition.

---

# Default Roles

The system seeds the following roles.

## Super Admin

Has every permission.

Used only for administration.

---

## Admin

Receives administrative permissions.

Can be customized later.

---

## User

Receives only standard user permissions.

---

# Default Permissions

Permissions are seeded automatically.

Examples

```
roles.read

roles.create

roles.update

roles.delete

permissions.read

permissions.create

assignments.assign-role

audit-logs.view

login-history.view

sessions.view
```

The naming convention is

```
resource.action
```

---

# Seeder

The project includes database seeders.

```
permissions.go

roles.go

assignments.go

users.go

seed.go
```

The seeding process is idempotent.

Running it multiple times does not create duplicates.

Execution order

```
Permissions

↓

Roles

↓

Role Permissions

↓

Users

↓

User Roles
```

---

# Default Administrator

The seeder creates a verified administrator account.

The account is assigned the

```
Super Admin
```

role.

The administrator receives every permission through role inheritance.

---

# Running the Seeder

Using Go

```bash
go run ./cmd/seed
```

Using Task

```bash
task seed
```

---

# Security Considerations

The authorization layer follows several principles.

- Authorization occurs on every protected request.

- Permissions are never trusted from JWT claims.

- Permissions are always loaded from the database.

- Roles inherit permissions.

- Users receive permissions through assigned roles.

- System roles cannot be modified.

- Duplicate role and permission assignments are prevented.

---

# Testing Recommendations

Recommended integration tests include

- Create Role

- Update Role

- Delete Role

- Create Permission

- Assign Role

- Remove Role

- Assign Permission

- Remove Permission

- Authorization middleware

- Super Admin access

- Forbidden access

- Authenticated but unverified user

---

# Future Improvements

Possible enhancements include

- Permission caching

- Wildcard permissions

```
users.*

roles.*
```

- Multi-tenancy support

- Organization-level roles

- Hierarchical roles

- Temporary permissions

- Permission audit reports

- Administrative dashboard

---

# Summary

The RBAC module provides a complete authorization system that integrates with the existing authentication flow.

Implemented features include:

- Database schema
- sqlc queries
- Domain entities
- Repository interfaces
- PostgreSQL repository implementations
- Authorization service
- HTTP handlers
- Authentication middleware
- Authorization middleware
- Route integration
- Database seeders

The module is fully integrated with the project's Clean Architecture and is designed to be extensible for future authorization requirements.