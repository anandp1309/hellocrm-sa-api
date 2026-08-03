# Feature Checkpoint

Before starting or continuing a feature:

1. Re-read `02_STACK.md` and `03_RULES.md`.
2. Inspect the related existing code.
3. Confirm the current dependency and folder structure.
4. Read the relevant architecture and building-block documents.
5. Reuse existing components, helpers, queries, and utilities.
6. Identify:
   - Validation requirements
   - Permission and object-access requirements
   - Tenant/project isolation
   - Audit events
   - Transaction boundaries
   - Notifications or background jobs
   - UI states
   - Test cases
7. Use migrations and SQLC for database changes.
8. Preserve API compatibility unless change is approved.
9. Identify material ambiguity before implementation.
10. Prepare a concise implementation plan.

Recommended prompt:

```text
Read `AI/05_CHECKPOINT.md`, inspect the existing implementation, identify all affected layers, and prepare the implementation plan for this feature:

[feature details]
```
