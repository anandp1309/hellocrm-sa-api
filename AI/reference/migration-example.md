# Migration Reference

Up migration:

```sql
CREATE TABLE customers (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    customer_number VARCHAR(30) NOT NULL,
    full_name VARCHAR(150) NOT NULL,
    mobile VARCHAR(20),
    email VARCHAR(254),
    archived_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT customers_tenant_number_uk
        UNIQUE (tenant_id, customer_number),

    CONSTRAINT customers_mobile_length_ck
        CHECK (mobile IS NULL OR char_length(mobile) BETWEEN 8 AND 20)
);

CREATE INDEX customers_tenant_active_created_idx
    ON customers (tenant_id, created_at DESC, id DESC)
    WHERE archived_at IS NULL;
```

Down migration:

```sql
DROP TABLE customers;
```

Production notes:

- Do not casually drop populated columns.
- Large-table changes may require staged migrations.
- Do not edit migrations already applied in shared environments.
