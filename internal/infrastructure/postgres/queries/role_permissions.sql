-- name: AssignPermissionToRole :exec
INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;

-- name: RemovePermissionFromRole :exec
DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2;

-- name: ListPermissionsByRole :many
SELECT p.* FROM permissions p JOIN role_permissions rp ON rp.permission_id = p.id WHERE rp.role_id = $1 ORDER BY p.resource, p.action;

-- name: ListRolesByPermission :many
SELECT r.* FROM roles r JOIN role_permissions rp ON rp.role_id = r.id WHERE rp.permission_id = $1 ORDER BY r.name;

-- name: HasPermission :one
SELECT COUNT(*) > 0 FROM permissions p 
JOIN role_permissions rp ON p.id = rp.permission_id 
JOIN roles r ON r.id = rp.role_id 
WHERE r.name = $1 AND p.resource = $2 AND p.action = $3;

-- name: RoleHasPermission :one
SELECT EXISTS (SELECT 1 FROM role_permissions WHERE role_id = $1 AND permission_id = $2);