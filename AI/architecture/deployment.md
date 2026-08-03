# Docker, CI/CD, and Deployment Standard

Status: Required

## Docker

- Use multi-stage builds.
- Use minimal trusted runtime images.
- Run as a non-root user.
- Do not include source-control metadata, secrets, or local development files.
- Define health and readiness endpoints.
- Support graceful shutdown.
- Pin base image versions according to repository policy.

## Configuration

- Use environment-based configuration.
- Validate required configuration at startup.
- Keep environment differences in configuration, not code branches.
- Never commit production secrets.

## CI requirements

Before merge:

- Format check
- Lint
- Unit tests
- Integration tests where configured
- SQLC generation check
- Migration validation
- Frontend type check
- Production build
- Dependency/security scan where available

## Deployment

- Use immutable build artifacts.
- Run migrations through a controlled deployment step.
- Separate migration execution from uncontrolled application startup when required.
- Use rolling or staged deployment where infrastructure supports it.
- Verify health after deployment.

## Rollback

Every release with schema or compatibility risk must define:

- Application rollback
- Migration roll-forward or rollback
- Data recovery approach
- Feature flag or operational disablement where relevant

## Backup and restore

- Backups are not complete until restore is tested.
- Define backup frequency, retention, encryption, and access.
- Test restore procedures periodically.
