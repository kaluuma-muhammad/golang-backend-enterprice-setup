# API Documentation

## Overview

The API is a RESTful JSON API built using Gin.

Base URL

```
/api/v1
```

All responses are JSON.

Authentication uses JWT Bearer Tokens.

---

# Authentication

Authenticated endpoints require

```
Authorization: Bearer <access_token>
```

---

# Route Groups

The API is divided into three groups.

## Public

No authentication required.

```
/health

/auth/*
```

---

## Protected

Authentication required.

```
/auth/activate-account

/auth/logout

/auth/logout-all
```

---

## Verified

Authentication and verified email required.

```
/user/*

/authorization/*
```

---

# Public Endpoints

## Health

```
GET /health
```

Returns API health status.

---

## Register

```
POST /auth/register
```

Creates a new account.

---

## Login

```
POST /auth/login
```

Authenticates a user.

Returns

- Access Token
- Refresh Token

---

## Refresh Token

```
POST /auth/refresh
```

Generates a new access token.

---

## Forgot Password

```
POST /auth/forgot-password
```

Sends password reset instructions.

---

## Verify Email

```
POST /auth/verify-email
```

Verifies the email verification code.

---

## Reset Password

```
POST /auth/reset-password
```

Resets the user password.

---

## Resend Verification

```
POST /auth/resend-verification-code
```

Sends another verification code.

---

# Protected Endpoints

Require authentication.

---

## Activate Account

```
POST /auth/activate-account
```

---

## Logout

```
POST /auth/logout
```

Logs out the current session.

---

## Logout All Sessions

```
POST /auth/logout-all
```

Terminates every active session.

---

# User Endpoints

Require

- Authentication
- Verified Email

---

## Current User

```
GET /user/me
```

---

## Update Profile

```
PUT /user/update-account
```

---

## Update Password

```
PUT /user/update-password
```

---

## Update Avatar

```
PUT /user/update-avatar
```

---

## Audit Logs

```
GET /user/audit-logs
```

---

## Login History

```
GET /user/login-history
```

---

## User Sessions

```
GET /user/sessions
```

---

## Current Sessions

```
GET /user/sessions/current
```

---

# Authorization API

All endpoints require authentication, email verification, and the corresponding RBAC permission.

---

## Roles

| Method | Endpoint | Permission |
|---------|----------|------------|
| GET | `/authorization/roles` | `roles.read` |
| GET | `/authorization/roles/get/:id` | `roles.read` |
| POST | `/authorization/roles/create` | `roles.create` |
| PUT | `/authorization/roles/update/:id` | `roles.update` |
| DELETE | `/authorization/roles/delete/:id` | `roles.delete` |

---

## Permissions

| Method | Endpoint | Permission |
|---------|----------|------------|
| GET | `/authorization/permissions` | `permissions.read` |
| GET | `/authorization/permissions/get/:id` | `permissions.read` |
| POST | `/authorization/permissions/create` | `permissions.create` |
| PUT | `/authorization/permissions/update/:id` | `permissions.update` |
| DELETE | `/authorization/permissions/delete/:id` | `permissions.delete` |

---

## Role Assignments

| Method | Endpoint | Permission |
|---------|----------|------------|
| POST | `/authorization/assignments/assign-role` | `users.assign-role` |
| POST | `/authorization/assignments/remove-role` | `users.unassign-role` |
| POST | `/authorization/assignments/assign-permission` | `roles.assign-permission` |
| POST | `/authorization/assignments/remove-permission` | `roles.unassign-permission` |
| GET | `/authorization/assignments/list-user-roles/:userId` | `users.read` |
| GET | `/authorization/assignments/list-user-permissions/:userId` | `users.read` |
| GET | `/authorization/assignments/list-role-permissions/:roleId` | `roles.read` |

---

## Authorization Checks

| Method | Endpoint | Permission |
|---------|----------|------------|
| GET | `/authorization/checks/user-has-role/:userId/:roleId` | `users.read` |
| GET | `/authorization/checks/user-has-permission/:userId/:permission` | `users.read` |
| GET | `/authorization/checks/role-has-permission/:roleId/:permissionId` | `roles.read` |

---

# Middleware Pipeline

Every request passes through middleware.

```
Request

↓

Request ID

↓

Logger

↓

Recovery

↓

CORS

↓

Rate Limiter

↓

Authentication (protected)

↓

Email Verification

↓

Permission Authorization

↓

Handler
```

---

# Response Format

Successful responses return JSON.

Example

```json
{
  "success": true,
  "data": {}
}
```

Errors return a consistent structure.

```json
{
  "success": false,
  "message": "Permission denied"
}
```

---

# HTTP Status Codes

| Status | Meaning |
|---------|---------|
| 200 | OK |
| 201 | Created |
| 204 | No Content |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 409 | Conflict |
| 422 | Validation Error |
| 429 | Too Many Requests |
| 500 | Internal Server Error |

---

# Rate Limiting

Routes are grouped into different rate-limiting policies.

- Public
- Protected
- Verified

Each group can have independent limits configured in the rate limiter.

---

# API Versioning

Current version

```
/api/v1
```

Future versions should follow the same convention.

```
/api/v2
```

---

# Future Improvements

- Swagger / OpenAPI generation
- API examples
- Postman collection
- Request/response schemas
- Error code reference
- Pagination guide
- Filtering and sorting conventions

---

# Summary

The API is organized into logical route groups with layered security:

```
Public
        ↓
Protected
        ↓
Verified
        ↓
Permission Protected
```

Every request follows the same lifecycle, ensuring consistent authentication, authorization, validation, rate limiting, and JSON responses across the entire application.