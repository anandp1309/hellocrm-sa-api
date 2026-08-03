# Import and Export Building Block

## Import

- Validate file type and size.
- Store the original securely.
- Parse asynchronously for large files.
- Validate every row.
- Produce a clear row-level error report.
- Do not partially commit unless the feature specification allows it.
- Use idempotency or duplicate detection.
- Record importer, scope, file, status, counts, and timestamps.
- Allow safe retry when possible.

## Export

- Apply the same permissions and filters as the source screen.
- Run large exports asynchronously.
- Use expiring authorised download links.
- Avoid exporting sensitive fields by default.
- Audit sensitive exports.
- Define retention and cleanup.

## Statuses

Recommended technical statuses:

```text
queued
processing
completed
completed_with_errors
failed
cancelled
```
