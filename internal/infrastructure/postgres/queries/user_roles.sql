-- name: AssignRoleToUser :exec
INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;

-- name: RemoveRoleFromUser :exec
DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2;

-- name: ListRolesByUser :many
SELECT r.* FROM roles r JOIN user_roles ur ON ur.role_id = r.id WHERE ur.user_id = $1 ORDER BY r.name;

-- name: ListUserPermissions :many
SELECT DISTINCT p.* FROM permissions p 
JOIN role_permissions rp ON rp.permission_id = p.id 
JOIN user_roles ur ON ur.role_id = rp.role_id 
WHERE ur.user_id = $1 ORDER BY p.resource, p.action;

-- name: CountUserRoles :one
SELECT COUNT(*) FROM user_roles WHERE user_id = $1;

-- name: CountUserPermissions :one
SELECT COUNT(*) FROM permissions p 
JOIN role_permissions rp ON rp.permission_id = p.id 
JOIN user_roles ur ON ur.role_id = rp.role_id 
WHERE ur.user_id = $1;

-- name: UserHasRole :one
SELECT EXISTS (SELECT 1 FROM user_roles WHERE user_id = $1 AND role_id = $2);

-- name: UserHasPermission :one
SELECT EXISTS (
    SELECT 1 FROM user_roles ur 
    JOIN role_permissions rp ON ur.role_id = rp.role_id 
    JOIN permissions p ON p.id = rp.permission_id 
    WHERE ur.user_id = $1 AND p.name = $2)
;

