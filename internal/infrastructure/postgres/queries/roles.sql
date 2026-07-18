-- name: CreateRole :one
INSERT INTO roles (id, name, description, is_system, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING *;

-- name: GetRoleByID :one
SELECT * FROM roles WHERE id = $1;

-- name: GetRoleByName :one
SELECT * FROM roles WHERE name = $1;

-- name: GetPaginatedRoles :many
SELECT * FROM roles LIMIT $1 OFFSET $2;

-- name: CountRoles :one
SELECT COUNT(*) FROM roles;

-- name: ListRoles :many
SELECT * FROM roles ORDER BY name;

-- name: UpdateRole :one
UPDATE roles SET name = $2, description = $3, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: DeleteRole :exec
DELETE FROM roles WHERE id = $1;