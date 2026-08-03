# REST API Standard

Status: Required

## URL design

- Version externally consumed APIs, for example `/api/v1`.
- Use plural resource names.
- Keep URLs resource-oriented.
- Use nested resources only when ownership is meaningful.
- Avoid verbs in URLs unless modelling a real command that cannot be expressed safely as resource state.

Examples:

```text
GET    /api/v1/customers
POST   /api/v1/customers
GET    /api/v1/customers/{customer_id}
PATCH  /api/v1/customers/{customer_id}
POST   /api/v1/bookings/{booking_id}/approve
```

## HTTP methods

- `GET`: Read only
- `POST`: Create or execute a non-idempotent command
- `PUT`: Replace a resource
- `PATCH`: Partial update
- `DELETE`: Delete only when deletion is valid; otherwise use an explicit archive/cancel command

## Success contract

Single resource:

```json
{
  "data": {
    "id": "0190f2b8-...",
    "name": "Example"
  },
  "meta": null
}
```

List response:

```json
{
  "data": [],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

## Error contract

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed.",
    "fields": {
      "email": "Enter a valid email address."
    },
    "request_id": "req_01..."
  }
}
```

The `fields` property may be omitted when no field-level errors exist.

## Status codes

Use meaningful status codes:

- `200` successful read/update
- `201` successful creation
- `204` successful action without response body
- `400` malformed request
- `401` unauthenticated
- `403` authenticated but forbidden
- `404` resource not found in accessible scope
- `409` conflict, duplicate, invalid transition
- `422` semantic validation failure
- `429` rate limited
- `500` unexpected server error
- `503` dependency unavailable

## Pagination

Use page-based pagination for standard administrative lists. Use cursor pagination when:

- Dataset is very large
- Records change frequently
- Stable traversal matters

Always enforce maximum page size.

## Filtering and sorting

- Whitelist filter fields.
- Whitelist sort fields.
- Do not pass raw client sort expressions into SQL.
- Define default order explicitly.
- Include a stable tiebreaker such as ID.

## Idempotency

Use idempotency keys for duplicate-sensitive actions such as:

- Payment creation
- Booking confirmation
- External webhook processing
- Notification submission
- Import initiation

## Versioning and compatibility

- Preserve response fields within a major API version.
- New optional fields are generally backward compatible.
- Renaming/removing fields requires a version or migration strategy.
- Document deprecation.

## OpenAPI

- Keep request, response, error, and auth requirements documented.
- Generated documentation must reflect actual behaviour.
