# Reporting Building Block

## Rules

- Define every metric and formula.
- Apply tenant/project and permission scope.
- Use server-side filtering.
- Use pagination for detail rows.
- Use background jobs for large exports.
- Avoid running reports directly against highly contended transactional paths when scale makes this unsafe.
- Introduce read replicas, materialised views, or analytical stores only after measured need and ADR approval.
- Include report period and generated timestamp.
- Ensure totals reconcile with detail.
- Audit sensitive exports.
