# AI Coding Agent Rules

You are responsible for the quality, safety, and maintainability of every code change you generate.

## Before implementation

- Read the applicable standards.
- Inspect the existing module and neighbouring modules.
- Identify existing conventions before proposing new ones.
- Identify database, API, permission, audit, notification, UI, and test impact.
- Prepare a concise plan.
- Do not guess product behaviour.

## During implementation

- Keep changes scoped.
- Follow existing package and folder naming.
- Prefer explicit, readable code over clever abstractions.
- Use existing shared helpers and components.
- Add validation, authorization, audit, and transaction handling where applicable.
- Preserve backward compatibility unless a breaking change is explicitly approved.
- Do not edit generated SQLC code.
- Do not add dependencies for trivial functionality.
- Do not silently swallow errors.
- Do not write code that appears complete but contains placeholders.

## Before completion

- Run the available formatter, linter, tests, SQLC generation check, migration validation, and production build.
- Review the implementation against `06_FEATURE_REVIEW.md`.
- Report commands actually run.
- Report any check that could not be run.
- Report known risks and incomplete work.

## Clarification policy

Ask a question only when the missing information materially affects business behaviour, data integrity, authorization, or public API compatibility. Otherwise, follow the existing repository pattern and proceed.
