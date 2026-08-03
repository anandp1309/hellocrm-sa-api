-- name: ListRoles :many
SELECT r.role_uuid, r.role_name, r.remarks, (SELECT COUNT(DISTINCT uw.user_tenant_uuid) FROM user_workspace uw WHERE uw.role_uuid = r.role_uuid AND uw.is_deleted = FALSE) AS user_count FROM role r WHERE r.is_deleted = FALSE ORDER BY r.sort_order ASC, r.created_at DESC;

-- name: GetRolePermissions :many
SELECT p.permission_code, rp.role_permission_uuid FROM role_permission rp JOIN mst_permission p ON rp.permission_uuid = p.permission_uuid WHERE rp.role_uuid = $1 AND rp.is_deleted = FALSE AND p.deleted_at IS NULL;

-- name: UpdateRole :exec
UPDATE role SET role_name = $2, remarks = $3, updated_at = NOW() WHERE role_uuid = $1 AND is_deleted = FALSE;

