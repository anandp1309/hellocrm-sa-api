-- name: ListAdminsPaginated :many
SELECT 
    u.user_uuid,
    u.first_name,
    u.last_name,
    ua.login_id AS username,
    u.email_address,
    r.role_name AS role,
    COALESCE(u_status.value_name, 'Unknown') AS status,
    CASE WHEN ua.is_mobile_verified THEN 'Enabled' ELSE 'Disabled' END AS mfa,
    ua.last_login_at,
    COUNT(*) OVER() AS total_records
FROM "user" u
JOIN user_auth ua ON u.user_uuid = ua.user_uuid
JOIN user_tenant ut ON u.user_uuid = ut.user_uuid
JOIN tenant t ON ut.tenant_uuid = t.tenant_uuid
LEFT JOIN user_workspace uw ON ut.user_tenant_uuid = uw.user_tenant_uuid
LEFT JOIN role r ON uw.role_uuid = r.role_uuid
LEFT JOIN mst_universal u_status ON u.user_status_universal_uuid = u_status.universal_uuid
WHERE 
    u.is_deleted = false
    AND t.tenant_code = 'SYS-001'
    AND ($1::text = '' OR u.first_name ILIKE '%' || $1 || '%' OR u.last_name ILIKE '%' || $1 || '%' OR u.email_address ILIKE '%' || $1 || '%' OR ua.login_id ILIKE '%' || $1 || '%')
    AND ($2::text = '' OR r.role_name = $2)
    AND ($3::text = '' OR u_status.value_name = $3)
    AND ($4::text = '' OR (CASE WHEN ua.is_mobile_verified THEN 'Enabled' ELSE 'Disabled' END) = $4)
ORDER BY u.created_at DESC
LIMIT $6 OFFSET $5;

-- name: GetAdminStats :one
SELECT 
    COUNT(*) AS total_users,
    COUNT(*) FILTER (WHERE u_status.value_name = 'Active') AS active_users,
    COUNT(*) FILTER (WHERE u_status.value_name = 'Inactive') AS inactive_users,
    COUNT(DISTINCT r.role_name) AS total_roles
FROM "user" u
JOIN user_tenant ut ON u.user_uuid = ut.user_uuid
JOIN tenant t ON ut.tenant_uuid = t.tenant_uuid
LEFT JOIN user_workspace uw ON ut.user_tenant_uuid = uw.user_tenant_uuid
LEFT JOIN role r ON uw.role_uuid = r.role_uuid
LEFT JOIN mst_universal u_status ON u.user_status_universal_uuid = u_status.universal_uuid
WHERE u.is_deleted = false AND t.tenant_code = 'SYS-001';

-- name: GetAdmin :one
SELECT 
    u.user_uuid,
    u.first_name,
    u.last_name,
    ua.login_id AS username,
    u.email_address,
    r.role_name AS role,
    r.role_uuid,
    u_status.value_name AS status,
    u.user_status_universal_uuid AS status_uuid,
    ua.is_mobile_verified AS mfa_enabled,
    ua.last_login_at
FROM "user" u
JOIN user_auth ua ON u.user_uuid = ua.user_uuid
JOIN user_tenant ut ON u.user_uuid = ut.user_uuid
JOIN tenant t ON ut.tenant_uuid = t.tenant_uuid
LEFT JOIN user_workspace uw ON ut.user_tenant_uuid = uw.user_tenant_uuid
LEFT JOIN role r ON uw.role_uuid = r.role_uuid
LEFT JOIN mst_universal u_status ON u.user_status_universal_uuid = u_status.universal_uuid
WHERE u.user_uuid = $1 AND u.is_deleted = false AND t.tenant_code = 'SYS-001';

-- name: DeleteAdmin :exec
UPDATE "user"
SET 
    is_deleted = true,
    updated_at = NOW(),
    updated_by_user_uuid = $2
WHERE user_uuid = $1;

-- name: CreateAdminUser :exec
INSERT INTO "user" (
    user_uuid, first_name, last_name, email_address, user_status_universal_uuid,
    created_at, created_by_user_uuid
) VALUES (
    $1, $2, $3, $4, $5,
    NOW(), $6
);

-- name: CreateAdminAuth :exec
INSERT INTO user_auth (
    user_auth_uuid, tenant_uuid, user_uuid, login_id, password_hash,
    created_at
) VALUES (
    $1, $2, $3, $4, $5,
    NOW()
);

-- name: CreateAdminTenant :exec
INSERT INTO user_tenant (
    user_tenant_uuid, tenant_uuid, user_uuid, employee_code, is_default_tenant,
    created_at
) VALUES (
    $1, $2, $3, $4, true,
    NOW()
);

-- name: CreateAdminWorkspace :exec
INSERT INTO user_workspace (
    user_workspace_uuid, tenant_uuid, user_tenant_uuid, workspace_uuid, role_uuid,
    is_default_workspace, created_at
) VALUES (
    $1, $2, $3, $4, $5,
    true, NOW()
);

-- name: UpdateAdminUser :exec
UPDATE "user" SET
    first_name = $2,
    last_name = $3,
    email_address = $4,
    user_status_universal_uuid = $5,
    updated_at = NOW(),
    updated_by_user_uuid = $6
WHERE user_uuid = $1;

-- name: UpdateAdminAuth :exec
UPDATE user_auth SET
    login_id = $2,
    password_hash = COALESCE(NULLIF($3::text, ''), password_hash),
    is_mobile_verified = $4,
    updated_at = NOW()
WHERE user_uuid = $1;

-- name: UpdateAdminWorkspace :exec
UPDATE user_workspace SET
    role_uuid = $2
WHERE user_tenant_uuid = (
    SELECT user_tenant_uuid FROM user_tenant WHERE user_uuid = $1 LIMIT 1
);

-- name: GetTenantByCode :one
SELECT tenant_uuid FROM tenant WHERE tenant_code = $1 LIMIT 1;


