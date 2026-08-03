# ADR-005: Use the standard REST API contract

Status: Accepted  
Date: 2026-07-12

## Context

Multiple products and developers require predictable request, success, list, and error behaviour.

## Decision

Use versioned REST APIs with:

- Resource-oriented endpoints
- Standard success envelopes
- Standard pagination metadata
- Stable machine-readable error codes
- OpenAPI documentation
- Backend-enforced authorization

## Alternatives considered

- GraphQL as default: rejected because it adds complexity not required by current products.
- Ad-hoc JSON per module: rejected because it creates frontend and maintenance inconsistency.
- HTTP 200 for all outcomes: rejected because it discards standard protocol semantics.

## Consequences

- Existing modules must not invent response formats.
- Breaking changes require versioning or migration.
