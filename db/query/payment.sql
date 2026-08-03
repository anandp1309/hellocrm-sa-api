-- name: ListPaymentsPaginated :many
SELECT 
    p.tenant_subscription_payment_uuid,
    p.payment_number AS transaction_id,
    t.tenant_name AS customer,
    t.email_address,
    COALESCE(u_status.value_name, 'Unknown') AS status,
    COALESCE(pl.plan_name, 'None') AS plan,
    p.amount,
    COALESCE(u_mode.value_name, 'Unknown') AS method,
    p.payment_date,
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
    AND ($2::text = '' OR u_status.value_name = $2)
    AND ($3::text = '' OR u_mode.value_name = $3)
    AND ($4::text = '' OR pl.plan_name = $4)
    AND (
        ($7::text = '' OR $8::text = '') 
        OR 
        (p.payment_date >= $7::date AND p.payment_date <= $8::date)
    )
ORDER BY p.payment_date DESC, p.created_at DESC
LIMIT $6 OFFSET $5;

-- name: GetPaymentStats :one
SELECT 
    COUNT(*) AS total_payments,
    COALESCE(SUM(amount) FILTER (WHERE u_status.value_name = 'Paid'), 0) AS total_received,
    COALESCE(SUM(amount) FILTER (WHERE u_status.value_name = 'Pending'), 0) AS pending_amount,
    COUNT(*) FILTER (WHERE u_status.value_name = 'Pending') AS pending_count,
    COALESCE(SUM(amount) FILTER (WHERE u_status.value_name = 'Failed'), 0) AS failed_amount,
    COUNT(*) FILTER (WHERE u_status.value_name = 'Failed') AS failed_count,
    COALESCE(SUM(amount) FILTER (WHERE u_status.value_name = 'Refunded'), 0) AS refunded_amount,
    COUNT(*) FILTER (WHERE u_status.value_name = 'Refunded') AS refunded_count
FROM tenant_subscription_payment p
LEFT JOIN mst_universal u_status ON p.payment_status_universal_uuid = u_status.universal_uuid
WHERE p.is_deleted = false;

-- name: CreatePayment :one
INSERT INTO tenant_subscription_payment (
    tenant_subscription_payment_uuid,
    payment_number,
    tenant_uuid,
    tenant_subscription_uuid,
    payment_status_universal_uuid,
    payment_mode_universal_uuid,
    payment_date,
    amount,
    remarks,
    created_at,
    created_by_user_uuid
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetPayment :one
SELECT * FROM tenant_subscription_payment
WHERE tenant_subscription_payment_uuid = $1 AND is_deleted = false;

-- name: UpdatePayment :one
UPDATE tenant_subscription_payment
SET
    payment_status_universal_uuid = COALESCE($2, payment_status_universal_uuid),
    payment_mode_universal_uuid = COALESCE($3, payment_mode_universal_uuid),
    payment_date = COALESCE($4, payment_date),
    amount = COALESCE($5, amount),
    remarks = COALESCE($6, remarks),
    updated_at = $7,
    updated_by_user_uuid = $8
WHERE tenant_subscription_payment_uuid = $1 AND is_deleted = false
RETURNING *;

-- name: DeletePayment :exec
UPDATE tenant_subscription_payment
SET
    is_deleted = true,
    deleted_at = $2,
    deleted_by_user_uuid = $3
WHERE tenant_subscription_payment_uuid = $1;
