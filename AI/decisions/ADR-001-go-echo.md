# ADR-001: Use Go and Echo

Status: Accepted  
Date: 2026-07-12

## Context

EZ Solutions requires a backend stack suitable for long-lived SaaS products, predictable performance, small deployment artifacts, and a relatively small engineering team.

## Decision

Use Go as the backend language and Echo as the HTTP framework.

## Alternatives considered

- Node.js frameworks: strong ecosystem but higher runtime and dependency complexity for this team.
- Java/Spring: mature but heavier operational and development footprint.
- Python frameworks: productive, but less suitable as the company-wide default for high-throughput transactional APIs.

## Consequences

- Developers must follow Go conventions rather than framework-heavy patterns.
- Business logic remains independent of Echo where practical.
- The company maintains one primary backend stack.
