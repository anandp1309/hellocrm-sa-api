# Tenant, Project, and Workspace Isolation

Status: Required when the product has scoped access

## Core rule

A record ID does not grant access. Every read and write must be scoped to the authenticated actor's permitted tenant, project, workspace, or organisation.

## Repository requirements

Repository methods must receive scope explicitly:

```go
type Scope struct {
    TenantID  uuid.UUID
    ProjectID *uuid.UUID
}
```

Queries must include scope:

```sql
WHERE tenant_id = $1
  AND id = $2
```

## Service requirements

Services must:

- Validate actor membership
- Validate action permission
- Pass explicit scope to repositories
- Prevent cross-scope reassignment
- Recheck scope inside transactions

## Database requirements

- Include scope in relevant unique constraints.
- Include scope in indexes for common queries.
- Do not rely only on frontend filtering.
- Avoid global IDs in cache keys without scope.
- Include scope in object-storage paths or metadata.
- Include scope in background jobs and audit entries.

## API behaviour

For inaccessible resources, return either `404` or `403` according to the repository's security policy. Do not reveal whether a resource exists in another tenant.

## Testing

Required tests:

- Tenant A cannot read Tenant B records.
- Tenant A cannot update/delete Tenant B records.
- Cross-tenant ID substitution fails.
- Background jobs cannot process another tenant accidentally.
- Unique values can coexist across tenants when intended.
