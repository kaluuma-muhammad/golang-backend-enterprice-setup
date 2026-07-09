-- name: CreateUser :exec
INSERT INTO users (
    id,
    email,
    password,
    first_name,
    last_name,
    is_verified,
    failed_login_attempts,
    locked_until,
    last_login_at,
    created_at,
    updated_at
)
VALUES (
    $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11
);

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserSecurity :one
SELECT id, failed_login_attempts, locked_until, last_login_at FROM users WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users
SET email = $2,
    first_name = $3,
    last_name = $4,
    is_verified = $5,
    updated_at = $6
WHERE id = $1;

-- name: VerifyUser :exec
UPDATE users SET is_verified = TRUE, updated_at = NOW() WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users SET password = $2, updated_at = NOW() WHERE id = $1;

-- name: IncrementFailedLoginAttempts :exec
UPDATE users
SET
    failed_login_attempts = failed_login_attempts + 1,
    updated_at = NOW()
WHERE id = $1;

-- name: ResetFailedLoginAttempts :exec
UPDATE users
SET
    failed_login_attempts = 0,
    locked_until = NULL,
    updated_at = NOW()
WHERE id = $1;

-- name: LockUserAccount :exec
UPDATE users
SET
    locked_until = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateLastLogin :exec
UPDATE users
SET
    last_login_at = NOW(),
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;