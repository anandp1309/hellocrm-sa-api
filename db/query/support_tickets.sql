-- name: CreateSupportTicket :one
INSERT INTO support_tickets (
    tenant_id, user_id, subject, description, priority
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetSupportTicket :one
SELECT * FROM support_tickets
WHERE id = $1 LIMIT 1;

-- name: ListOpenSupportTickets :many
SELECT * FROM support_tickets
WHERE status = 'Open' OR status = 'In Progress'
ORDER BY created_at DESC;

-- name: ListClosedSupportTickets :many
SELECT * FROM support_tickets
WHERE status = 'Closed' OR status = 'Resolved'
ORDER BY created_at DESC;

-- name: UpdateSupportTicketStatus :one
UPDATE support_tickets
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ReopenSupportTicket :one
UPDATE support_tickets
SET status = 'Open', 
    reopen_count = reopen_count + 1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: RateSupportTicket :one
UPDATE support_tickets
SET customer_satisfaction = $2,
    rating = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
