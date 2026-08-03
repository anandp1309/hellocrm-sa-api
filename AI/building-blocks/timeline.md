# Timeline and Activity Building Block

A user-facing timeline explains meaningful business activity.

## Rules

- Keep chronological ordering deterministic.
- Show actor, action, timestamp, and useful context.
- Avoid exposing technical log details.
- Do not use the timeline as the immutable audit store.
- Support pagination for long histories.
- Attachments or links require normal authorization.
- System-generated and user-generated events should be distinguishable.
- Corrections should append new events rather than silently rewriting history.
