# ADR-003: Use SQLC and prohibit ORM

Status: Accepted  
Date: 2026-07-12

## Context

The team needs explicit, reviewable SQL and compile-time generated Go types without runtime ORM behaviour.

## Decision

Use SQLC for application database access. ORMs are prohibited.

A query builder may be approved for genuinely dynamic reporting queries, but it must not become a second persistence architecture.

## Alternatives considered

- GORM/Ent/Bun: convenient but introduce abstraction, hidden query behaviour, and a second modelling layer.
- Handwritten row scanning: explicit but repetitive and more error-prone.
- Stored procedures for all logic: rejected as the default because product logic belongs in services.

## Consequences

- Developers must understand SQL.
- Generated code must remain current.
- Complex dynamic reporting requires deliberate design.
