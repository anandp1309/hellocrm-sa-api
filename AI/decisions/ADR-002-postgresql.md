# ADR-002: Use PostgreSQL

Status: Accepted  
Date: 2026-07-12

## Context

Products require transactional integrity, reporting, rich indexing, JSON support when needed, and long-term maintainability.

## Decision

Use PostgreSQL as the default relational database.

## Alternatives considered

- MariaDB/MySQL: capable, but PostgreSQL is preferred for advanced constraints, indexing, JSONB, and transactional features.
- NoSQL-first storage: rejected for core transactional business data.
- Multiple databases per product: rejected unless a measured requirement justifies it.

## Consequences

- Teams need PostgreSQL-specific operational knowledge.
- Schema design should use relational constraints.
- JSONB is available but not a substitute for proper schema design.
