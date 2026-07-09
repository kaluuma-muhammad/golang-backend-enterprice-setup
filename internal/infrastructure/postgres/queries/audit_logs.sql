-- name: CreateAuditLog :exec
INSERT INTO audit_logs (
    id,
    user_id,
    action,
    entity_type,
    entity_id,
    ip_address,
    user_agent,
    metadata,
    created_at
)
VALUES (
    $1,$2,$3,$4,$5,$6,$7,$8,$9
);

-- name: GetAuditLogsByUser :many
SELECT * FROM audit_logs WHERE user_id = $1 ORDER BY created_at DESC;

-- name: GetAuditLogs :many
SELECT * FROM audit_logs ORDER BY created_at DESC;

-- name: GetAuditLogByID :one
SELECT * FROM audit_logs WHERE id = $1;
