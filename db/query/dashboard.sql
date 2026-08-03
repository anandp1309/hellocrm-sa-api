-- name: GetDashboardStats :one
SELECT 
    (SELECT count(*) FROM "user")::int AS total_users,
    (SELECT count(*) FROM tenant)::int AS total_tenants,
    (SELECT count(*) FROM tenant_subscription)::int AS total_subscriptions
;
