# ADR-004: Use UUIDv7 primary identifiers

Status: Accepted  
Date: 2026-07-12

## Context

Products need globally unique identifiers while preserving better index locality than random UUIDv4 values.

## Decision

Use UUIDv7 for new primary identifiers.

## Alternatives considered

- Auto-increment integers: simple but expose sequence and complicate distributed creation.
- UUIDv4: unique but less index-friendly.
- ULID: suitable, but UUIDv7 aligns with the UUID ecosystem and time ordering.

## Consequences

- Use one approved UUIDv7 helper.
- Business-readable numbers remain separate.
- Do not rely on UUID timestamp ordering as a business sequence.
