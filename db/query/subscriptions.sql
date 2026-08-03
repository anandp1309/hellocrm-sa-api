-- name: ListSubscriptions :many
SELECT 
    t.tenant_name AS customer_name,
    t.email_address AS customer_email,
    ts.subscription_number,
    p.plan_name,
    u_type.value_name AS plan_type,
    u_cycle.value_name AS billing_cycle,
    u_status.value_name AS status,
    ts.subscription_start_date AS start_date,
    ts.subscription_end_date AS next_billing_date,
    pp.price_amount AS mrr,
    u_payment.value_name AS payment_status,
    COUNT(*) OVER() AS total_records
FROM tenant_subscription ts
JOIN tenant t ON ts.tenant_uuid = t.tenant_uuid
JOIN mst_plan_price pp ON ts.plan_price_uuid = pp.plan_price_uuid
JOIN mst_plan p ON pp.plan_uuid = p.plan_uuid
LEFT JOIN mst_universal u_type ON p.plan_type_universal_uuid = u_type.universal_uuid
LEFT JOIN mst_universal u_cycle ON pp.billing_cycle_universal_uuid = u_cycle.universal_uuid
LEFT JOIN tenant_subscription_payment tsp ON ts.tenant_subscription_uuid = tsp.tenant_subscription_uuid
LEFT JOIN mst_universal u_status ON t.tenant_status_universal_uuid = u_status.universal_uuid
LEFT JOIN mst_universal u_payment ON tsp.payment_status_universal_uuid = u_payment.universal_uuid
WHERE 
    (@search_name::text = '' OR t.tenant_name ILIKE '%' || @search_name || '%' OR ts.subscription_number ILIKE '%' || @search_name || '%')
    AND (@status::text = '' OR u_status.value_name = @status)
    AND (@plan_name::text = '' OR p.plan_name = @plan_name)
    AND (@billing_cycle::text = '' OR u_cycle.value_name = @billing_cycle)
    AND (@payment_status::text = '' OR u_payment.value_name = @payment_status)
ORDER BY ts.created_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateSubscription :one
INSERT INTO tenant_subscription (
    tenant_subscription_uuid, 
    subscription_number, 
    tenant_uuid, 
    plan_price_uuid, 
    subscription_start_date, 
    subscription_end_date, 
    amount_paid,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, NOW()
) RETURNING tenant_subscription_uuid;

-- name: CancelSubscription :exec
UPDATE tenant_subscription
SET is_deleted = true, deleted_at = NOW()
WHERE subscription_number = $1;
