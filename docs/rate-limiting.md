# Rate Limiting System

## Overview

The rate limiting system protects the API from abuse, brute-force attacks, and excessive traffic. It is implemented at the **API Gateway level** and applies limits based on:

* IP address (unauthenticated users)
* User ID (authenticated users)
* Route and request type

---

## Architecture

```
Client
   ↓
API Gateway
   ├── Authentication
   ├── Rate Limiting  ← (THIS)
   ├── Logging
   └── Forward to Services
```

All downstream services are protected because requests are filtered at the gateway.

---

## Key Features

### 1. Token Bucket Algorithm

The system uses Go's `rate.Limiter`:

* Smooth request handling
* Allows short bursts
* Enforces long-term limits

---

### 2. Multi-Level Limiting

#### Route-Based Limits (Highest Priority)

Example:

```
POST /auth/login → 5 requests/min/IP
POST /auth/register → 3 requests/min/IP
```

---

#### Group-Based Limits

| Group     | Description           |
| --------- | --------------------- |
| Public    | Unauthenticated users |
| Protected | Authenticated users   |
| Verified  | Verified users        |

---

#### Default Limit

Fallback when no specific rule exists.

---

## Identity Strategy

### Unauthenticated Users

```
Key = METHOD:PATH:ip:CLIENT_IP
```

Example:

```
POST:/auth/login:ip:192.168.1.1
```

---

### Authenticated Users

```
Key = METHOD:PATH:user:USER_ID
```

Example:

```
GET:/orders:user:abc-123
```

This avoids penalizing multiple users sharing the same IP.

---

## Components

### 1. Config

Defines limits:

* Requests per window
* Burst size
* Route-specific rules
* Group limits

---

### 2. Key Generator

Responsible for:

* Building unique keys per request
* Handling proxy headers (`X-Forwarded-For`)

---

### 3. Store Interface

Abstract storage layer:

* Current: In-memory store
* Future: Redis (distributed support)

---

### 4. Memory Store

* Thread-safe (`sync.RWMutex`)
* O(1) lookups
* Stores limiter instances per key

---

### 5. Limiter Wrapper

Wraps `rate.Limiter`:

* Handles token bucket logic
* Calculates retry delay
* Provides limit metadata

---

### 6. Service Layer

Core logic:

* Resolves correct limit (route → group → default)
* Creates or retrieves limiter
* Tracks last usage
* Decides allow/deny

---

### 7. Cleanup Worker

* Runs periodically
* Removes inactive limiters
* Prevents memory leaks

---

### 8. Middleware

Integrated into Gin:

* Extracts user identity
* Applies rate limiting
* Returns `429 Too Many Requests` when exceeded

---

## Request Flow

```
Incoming Request
      ↓
Identify User (IP or user_id)
      ↓
Generate Key
      ↓
Resolve Limit (route/group/default)
      ↓
Check Limiter
      ↓
Allowed? ── Yes → Continue
         └─ No  → 429 Response
```

---

## Response Example

### Allowed

```
200 OK
X-RateLimit-Limit: 100
```

---

### Blocked

```
429 Too Many Requests
Retry-After: 10

{
  "error": "rate_limit_exceeded",
  "message": "Too many requests. Please try again later.",
  "retry_after": 10
}
```

---

## Cleanup Strategy

* Each limiter tracks last usage
* Inactive limiters are removed after a defined time
* Prevents unbounded memory growth

---

## Design Benefits

* High performance (in-memory)
* Scalable (can switch to Redis)
* Flexible (per-route + per-user limits)
* Secure (prevents brute-force attacks)
* Clean architecture (separated components)

---

## Future Improvements

* Redis-based distributed rate limiting
* Per-role limits (Admin, Customer, Rider)
* API key-based limits
* Global system throttling
* Rate limit dashboards & metrics

---

## Summary

The rate limiting system provides:

* Protection against abuse
* Fine-grained traffic control
* Seamless integration with authentication
* Scalable design for future growth

It is a critical component for maintaining API stability and security.
