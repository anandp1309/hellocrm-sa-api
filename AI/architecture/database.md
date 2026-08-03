# PostgreSQL and SQLC Standard

Status: Required

## Naming

- Use `snake_case`.
- Use plural table names.
- Use descriptive constraint and index names.
- Avoid ambiguous columns such as `status` without a documented domain meaning.

## Identifiers

- New primary IDs use UUIDv7.
- Generate IDs through one approved company helper.
- Do not mix random UUID generation methods across modules.
- Business display numbers are separate from primary IDs.

## Required timestamps

Normal mutable business tables should generally include:

```sql
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
```

Add `archived_at`, `created_by`, and `updated_by` only where the domain requires them. Do not add columns mechanically without use.

## Money

- Use `NUMERIC(18,2)` for ordinary INR amounts.
- Never use floating-point types for money.
- Define rounding behaviour in the business rule.
- Store currency code when multiple currencies are possible.

## Nullability

- Use `NOT NULL` when absence is not valid.
- Do not use empty strings as substitutes for null.
- Nullable fields must have a real domain meaning.
- Defaults must not hide missing required input.

## Constraints

Use database constraints for durable invariants:

- Primary keys
- Foreign keys
- Unique constraints
- Check constraints
- Non-null constraints

Foreign keys are required by default. Deviations require a documented architecture decision.

## Tenant and project scope

For multi-tenant data:

- Include tenant/project scope in every relevant query.
- IDs alone never establish access.
- Unique constraints should include tenant/project scope where appropriate.
- Background jobs must retain tenant/project context.
- Cross-tenant joins are prohibited unless explicitly authorised.
- Repository methods should accept scope explicitly.

Example:

```sql
-- name: GetCustomer :one
SELECT id, tenant_id, customer_number, full_name, mobile, email
FROM customers
WHERE tenant_id = $1
  AND id = $2
  AND archived_at IS NULL;
```

## JSONB

Use JSONB for:

- External payload snapshots
- Flexible metadata
- Rarely queried extension data

Do not place core searchable, sortable, unique, financial, or permission-related fields only inside JSONB.

## ENUMs

PostgreSQL ENUM types are prohibited. Prefer:

- Lookup tables for managed values
- Text columns with check constraints for stable technical states

## Migrations

- Every schema change uses a numbered/versioned migration.
- Migrations must be deterministic.
- Do not edit a migration already applied in shared environments.
- Large table changes require an online-safe migration plan.
- Destructive changes require backup and rollback/roll-forward planning.
- Data migrations should be resumable where possible.

## SQLC

- All application SQL belongs in SQLC query files.
- Generated code must not be manually edited.
- Do not use `SELECT *`.
- Select only required columns.
- Name queries by use case.
- CI must fail when generated SQLC output is stale.

## Transactions

The service/application layer owns transaction boundaries.

Use transactions for:

- Multi-table writes
- Status transition plus audit event
- Payment confirmation plus ledger posting
- Number allocation plus record creation
- Inventory or booking state changes
- Any operation that must succeed or fail as one unit

Keep transactions short. Do not call slow external services while holding database locks.

## Concurrency

Use database mechanisms rather than optimistic assumptions:

- Unique constraints
- `SELECT ... FOR UPDATE` where justified
- Version columns for optimistic concurrency
- Advisory locks only with documented reasoning
- Idempotency keys for duplicate-sensitive operations

## Number series

Never use `MAX(number) + 1`.

Use one of:

- PostgreSQL sequence for globally increasing numbers
- Transactionally locked number-series row for tenant/project/year-specific formats
- Preallocated ranges only when scale justifies them

The business number and primary UUID must remain separate.

## Indexes

- Index based on real query patterns.
- Composite index order must match filters and sorting.
- Index foreign-key lookup columns where needed.
- Use partial indexes for active/pending subsets where justified.
- Review write overhead.
- Use `EXPLAIN (ANALYZE, BUFFERS)` for critical or slow queries.

## Deletion and archival

- Distinguish delete, archive, deactivate, cancel, reverse, and expire.
- Do not soft-delete every table by default.
- Financial and audit records normally use reversal or archival, not physical deletion.
- Cascading deletes require explicit review.
