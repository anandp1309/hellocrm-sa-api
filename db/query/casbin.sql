-- name: GetAllRolePermissions :many
SELECT 
    rp.role_uuid,
    p.permission_code
FROM role_permission rp
JOIN role r ON rp.role_uuid = r.role_uuid
JOIN mst_permission p ON rp.permission_uuid = p.permission_uuid
WHERE rp.is_granted = TRUE AND rp.is_deleted = FALSE AND r.is_deleted = FALSE AND p.deleted_at IS NULL;
