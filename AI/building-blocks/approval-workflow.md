# Approval Workflow Building Block

Use for maker-checker, request-review, booking approval, payment confirmation, or other controlled transitions.

## Required model

Define explicit states. Example:

```text
draft -> submitted -> approved
                 \-> rejected
                 \-> correction_requested -> submitted
```

The actual states depend on the product requirement.

## Rules

- State transitions must be validated in the service.
- The approver must have explicit permission.
- Maker and checker separation must be enforced where required.
- Approval must run in a transaction.
- Update state, business side effects, and audit entry atomically.
- Use row locking or optimistic concurrency to prevent double approval.
- Repeated approval requests must be idempotent or return a stable conflict.
- Rejection/correction reason may be required by the feature specification.
- Notifications should occur after commit, normally through a job/outbox mechanism.
- Never allow the frontend to write arbitrary state values.

## Audit

Record:

- Previous state
- New state
- Actor
- Timestamp
- Reason/remarks
- Relevant entity
