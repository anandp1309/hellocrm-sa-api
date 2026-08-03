-- name: CreateMstPlan :exec
INSERT INTO mst_plan (
    plan_uuid, plan_name, remarks, max_users, default_storage_bytes, default_sms_credits, default_whatsapp_credits, default_email_credits, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, now(), now()
);

-- name: CreateMstPlanPrice :exec
INSERT INTO mst_plan_price (
    plan_price_uuid, plan_uuid, price_amount, created_at, updated_at
) VALUES (
    $1, $2, $3, now(), now()
);

-- name: GetPlanByID :one
SELECT p.plan_uuid as id, p.plan_name as name, COALESCE(mt.value_name, 'Unknown') as plan_type, p.remarks as description, pp.price_amount as price, COALESCE(mt_cycle.value_name, 'Monthly') as interval, p.max_users, p.default_storage_bytes, p.default_sms_credits, p.default_whatsapp_credits, p.default_email_credits, p.created_at, p.updated_at
FROM mst_plan p
LEFT JOIN mst_plan_price pp ON p.plan_uuid = pp.plan_uuid
LEFT JOIN mst_universal mt ON p.plan_type_universal_uuid = mt.universal_uuid
LEFT JOIN mst_universal mt_cycle ON pp.billing_cycle_universal_uuid = mt_cycle.universal_uuid
WHERE p.plan_uuid = $1
LIMIT 1;

-- name: ListPlans :many
SELECT p.plan_uuid as id, p.plan_name as name, COALESCE(mt.value_name, 'Unknown') as plan_type, p.remarks as description, pp.price_amount as price, COALESCE(mt_cycle.value_name, 'Monthly') as interval, p.max_users, p.default_storage_bytes, p.default_sms_credits, p.default_whatsapp_credits, p.default_email_credits, p.created_at, p.updated_at
FROM mst_plan p
LEFT JOIN mst_plan_price pp ON p.plan_uuid = pp.plan_uuid
LEFT JOIN mst_universal mt ON p.plan_type_universal_uuid = mt.universal_uuid
LEFT JOIN mst_universal mt_cycle ON pp.billing_cycle_universal_uuid = mt_cycle.universal_uuid
ORDER BY p.created_at DESC;

-- name: UpdateMstPlan :exec
UPDATE mst_plan
SET plan_name = $2, remarks = $3, max_users = $4, default_storage_bytes = $5, default_sms_credits = $6, default_whatsapp_credits = $7, default_email_credits = $8, updated_at = now()
WHERE plan_uuid = $1;

-- name: UpdateMstPlanPrice :exec
UPDATE mst_plan_price
SET price_amount = $2, updated_at = now()
WHERE plan_uuid = $1;

-- name: DeleteMstPlan :exec
DELETE FROM mst_plan
WHERE plan_uuid = $1;
-- name: GetFirstPlanPriceByPlan :one
SELECT plan_price_uuid FROM mst_plan_price WHERE plan_uuid = $1 LIMIT 1;
