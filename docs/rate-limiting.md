# Rate Limiting System

## Overview

The rate limiting system protects the API from abuse, brute-force attacks, denial-of-service attempts, and excessive traffic. It is implemented as middleware at the API gateway, ensuring requests are validated before reaching the application layer.

The implementation supports multiple storage backends:

- **Memory** – Ideal for local development and single-instance deployments.
- **Redis** – Distributed rate limiting for Docker and production environments.

The storage backend is configurable without requiring code changes.

---

# Architecture

```
                       Client
                          │
                          ▼
                 Gin HTTP Middleware
                          │
                          ▼
               RateLimitMiddleware
                          │
                          ▼
               RateLimiter Interface
                 ┌────────┴────────┐
                 │                 │
                 ▼                 ▼
          Memory Service      Redis Service
                 │                 │
                 ▼                 ▼
          Memory Store         Redis Cache
```

The middleware depends only on the `RateLimiter` interface, allowing different implementations to be swapped through dependency injection.

---

# Features

The rate limiting system provides:

- Token Bucket rate limiting
- Route-specific limits
- Group-based limits
- User/IP identification
- Pluggable storage backends
- Automatic cleanup (memory backend)
- Distributed support (Redis backend)
- Configurable without recompilation

---

# Token Bucket Algorithm

The memory implementation uses Go's `golang.org/x/time/rate` package.

Benefits include:

- Smooth request processing
- Burst support
- Constant-time checks
- Efficient memory usage

---

# Limiting Strategy

Rate limits are resolved in the following order.

```
Route Limit
      │
      ▼
Group Limit
      │
      ▼
Default Limit
```

The first matching rule is applied.

---

# Route-Based Limits

Routes that require additional protection can have dedicated limits.

Example:

```
POST /api/v1/auth/login
→ 5 requests/minute
```

```
POST /api/v1/auth/register
→ 3 requests/minute
```

These limits are typically used for authentication endpoints that are vulnerable to brute-force attacks.

---

# Group-Based Limits

Application routes are divided into logical groups.

| Group | Description |
|---------|------------|
| Public | Unauthenticated users |
| Protected | Authenticated users |
| Verified | Verified users |

Each group has its own configurable request limits.

---

# Default Limits

If no route-specific or group-specific configuration exists, the default limit is applied.

This guarantees every endpoint is protected.

---

# Identity Strategy

## Unauthenticated Users

Requests are tracked using the client's IP address.

```
METHOD:PATH:ip:CLIENT_IP
```

Example

```
POST:/api/v1/auth/login:ip:192.168.1.20
```

---

## Authenticated Users

Authenticated requests use the user identifier.

```
METHOD:PATH:user:USER_ID
```

Example

```
GET:/api/v1/orders:user:997c0ac0
```

Using the user ID prevents multiple authenticated users behind the same network from affecting one another.

---

# Components

## Config

Defines:

- Route limits
- Group limits
- Default limits
- Cleanup intervals
- Maximum idle time

---

## Key Generator

Responsible for:

- Generating consistent request keys
- Using User IDs for authenticated requests
- Falling back to IP addresses
- Supporting reverse proxy headers

Supported headers include:

- X-Forwarded-For
- X-Real-IP

---

## RateLimiter Interface

The middleware depends on a common interface instead of a concrete implementation.

This allows switching between memory and Redis without modifying the middleware.

---

## Memory Service

The in-memory implementation uses:

- Token Bucket algorithm
- `sync.RWMutex`
- In-memory map
- Background cleanup worker

Recommended for:

- Local development
- Single server deployments
- Testing

---

## Memory Store

Stores limiter instances keyed by request identity.

Characteristics:

- Thread-safe
- Constant-time lookup
- Automatic limiter creation
- Periodic cleanup

---

## Redis Service

The Redis implementation stores counters inside Redis instead of application memory.

Characteristics:

- Shared across multiple application instances
- Supports horizontal scaling
- No in-memory cleanup required
- Uses Redis key expiration

Recommended for:

- Docker deployments
- Kubernetes
- Production environments
- Load-balanced APIs

---

## Cache Abstraction

The Redis implementation depends on the application's cache abstraction rather than the Redis client directly.

```
RedisService
       │
       ▼
common.Cache
       │
       ▼
Redis Client
```

This keeps the application independent of a specific cache implementation.

---

## Cleanup Worker

Only the memory implementation requires cleanup.

The worker:

- Runs periodically
- Removes inactive limiters
- Prevents memory leaks

Redis does not require cleanup because key expiration is handled automatically.

---

# Request Flow

```
Incoming Request
        │
        ▼
Authentication
        │
        ▼
Determine Identity
(IP or User ID)
        │
        ▼
Generate Request Key
        │
        ▼
Resolve Applicable Limit
(Route → Group → Default)
        │
        ▼
RateLimiter Interface
        │
        ├───────────────┐
        ▼               ▼
Memory Service    Redis Service
        │               │
        ▼               ▼
Decision (Allow/Deny)
        │
        ▼
Continue or Return 429
```

---

# Response Headers

Successful requests include:

```
X-RateLimit-Limit
```

When the request exceeds the configured limit:

```
Retry-After
```

is returned.

Example:

```
HTTP/1.1 429 Too Many Requests

Retry-After: 12

{
    "error": "rate_limit_exceeded",
    "message": "Too many requests. Please try again later.",
    "retry_after": 12
}
```

---

# Configuration

The implementation supports multiple backends.

Example:

```env
RATE_LIMIT_STORE=memory
```

or

```env
RATE_LIMIT_STORE=redis
```

Redis configuration:

```env
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

For Docker:

```env
RATE_LIMIT_STORE=redis
REDIS_HOST=redis
```

No code changes are required when switching implementations.

---

# Local Development

Recommended configuration:

```
RATE_LIMIT_STORE=memory
```

Advantages:

- No Redis dependency
- Faster startup
- Easier debugging

---

# Production

Recommended configuration:

```
RATE_LIMIT_STORE=redis
```

Advantages:

- Distributed rate limiting
- Shared counters
- Horizontal scalability
- Consistent limits across multiple API instances

---

# Design Principles

The implementation follows several design principles.

- Dependency Injection
- Interface-based design
- Separation of concerns
- Single Responsibility Principle
- Configurable infrastructure
- Pluggable storage backend

---

# Future Enhancements

With Redis integrated, additional features can reuse the same infrastructure:

- Permission caching (RBAC)
- User profile caching
- JWT blacklist
- Refresh token storage
- Email verification cache
- Password reset cache
- API response caching
- Distributed locking
- Background jobs
- Pub/Sub messaging
- Metrics and dashboards

---

# Summary

The rate limiting system provides a flexible, enterprise-ready solution for protecting the API.

Key characteristics include:

- Multi-level rate limiting
- Token Bucket algorithm
- Memory and Redis implementations
- Interface-driven architecture
- Dependency injection
- Horizontal scalability
- Environment-based configuration
- Automatic cleanup for memory storage
- Distributed support through Redis

This design allows the application to start with a lightweight in-memory implementation during development while seamlessly switching to Redis for distributed production deployments without requiring changes to the application code.