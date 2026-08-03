# CRUD / Master Data Building Block

Use this for controlled master records and simple business entities.

## Required capabilities

- Paginated list
- Search/filter where useful
- Create
- View
- Update
- Archive/restore when domain-valid
- Permission enforcement
- Validation
- Audit for important changes

## Rules

- Do not assume physical delete is valid.
- Use uniqueness constraints, not only pre-check queries.
- Handle duplicate conflicts gracefully.
- Keep list queries scoped and paginated.
- Use explicit allowed sort fields.
- Use PATCH semantics for partial updates only when the API contract supports it.
- Prevent updates to immutable fields.
- Include optimistic concurrency if lost updates are a realistic risk.

## Tests

- Create valid record
- Validation failure
- Duplicate handling
- Permission denial
- Cross-tenant access denial
- Archive and restore behaviour
