# SQLC Reference

Query file:

```sql
-- name: GetCustomer :one
SELECT
    id,
    tenant_id,
    customer_number,
    full_name,
    mobile,
    email,
    created_at,
    updated_at
FROM customers
WHERE tenant_id = $1
  AND id = $2
  AND archived_at IS NULL;

-- name: ListCustomers :many
SELECT
    id,
    customer_number,
    full_name,
    mobile,
    email,
    created_at
FROM customers
WHERE tenant_id = $1
  AND archived_at IS NULL
  AND (
      sqlc.narg(search)::text IS NULL
      OR full_name ILIKE '%' || sqlc.narg(search)::text || '%'
      OR mobile ILIKE '%' || sqlc.narg(search)::text || '%'
  )
ORDER BY created_at DESC, id DESC
LIMIT $2 OFFSET $3;
```

Rules demonstrated:

- Explicit columns
- Tenant scope
- Archived filter
- Stable ordering
- Pagination
- No raw client-provided sort expression
