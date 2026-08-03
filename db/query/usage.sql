-- name: GetUsageStats :one
SELECT 
    (SELECT count(*) FROM tenant_subscription)::int AS active_subscriptions,
    (SELECT count(*) FROM "user")::int AS total_users,
    (SELECT COALESCE(sum(file_size), 0) FROM media)::bigint AS total_storage_bytes,
    (SELECT count(*) FROM message_log)::int AS total_messages_used
;
