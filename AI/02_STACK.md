# Approved Technology Stack

Status: Active  
Version: 1.0.0  
Last reviewed: 2026-07-12  
Next review: 2027-01-15

## Backend

- Language: Go
- HTTP framework: Echo
- Application style: Modular monolith by default
- Database access: SQLC
- ORM: Prohibited
- API style: REST
- API specification: OpenAPI
- Validation: Explicit request validation at the application boundary
- Logging: Structured logging
- Testing: Go test with unit and integration coverage based on risk

## Database

- PostgreSQL
- Primary identifiers: UUIDv7
- Money: `NUMERIC(18,2)` unless a documented domain requires different precision
- Timestamps: `TIMESTAMPTZ`, stored in UTC
- Schema changes: Versioned migrations only
- PostgreSQL ENUM types: Prohibited
- Core searchable business data: Normalised relational columns
- JSONB: Allowed only for genuinely flexible or external payload data

## Frontend

- SvelteKit
- TypeScript strict mode
- Server-side rendering and server actions where appropriate
- Shared reusable UI components
- Typed API client and typed request/response models
- Accessibility and responsive behaviour are required

## Storage and operations

- S3-compatible object storage
- Docker
- CI validation before merge
- Health and readiness checks
- Graceful shutdown
- Environment-based configuration
- Secrets outside source control and container images

## Version policy

Exact approved versions must be recorded in each repository. Use stable, supported releases. Major upgrades require:

- Compatibility review
- Migration plan
- Test plan
- Rollback plan
- ADR and approval
