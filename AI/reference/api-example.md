# API Response Reference

## Single resource

```json
{
  "data": {
    "id": "0190f2b8-7d7a-7a0d-bbe4-2b8dd9116f01",
    "customer_number": "CUS-000123",
    "full_name": "Asha Mehta"
  },
  "meta": null
}
```

## Paginated list

```json
{
  "data": [
    {
      "id": "0190f2b8-7d7a-7a0d-bbe4-2b8dd9116f01",
      "customer_number": "CUS-000123",
      "full_name": "Asha Mehta"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 1,
    "total_pages": 1
  }
}
```

## Validation error

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed.",
    "fields": {
      "mobile": "Enter a valid mobile number."
    },
    "request_id": "req_01J2..."
  }
}
```

## Conflict

```json
{
  "error": {
    "code": "INVALID_STATE",
    "message": "This request has already been approved.",
    "request_id": "req_01J2..."
  }
}
```
