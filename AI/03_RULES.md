# Non-Negotiable Engineering Rules

1. Inspect existing code before editing.
2. Preserve the approved repository architecture.
3. Do not introduce an alternative architecture, framework, or database access layer.
4. Use SQLC for application database access.
5. Do not use an ORM.
6. Use versioned migrations for every schema change.
7. Use UUIDv7 for new primary identifiers.
8. Store timestamps in UTC using `TIMESTAMPTZ`.
9. Use `NUMERIC(18,2)` for normal financial values.
10. Do not use PostgreSQL ENUM types.
11. Do not put business logic in HTTP handlers.
12. Do not put business logic in generic frontend components.
13. Validate all untrusted input on the server.
14. Enforce authentication and authorization on the backend.
15. Hidden or disabled frontend controls are not security controls.
16. Enforce tenant, project, workspace, and object ownership in repository queries and services.
17. Add audit logging for sensitive and business-critical actions.
18. Use explicit transactions for multi-step business operations.
19. Never use `SELECT *` in production queries.
20. Paginate potentially large lists.
21. Do not expose SQL errors, stack traces, or infrastructure details to clients.
22. Do not hardcode secrets, environment URLs, IDs, or business configuration.
23. Reuse existing components, helpers, and utilities before creating new ones.
24. Avoid unnecessary abstractions, interfaces, generic repositories, and premature generalisation.
25. Prefer a modular monolith; microservices require an ADR and explicit approval.
26. Handle loading, empty, success, validation, and error states in UI features.
27. Add tests for critical business rules, permissions, tenant isolation, and transaction behaviour.
28. Generated code must pass formatting, lint, tests, and production build.
29. Never claim checks passed unless they were executed.
30. Do not create fake implementations, silent fallbacks, placeholders, or hidden TODOs.
31. Never change product behaviour without explicit approval.
32. Never generate sequential business numbers using `MAX(number) + 1`.
33. Financial and audit records must not be physically deleted without an approved retention rule.
34. Errors should be logged once at the application boundary, not repeatedly at every layer.
35. Security-sensitive changes require explicit review.
