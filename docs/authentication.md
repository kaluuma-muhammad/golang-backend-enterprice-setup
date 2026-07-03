# Authentication Module Architecture

## Overview

The authentication module is responsible for user authentication, session management, JWT authentication, refresh token rotation, email verification, and account security.

The design follows a layered architecture inspired by Domain-Driven Design (DDD) and Clean Architecture.

```
HTTP Request
      │
      ▼
HTTP Handler
      │
      ▼
Application Service
      │
      ▼
Domain Interfaces
      │
      ▼
Infrastructure Repositories
      │
      ▼
PostgreSQL
```

The application layer contains business rules.

The infrastructure layer contains implementations such as PostgreSQL, JWT generation, email delivery, and logging.

The domain layer contains entities and repository contracts.

---

# Database Design

## Users

The `users` table stores permanent account information.

```
users
─────────────────────────────────────────────
id
email
password
first_name
last_name
is_verified
created_at
updated_at
```

### Purpose

The users table represents the identity of a person using the system.

Passwords are stored as bcrypt hashes.

The `is_verified` flag determines whether the user has verified their email address.

A user record never stores JWTs, refresh tokens, or verification codes.

---

## Sessions

The `sessions` table stores login sessions.

```
sessions
─────────────────────────────────────────────
id
user_id
refresh_token
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

Every successful login creates a new session.

Each browser, mobile device, or computer has its own session.

Sessions allow:

* Multiple device login
* Logout from one device
* Logout from all devices
* Refresh token rotation
* Session revocation
* Audit information

Refresh tokens are never stored in plain text.

Only SHA-256 hashes are stored.

---

## Tokens

The `tokens` table stores temporary single-use security tokens.

```
tokens
─────────────────────────────────────────────
id
user_id
type
token_hash
expires_at
used_at
created_at
```

### Supported Token Types

Current:

* email_verification

Future:

* password_reset
* email_change
* invitation
* magic_link
* two_factor_recovery

Unlike sessions, tokens are temporary and disposable.

---

# Authentication Workflow

## Registration

```
Client
   │
POST /auth/register
   │
   ▼
Validate Request
   │
   ▼
Check Email Exists
   │
   ▼
Hash Password
   │
   ▼
Create User
   │
   ▼
Generate Refresh Token
   │
   ▼
Hash Refresh Token
   │
   ▼
Create Session
   │
   ▼
Generate Access Token
   │
   ▼
Generate Verification Code
   │
   ▼
Hash Verification Code
   │
   ▼
Delete Previous Verification Token
   │
   ▼
Store Verification Token
   │
   ▼
Send Verification Email
   │
   ▼
Return
```

Response:

```
{
    "user": { ... },
    "access_token": "...",
    "refresh_token": "...",
    "expires_in": 900
}
```

The user is authenticated immediately after registration but is not yet verified.

---

# Login

```
POST /auth/login
```

Workflow

```
Find User
      │
Verify Password
      │
Generate Refresh Token
      │
Hash Refresh Token
      │
Create Session
      │
Generate Access Token
      │
Return Tokens
```

Each login creates a completely new session.

Existing sessions remain active.

---

# Refresh Token

Refresh tokens are rotated.

Workflow

```
Receive Refresh Token
      │
Hash Token
      │
Find Session
      │
Validate Session
      │
Generate New Refresh Token
      │
Hash New Token
      │
Update Session
      │
Generate New Access Token
      │
Return Both Tokens
```

Old refresh tokens immediately become invalid.

This protects against replay attacks.

---

# Logout

Logout only affects the current session.

```
Receive JWT
      │
Extract Session ID
      │
Delete Session
      │
Return Success
```

Deleting the session invalidates every refresh token belonging to that session.

The access token naturally expires after fifteen minutes.

---

# Logout All Devices

Future endpoint:

```
POST /auth/logout-all
```

Workflow

```
Delete All Sessions
      │
Return Success
```

Every refresh token belonging to the user immediately becomes invalid.

---

# Email Verification

After registration, a six-digit verification code is generated.

Example

```
483921
```

Only the SHA-256 hash is stored.

The plain code only exists in the email.

Workflow

```
Generate Code
      │
Hash Code
      │
Delete Previous Verification Token
      │
Store Hash
      │
Send Email
```

---

# Verify Email

```
POST /auth/verify-email
```

Workflow

```
Authenticated User
      │
Receive Verification Code
      │
Hash Code
      │
Find Verification Token
      │
Validate Owner
      │
Check Used
      │
Check Expiration
      │
Mark User Verified
      │
Mark Token Used
      │
Return Success
```

No new JWT is issued.

The user's verification status is loaded from the database on every authenticated request.

---

# Resend Verification

```
POST /auth/resend-verification
```

Workflow

```
Delete Existing Verification Token
      │
Generate New Code
      │
Store Hash
      │
Send Email
```

Only one active verification code exists at any time.

---

# JWT Authentication

JWT access tokens contain only:

```
User ID
Session ID
Issued At
Expires At
```

They intentionally do not contain:

* email
* role
* permissions
* is_verified

This prevents stale authorization data inside tokens.

---

# Authentication Middleware

Authentication middleware performs the following steps.

```
Read Authorization Header
      │
Validate Bearer Format
      │
Verify JWT
      │
Load User
      │
Load Session
      │
Store User In Context
      │
Store Session In Context
      │
Continue
```

Handlers never parse JWTs directly.

They retrieve the authenticated user from the request context.

---

# Authorization Middleware

Authentication and authorization are separated.

```
RequireAuth()
```

Ensures the request has a valid authenticated user.

```
RequireVerified()
```

Ensures the authenticated user has verified their email.

Future middleware may include:

* RequireRole()
* RequirePermission()
* RequireAdmin()

---

# Password Security

Passwords are hashed using bcrypt.

Plain passwords are never stored.

Verification uses bcrypt comparison.

---

# Refresh Token Security

Refresh tokens are generated using cryptographically secure random bytes.

Only SHA-256 hashes are stored in the database.

A database leak cannot reveal usable refresh tokens.

---

# Verification Code Security

Verification codes are six-digit numeric values.

Only their SHA-256 hashes are stored.

Verification codes:

* expire after fifteen minutes
* are single use
* belong to one user
* are invalidated when regenerated

---

# Session Security

Every login creates a unique session.

Each session has:

* its own refresh token
* expiration date
* device information
* IP address
* user agent

Compromising one refresh token does not affect other devices.

---

# Email Infrastructure

The email subsystem is infrastructure only.

```
Application
      │
      ▼
Email Service
      │
      ▼
Provider
      │
      ▼
SMTP
```

The email package knows how to send emails.

The authentication module decides which templates to send.

Current templates:

* Email Verification

Future templates:

* Password Reset
* Welcome Email
* Invitation
* Email Change Confirmation

---

# Future Enhancements

The architecture has been designed to support additional authentication features without major structural changes.

Planned additions include:

* Forgot Password
* Password Reset
* Email Change Verification
* Multi-Factor Authentication (MFA)
* OAuth Providers (Google, GitHub, Microsoft)
* Role-Based Access Control (RBAC)
* Permission-Based Authorization
* API Keys
* Device Management
* Session History
* Login Notifications
* Security Audit Logs
* Rate Limiting
* Account Lockout
* Background Cleanup Workers for Expired Tokens and Sessions

Because the authentication module separates domain logic, application services, infrastructure, and interfaces, these features can be added incrementally while preserving a clean, maintainable architecture.
