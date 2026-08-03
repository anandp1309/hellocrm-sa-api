-- name: ListInvoicesPaginated :many
SELECT 
    p.tenant_subscription_payment_uuid AS invoice_uuid,
    p.payment_number AS invoice_no,
    t.tenant_name AS customer,
    t.email_address,
    COALESCE(u_type.value_name, 'Subscription') AS billing_type,
    p.created_at AS invoice_date,
    CAST(p.created_at + INTERVAL '10 days' AS timestamptz) AS due_date,
    p.amount,
    COALESCE(u_status.value_name, 'Unknown') AS status,
    p.payment_date,
    COUNT(*) OVER() AS total_records
FROM tenant_subscription_payment p
JOIN tenant t ON p.tenant_uuid = t.tenant_uuid
LEFT JOIN tenant_subscription ts ON p.tenant_subscription_uuid = ts.tenant_subscription_uuid
LEFT JOIN mst_plan_price pp ON ts.plan_price_uuid = pp.plan_price_uuid
LEFT JOIN mst_plan pl ON pp.plan_uuid = pl.plan_uuid
LEFT JOIN mst_universal u_status ON p.payment_status_universal_uuid = u_status.universal_uuid
LEFT JOIN mst_universal u_type ON pl.plan_uuid = u_type.universal_uuid -- Mocking billing type
WHERE 
    p.is_deleted = false
    AND ($1::text = '' OR p.payment_number ILIKE '%' || $1 || '%' OR t.tenant_name ILIKE '%' || $1 || '%')
    AND ($2::text = '' OR u_status.value_name = $2)
    AND (
        ($5::text = '' OR $6::text = '') 
        OR 
        (p.created_at >= $5::date AND p.created_at <= $6::date)
    )
ORDER BY p.created_at DESC
LIMIT $4 OFFSET $3;

-- name: GetInvoiceStats :one
SELECT 
    COUNT(*) AS total_invoices,
    COALESCE(SUM(amount) FILTER (WHERE u_status.value_name = 'Paid'), 0) AS paid_amount,
    COUNT(*) FILTER (WHERE u_status.value_name = 'Paid') AS paid_count,
    COALESCE(SUM(amount) FILTER (WHERE u_status.value_name = 'Pending'), 0) AS pending_amount,
    COUNT(*) FILTER (WHERE u_status.value_name = 'Pending') AS pending_count,
    COALESCE(SUM(amount) FILTER (WHERE u_status.value_name = 'Failed' OR u_status.value_name = 'Overdue'), 0) AS overdue_amount,
    COUNT(*) FILTER (WHERE u_status.value_name = 'Failed' OR u_status.value_name = 'Overdue') AS overdue_count
FROM tenant_subscription_payment p
LEFT JOIN mst_universal u_status ON p.payment_status_universal_uuid = u_status.universal_uuid
WHERE p.is_deleted = false;

-- name: GetInvoice :one
SELECT 
    p.tenant_subscription_payment_uuid AS invoice_uuid,
    p.payment_number AS invoice_no,
    t.tenant_name AS customer,
    t.email_address,
    COALESCE(u_type.value_name, 'Subscription') AS billing_type,
    p.created_at AS invoice_date,
    CAST(p.created_at + INTERVAL '10 days' AS timestamptz) AS due_date,
    p.amount,
    COALESCE(u_status.value_name, 'Unknown') AS status,
    p.payment_date
FROM tenant_subscription_payment p
JOIN tenant t ON p.tenant_uuid = t.tenant_uuid
LEFT JOIN tenant_subscription ts ON p.tenant_subscription_uuid = ts.tenant_subscription_uuid
LEFT JOIN mst_plan_price pp ON ts.plan_price_uuid = pp.plan_price_uuid
LEFT JOIN mst_plan pl ON pp.plan_uuid = pl.plan_uuid
LEFT JOIN mst_universal u_status ON p.payment_status_universal_uuid = u_status.universal_uuid
LEFT JOIN mst_universal u_type ON pl.plan_uuid = u_type.universal_uuid
WHERE p.tenant_subscription_payment_uuid = $1 AND p.is_deleted = false;
