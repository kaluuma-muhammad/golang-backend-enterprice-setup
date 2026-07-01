-- name: CreateSession :exec

INSERT INTO sessions (
    id,
    user_id,
    refresh_token,
    user_agent,
    ip_address,
    device_name,
    expires_at,
    last_used_at,
    revoked_at,
    revoked_reason,
    created_at,
    updated_at
)
VALUES (
    $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12
);

-- name: GetSessionByID :one

SELECT * FROM sessions WHERE id = $1;

-- name: GetSessionByToken :one

SELECT * FROM sessions WHERE refresh_token = $1;

-- name: UpdateSessionRefreshToken :exec

UPDATE sessions SET refresh_token = $2, expires_at = $3, last_used_at = NOW(), updated_at = NOW() WHERE id = $1;

-- name: RevokeRefreshToken :exec

UPDATE sessions SET revoked_at = NOW(), revoked_reason = $2, updated_at = NOW() WHERE id = $1;

-- name: RevokeAllSessions :exec

UPDATE sessions SET revoked_at = NOW(), revoked_reason = $2, updated_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL;

-- name: UpdateSessionLastUsed :exec

UPDATE sessions SET last_used_at = NOW(), updated_at = NOW() WHERE id = $1;

-- name: DeleteSessionByID :exec

DELETE FROM sessions WHERE id = $1;

-- name: DeleteSessionsByUserID :exec

DELETE FROM sessions WHERE user_id = $1;