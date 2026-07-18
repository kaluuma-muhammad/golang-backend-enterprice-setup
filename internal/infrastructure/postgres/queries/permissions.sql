-- name: CreatePermission :one
INSERT INTO permissions (id, resource, action, name, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING *;

-- name: GetPermissionByID :one
SELECT * FROM permissions WHERE id = $1;

-- name: GetPermissionByName :one
SELECT * FROM permissions WHERE name = $1;

-- name: GetPaginatedPermissions :many
SELECT * FROM permissions LIMIT $1 OFFSET $2;

-- name: CountPermissions :one
SELECT COUNT(*) FROM permissions;

-- name: ListPermissions :many
SELECT * FROM permissions ORDER BY resource, action;

-- name: UpdatePermission :one
UPDATE permissions SET resource = $2, action = $3, name = $4, description = $5, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: DeletePermission :exec
DELETE FROM permissions WHERE id = $1;