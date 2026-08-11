-- name: ListRoles :many
SELECT 
    r.role_uuid, 
    r.role_name, 
    r.remarks, 
    r.created_by_user_uuid,
    c.first_name AS created_by_first_name,
    c.last_name AS created_by_last_name,
    r.updated_by_user_uuid,
    u.first_name AS updated_by_first_name,
    u.last_name AS updated_by_last_name,
    (SELECT COUNT(DISTINCT uw.user_tenant_uuid) FROM user_workspace uw WHERE uw.role_uuid = r.role_uuid AND uw.is_deleted = FALSE) AS user_count 
FROM role r 
LEFT JOIN "user" c ON r.created_by_user_uuid = c.user_uuid
LEFT JOIN "user" u ON r.updated_by_user_uuid = u.user_uuid
WHERE r.is_deleted = FALSE 
ORDER BY r.sort_order ASC, r.created_at DESC;

-- name: GetRolePermissions :many
SELECT 
    r.role_name,
    m.module_name,
    p.permission_name AS access_right,
    p.permission_code,
    rp.role_permission_uuid,
    r.role_uuid,
    m.module_uuid,
    p.permission_uuid,
    rp.tenant_uuid,
    rp.is_granted
FROM role_permission rp
JOIN role r ON rp.role_uuid = r.role_uuid
JOIN mst_permission p ON rp.permission_uuid = p.permission_uuid
JOIN mst_module m ON rp.module_uuid = m.module_uuid
WHERE rp.role_uuid = $1 AND p.deleted_at IS NULL
ORDER BY m.module_name ASC, p.permission_name ASC;

-- name: UpdateRole :exec
UPDATE role SET role_name = $2, remarks = $3, updated_at = NOW() WHERE role_uuid = $1 AND is_deleted = FALSE;

-- name: ListModules :many
SELECT module_uuid, module_name FROM mst_module WHERE is_deleted = FALSE ORDER BY module_name ASC;

-- name: CreateNewRole :one
INSERT INTO role (role_uuid, tenant_uuid, role_name, remarks, created_at, sort_order, is_system) 
VALUES ($1, $2, $3, $4, NOW(), 0, false) RETURNING role_uuid;

-- name: GetPermissionByCode :one
SELECT permission_uuid FROM mst_permission WHERE permission_code = $1;

-- name: CreateRolePermission :exec
INSERT INTO role_permission (role_permission_uuid, role_uuid, module_uuid, permission_uuid, is_granted, created_at)
VALUES ($1, $2, $3, $4, TRUE, NOW());

-- name: DeleteRolePermissions :exec
UPDATE role_permission SET is_granted = FALSE, updated_at = NOW() WHERE role_uuid = $1;

-- name: ListAllPermissions :many
SELECT p.permission_code, p.permission_uuid 
FROM mst_permission p
WHERE p.deleted_at IS NULL;

-- name: AddRolePermission :exec
WITH updated AS (
    UPDATE role_permission
    SET is_granted = TRUE, is_deleted = FALSE, updated_at = NOW()
    WHERE role_uuid = $2 AND module_uuid = $3 AND permission_uuid = $4
    RETURNING role_permission_uuid
)
INSERT INTO role_permission (role_permission_uuid, tenant_uuid, role_uuid, module_uuid, permission_uuid, is_granted, created_at)
SELECT $1, (SELECT tenant_uuid FROM role WHERE role_uuid = $2), $2, $3, $4, TRUE, NOW()
WHERE NOT EXISTS (SELECT 1 FROM updated);

-- name: DeleteRolePermission :exec
UPDATE role_permission SET is_granted = FALSE, updated_at = NOW()
WHERE role_uuid = $1 AND module_uuid = $2 AND permission_uuid = $3;

