# Number Series Building Block

Use for booking numbers, receipts, demands, customers, work orders, invoices, or other business-readable identifiers.

## Core rules

- Primary keys remain UUIDv7.
- Business numbers are separate.
- Never use `MAX(number) + 1`.
- Number allocation must be safe under concurrency.
- Define scope: global, tenant, project, branch, financial year, or document type.
- Define reset rules.
- Define whether gaps are acceptable.
- Issued numbers must not normally be reused.

## Transactional counter pattern

A counter row may contain:

```text
scope_id
series_type
financial_year
prefix
next_value
padding
```

Allocation:

1. Begin transaction.
2. Lock the series row with `FOR UPDATE`.
3. Read and increment `next_value`.
4. Build the formatted number.
5. Insert the business record.
6. Commit.

If gap-free legal numbering is required, document the specific accounting/legal rule and design accordingly.
