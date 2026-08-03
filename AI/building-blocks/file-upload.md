# File Upload Building Block

## Storage

- Use S3-compatible object storage.
- Keep objects private by default.
- Use generated object keys.
- Store original filename as metadata.
- Never trust user-provided paths or filenames.

## Validation

- Maximum size
- Allowed extension
- Detected MIME/content type
- Expected content structure where applicable
- Malware-scanning integration point

## Security

- Authorise upload initiation and download.
- Use expiring presigned URLs where suitable.
- Prevent cross-tenant access.
- Do not expose storage credentials to the browser.
- Audit sensitive downloads where required.

## Lifecycle

- Track owner entity and scope.
- Clean abandoned temporary uploads.
- Define retention and archival.
- Handle replacement without orphaning files.
