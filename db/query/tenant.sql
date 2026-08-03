-- name: ListTenantsPaginated :many
SELECT 
    t.tenant_uuid,
    t.tenant_name,
    t.email_address,
    COALESCE(u_status.value_name, 'Unknown') AS status,
    COALESCE(u_type.value_name, 'Unknown') AS plan_type,
    COALESCE(p.plan_name, 'None') AS plan,
    COALESCE(u_cycle.value_name, 'Monthly') AS billing_cycle,
    ts.subscription_start_date AS start_date,
    ts.subscription_end_date AS next_renewal,
    COALESCE(pp.price_amount, 0) AS mrr,
    COUNT(*) OVER() AS total_records
FROM tenant t
LEFT OUTER JOIN mst_universal u_status ON t.tenant_status_universal_uuid = u_status.universal_uuid
LEFT OUTER JOIN tenant_subscription ts ON t.tenant_uuid = ts.tenant_uuid AND ts.is_deleted = false
LEFT OUTER JOIN mst_plan_price pp ON ts.plan_price_uuid = pp.plan_price_uuid
LEFT OUTER JOIN mst_plan p ON pp.plan_uuid = p.plan_uuid
LEFT OUTER JOIN mst_universal u_type ON p.plan_type_universal_uuid = u_type.universal_uuid
LEFT OUTER JOIN mst_universal u_cycle ON pp.billing_cycle_universal_uuid = u_cycle.universal_uuid
WHERE 
    t.is_deleted = false
    AND ($1::text = '' OR t.tenant_name ILIKE '%' || $1 || '%' OR t.email_address ILIKE '%' || $1 || '%')
    AND ($2::text = '' OR u_status.value_name = $2)
    AND ($3::text = '' OR u_type.value_name = $3)
    AND ($4::text = '' OR p.plan_name = $4)
    AND ($5::text = '' OR u_cycle.value_name = $5)
ORDER BY t.created_at DESC
LIMIT $7 OFFSET $6;

-- name: GetTenantStats :one
SELECT 
    COUNT(*) AS total_customers,
    COUNT(*) FILTER (WHERE u_status.value_name = 'Active') AS active_customers,
    COUNT(*) FILTER (WHERE u_status.value_name = 'Trial') AS trial_customers,
    COUNT(*) FILTER (WHERE u_status.value_name IN ('Inactive', 'Suspended')) AS inactive_customers,
    COALESCE(SUM(pp.price_amount) FILTER (WHERE u_status.value_name = 'Active'), 0) AS total_mrr
FROM tenant t
LEFT OUTER JOIN mst_universal u_status ON t.tenant_status_universal_uuid = u_status.universal_uuid
LEFT OUTER JOIN tenant_subscription ts ON t.tenant_uuid = ts.tenant_uuid AND ts.is_deleted = false
LEFT OUTER JOIN mst_plan_price pp ON ts.plan_price_uuid = pp.plan_price_uuid
WHERE t.is_deleted = false;

-- name: CreateTenant :exec
INSERT INTO tenant (
    tenant_uuid, tenant_code, tenant_id, tenant_name, email_address,
    contact_person_name, mobile_number, country_name, state_name, city_name, address, gst_number, remarks,
    created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, now(), now()
);

-- name: UpdateTenant :exec
UPDATE tenant
SET tenant_name = $2, email_address = $3, updated_at = now()
WHERE tenant_uuid = $1 AND is_deleted = false;

-- name: DeleteTenant :exec
UPDATE tenant
SET tenant_status_universal_uuid = (
    SELECT universal_uuid FROM mst_universal 
    WHERE value_name = 'Inactive' 
    LIMIT 1
), updated_at = now()
WHERE tenant.tenant_uuid = $1 AND tenant.is_deleted = false;

-- name: GetTenantByID :one
SELECT 
    tenant_uuid, tenant_code, tenant_id, tenant_name, email_address,
    contact_person_name, mobile_number, country_name, state_name, city_name, address, gst_number, remarks,
    created_at, updated_at
FROM tenant
WHERE tenant_uuid = $1 AND is_deleted = false;
