-- name: CreateLoginHistory :exec
INSERT INTO login_histories (
    id,
    user_id,
    session_id,
    status,
    ip_address,
    user_agent,
    device_name,
    failure_reason,
    login_at,
    logout_at,
    created_at,
    updated_at
)
VALUES (
    $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12
);

-- name: GetLoginHistoryByUserPaginated :many
SELECT * FROM login_histories WHERE user_id = $1 ORDER BY login_at DESC LIMIT $2 OFFSET $3;

-- name: CountLoginHistoryByUser :one
SELECT COUNT(*) FROM login_histories WHERE user_id = $1;

-- name: GetLoginHistoryByUser :many
SELECT * FROM login_histories WHERE user_id = $1 ORDER BY login_at DESC;

-- name: MarkLoginHistoryLogout :exec
UPDATE login_histories SET logout_at = NOW() WHERE session_id = $1 AND logout_at IS NULL;

-- name: GetLatestSuccessfulLogin :one
SELECT * FROM login_histories WHERE user_id = $1 AND status = 'success' ORDER BY login_at DESC LIMIT 1;