# API Contract Template

## Endpoint

```text
METHOD /api/v1/resource
```

## Purpose

## Authentication and permission

## Path parameters

## Query parameters

## Request body

```json
{}
```

## Success response

Status:

```json
{
  "data": {},
  "meta": null
}
```

## Validation errors

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed.",
    "fields": {},
    "request_id": "req_01..."
  }
}
```

## Other error codes

| Status | Code | Condition |
|---:|---|---|

## Idempotency

## Tenant/project scope

## Audit event

## Compatibility notes
