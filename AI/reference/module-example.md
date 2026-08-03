# Reference Module Flow

A normal full-stack feature is implemented in this order:

1. Confirm feature specification.
2. Identify permissions and tenant/project scope.
3. Design schema, constraints, and indexes.
4. Create migration.
5. Write SQLC queries.
6. Generate SQLC code.
7. Implement repository.
8. Implement service and transaction logic.
9. Implement handler and route.
10. Add standard API mapping.
11. Add audit and background-job integration where required.
12. Add backend tests.
13. Implement typed frontend API access.
14. Build UI using shared components.
15. Add loading, empty, validation, success, and error states.
16. Run `AI/06_FEATURE_REVIEW.md`.
17. Submit PR using the template.

## Example service signature

```go
func (s *CustomerService) Create(
    ctx context.Context,
    actor auth.Actor,
    input CreateCustomerInput,
) (Customer, error)
```

## Example transaction rule

If creation allocates a number series and writes an audit entry, all three steps must commit atomically.
