-- name: ListGatewayTransactionsPaginated :many
SELECT 
    p.tenant_subscription_payment_uuid,
    p.payment_number AS transaction_id,
    COALESCE(u_mode.value_name, 'Unknown') AS gateway,
    t.tenant_name AS customer,
    t.email_address,
    COALESCE(pl.plan_name, 'None') AS plan_name,
    p.amount,
    COALESCE(u_status.value_name, 'Unknown') AS payment_status,
    p.created_at,
    COUNT(*) OVER() AS total_records
FROM tenant_subscription_payment p
JOIN tenant t ON p.tenant_uuid = t.tenant_uuid
LEFT JOIN tenant_subscription ts ON p.tenant_subscription_uuid = ts.tenant_subscription_uuid
LEFT JOIN mst_plan_price pp ON ts.plan_price_uuid = pp.plan_price_uuid
LEFT JOIN mst_plan pl ON pp.plan_uuid = pl.plan_uuid
LEFT JOIN mst_universal u_status ON p.payment_status_universal_uuid = u_status.universal_uuid
LEFT JOIN mst_universal u_mode ON p.payment_mode_universal_uuid = u_mode.universal_uuid
WHERE 
    p.is_deleted = false
    AND ($1::text = '' OR p.payment_number ILIKE '%' || $1 || '%' OR t.tenant_name ILIKE '%' || $1 || '%')
    AND ($2::text = '' OR u_mode.value_name = $2)
    AND ($3::text = '' OR u_status.value_name = $3)
    AND (
        ($6::text = '' OR $7::text = '') 
        OR 
        (p.created_at >= $6::date AND p.created_at <= $7::date)
    )
ORDER BY p.created_at DESC
LIMIT $5 OFFSET $4;

-- name: GetGatewayStats :one
SELECT 
    COUNT(*) AS total_transactions,
    COALESCE(SUM(amount), 0) AS total_amount,
    COALESCE(SUM(amount) FILTER (WHERE u_status.value_name = 'Paid'), 0) AS success_amount,
    COUNT(*) FILTER (WHERE u_status.value_name = 'Paid') AS success_count,
    COALESCE(SUM(amount) FILTER (WHERE u_status.value_name = 'Failed'), 0) AS failed_amount,
    COUNT(*) FILTER (WHERE u_status.value_name = 'Failed') AS failed_count,
    COALESCE(SUM(amount) FILTER (WHERE u_status.value_name = 'Pending'), 0) AS pending_amount,
    COUNT(*) FILTER (WHERE u_status.value_name = 'Pending') AS pending_count
FROM tenant_subscription_payment p
LEFT JOIN mst_universal u_status ON p.payment_status_universal_uuid = u_status.universal_uuid
WHERE p.is_deleted = false;

-- name: GetGatewayTransaction :one
SELECT 
    p.tenant_subscription_payment_uuid,
    p.payment_number AS transaction_id,
    COALESCE(u_mode.value_name, 'Unknown') AS gateway,
    t.tenant_name AS customer,
    t.email_address,
    COALESCE(pl.plan_name, 'None') AS plan_name,
    p.amount,
    COALESCE(u_status.value_name, 'Unknown') AS payment_status,
    p.created_at
FROM tenant_subscription_payment p
JOIN tenant t ON p.tenant_uuid = t.tenant_uuid
LEFT JOIN tenant_subscription ts ON p.tenant_subscription_uuid = ts.tenant_subscription_uuid
LEFT JOIN mst_plan_price pp ON ts.plan_price_uuid = pp.plan_price_uuid
LEFT JOIN mst_plan pl ON pp.plan_uuid = pl.plan_uuid
LEFT JOIN mst_universal u_status ON p.payment_status_universal_uuid = u_status.universal_uuid
LEFT JOIN mst_universal u_mode ON p.payment_mode_universal_uuid = u_mode.universal_uuid
WHERE p.tenant_subscription_payment_uuid = $1 AND p.is_deleted = false;
