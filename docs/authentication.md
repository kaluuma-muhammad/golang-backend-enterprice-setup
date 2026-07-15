# Authentication Module Architecture

## Overview

The authentication module handles:

* User authentication
* Session management
* JWT-based access control
* Refresh token rotation
* Email verification
* Login tracking and auditing
* Device-aware sessions

The system follows **Clean Architecture (DDD-inspired)** with clear separation between:

* Domain
* Application
* Infrastructure
* Interfaces

---

## Architecture

```
HTTP Request
      │
      ▼
HTTP Handler
      │
      ▼
Application Service (Auth, Token, Security)
      │
      ▼
Domain Interfaces
      │
      ▼
Infrastructure (Postgres, JWT, Email, Logger)
```

---

# Database Design

## Users

```
users
──────────────────────────────
id
email
password
first_name
last_name
is_verified
created_at
updated_at
```

### Notes

* Passwords are hashed using bcrypt
* `is_verified` is used by authorization middleware
* No tokens or sessions are stored here

---

## Sessions

```
sessions
──────────────────────────────
id
user_id
refresh_token (hashed)
user_agent
ip_address
device_name
last_used_at
expires_at
revoked_at
revoked_reason
created_at
updated_at
```

### Purpose

Each login creates a session.

Sessions enable:

* Multi-device login
* Token rotation
* Logout per device
* Logout all devices
* Session-level security

---

## Tokens (Temporary)

```
tokens
──────────────────────────────
id
user_id
type
token_hash
expires_at
used_at
created_at
```

### Supported Types

* email_verification

Future-ready for:

* password_reset
* email_change
* magic_link

---

## Login History (NEW)

```
login_histories
──────────────────────────────
id
user_id
ip_address
user_agent
status
created_at
```

### Purpose

Tracks every login attempt:

* Success
* Failure

Used for:

* Security monitoring
* Fraud detection
* Audit trail

---

## Audit Logs (NEW)

```
audit_logs
──────────────────────────────
id
user_id
action
entity_type
created_at
```

### Purpose

Tracks important system actions:

* USER_LOGIN
* USER_LOGOUT
* EMAIL_VERIFIED

---

# Authentication Workflow

## Registration

```
POST /auth/register
```

Flow:

```
Validate Request
      │
Check Email Exists
      │
Hash Password (Security Service)
      │
Create User
      │
Create Session
      │
Generate Access Token (JWT)
      │
Generate Verification Code
      │
Store Token (hashed)
      │
Send Email
      │
Return Tokens
```

User is authenticated immediately but not verified.

---

## Login

```
POST /auth/login
```

Flow:

```
Find User
      │
Verify Password (bcrypt)
      │
Create Session (with device info)
      │
Generate Access Token
      │
Store Login History
      │
Create Audit Log
      │
Return Tokens
```

---

## Device Awareness (NEW)

Each session includes:

* IP address
* User agent
* Device name

Extracted using:

```
auth/device_info.go
```

---

## Refresh Token Rotation

```
POST /auth/refresh
```

Flow:

```
Receive Refresh Token
      │
Hash Token
      │
Find Session
      │
Validate:
   - Not expired
   - Not revoked
      │
Generate New Refresh Token
      │
Update Session
      │
Generate New Access Token
      │
Return Both
```

Old refresh tokens become invalid immediately.

---

## Logout

```
POST /auth/logout
```

Flow:

```
Extract Session ID (from JWT)
      │
Revoke Session
      │
Create Audit Log
      │
Return Success
```

---

## Logout All Devices

```
POST /auth/logout-all
```

Flow:

```
Delete All Sessions for User
      │
Return Success
```

---

# Email Verification

## Generate Code

* 6-digit numeric code
* Stored as SHA-256 hash
* Expires after 15 minutes

---

## Verify Email

```
POST /auth/verify-email
```

Flow:

```
Authenticated User
      │
Receive Code
      │
Hash Code
      │
Find Token
      │
Validate:
   - Not expired
   - Not used
      │
Mark User Verified
      │
Mark Token Used
      │
Create Audit Log
      │
Return Success
```

---

# JWT Authentication

JWT contains:

```
user_id
session_id
issued_at
expires_at
```

### Important Design Choice

JWT does NOT contain:

* role
* permissions
* email
* is_verified

This ensures:

* No stale authorization data
* Always validated against DB

---

# Authentication Middleware

Flow:

```
Read Authorization Header
      │
Validate Bearer Token
      │
Verify JWT Signature
      │
Extract Claims
      │
Load Session
      │
Validate Session
      │
Load User
      │
Store in Context
      │
Continue
```

---

## Context Values (IMPORTANT)

The middleware injects:

```
"user"    → full user object
"session" → session object
```

Used by:

* Handlers
* Rate Limiter
* Future RBAC system

---

# Authorization Middleware

```
RequireAuth()
```

→ ensures authenticated user exists

```
RequireVerified()
```

→ ensures email is verified

---

# Security Layer

Located in:

```
internal/application/security/
```

Handles:

* Password hashing
* Password validation
* Input filtering

---

# Token Service Layer

Located in:

```
internal/application/token/
```

Handles:

* Access token generation
* Refresh logic
* Token validation

---

# Security Summary

* Passwords hashed with bcrypt
* Refresh tokens stored as SHA-256 hashes
* JWTs short-lived (15 min)
* Sessions tracked per device
* Tokens are single-use (verification)
* Login history tracked
* Audit logs recorded

---

# Design Strengths

* Clean separation of concerns
* Device-aware authentication
* Secure token handling
* Audit-ready system
* Extensible for RBAC and permissions
* Ready for distributed scaling

---

# Future Improvements

* Role-Based Access Control (RBAC)
* Permission system
* Two-factor authentication (2FA)
* Session dashboard (active devices)
* Suspicious login detection

---

# Summary

The authentication module provides a **secure, scalable, and extensible foundation** for:

* Identity management
* Session control
* API security
* Event auditing

It is tightly integrated with middleware and designed to support future authorization features.
