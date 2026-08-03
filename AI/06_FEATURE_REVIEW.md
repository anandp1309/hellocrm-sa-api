# Feature Review

A feature is not complete until this review has been performed.

## Functional completeness

- Approved requirement is fully implemented.
- No unapproved business behaviour was introduced.
- Status transitions and edge cases are handled.
- Failure behaviour is predictable.
- Out-of-scope work was not added.

## Backend

- Handlers contain transport logic only.
- Services contain business rules and transaction ownership.
- Repositories use SQLC.
- Validation is complete.
- Authentication and authorization are enforced.
- Tenant/project/object ownership is enforced.
- Sensitive actions are audited.
- Errors are mapped safely.
- Context cancellation and timeouts are respected.
- Duplicate-sensitive operations are idempotent where required.

## Database

- Migration is included.
- Constraints express real invariants.
- Indexes match actual query patterns.
- SQLC generation is current.
- No `SELECT *`.
- Money, timestamps, IDs, and nullability follow standards.
- Concurrency and locking behaviour is safe.
- Number series do not use `MAX + 1`.
- Rollback or roll-forward risk is understood.

## API

- Endpoint naming and HTTP methods follow the standard.
- Request and response models are typed.
- Error codes are stable.
- Pagination, filtering, and sorting are controlled.
- Breaking changes are documented.
- OpenAPI is updated where applicable.

## Frontend

- Shared components are reused.
- Loading, empty, success, validation, and error states exist.
- Duplicate submissions are prevented.
- Responsive behaviour is checked.
- Destructive actions require confirmation.
- Permissions are reflected in the UI without replacing backend enforcement.
- Accessibility basics are satisfied.

## Quality

- Unit and integration tests cover critical behaviour.
- Permission and tenant-isolation tests exist where applicable.
- Lint, tests, SQLC check, and production build were run.
- No debug code, fake data, unused code, hidden TODOs, or temporary bypasses remain.
- Documentation is updated where required.

## Required completion report

1. Files changed
2. Standards verified
3. Commands/checks executed
4. Migration or compatibility impact
5. Remaining risks or incomplete items
