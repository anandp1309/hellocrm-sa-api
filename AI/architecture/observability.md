# Logging, Audit, Metrics, and Health

Status: Required

## Structured application logs

Include:

- Timestamp
- Level
- Message
- Request/correlation ID
- Service/module
- Tenant/project context where safe
- Operation
- Error classification

Do not log secrets or full sensitive payloads.

## Logging ownership

- Lower layers wrap and return errors.
- The application boundary logs the final error once.
- Avoid duplicate logs for the same failure.

## Audit logs

Audit logs are required for:

- Approval/rejection
- Payment confirmation/reversal
- Role and permission changes
- Sensitive configuration changes
- Status changes with business impact
- Record archival/restoration
- Sensitive document access where required

Include actor, scope, action, entity, entity ID, timestamp, and relevant before/after details.

## Business activity timeline

User-facing activity timelines are separate from audit and technical logs. They may be curated for readability.

## Metrics

Track meaningful operational signals:

- Request rate and latency
- Error rate
- Database pool usage
- Background-job success/failure
- Queue depth
- External dependency latency
- Critical business workflow failure counts

## Health endpoints

- Liveness: process is running.
- Readiness: application can serve traffic.
- Do not expose secrets or detailed infrastructure information.
