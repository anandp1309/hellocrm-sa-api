# Start Here

This folder defines how all EZ Solutions software must be engineered. It is product-independent and reusable across current and future repositories.

It does **not** define product functionality. Product workflows, roles, screens, business rules, field definitions, and acceptance criteria belong in the project's own documentation.

## Developer workflow

### At the beginning of a new AI session

Use:

```text
Read `AI/00_SYSTEM_PROMPT.md` and `AI/01_START_HERE.md` completely.
Then inspect the repository, explain what you understood, identify the relevant standards, and prepare an implementation plan before changing code.
```

### Before every feature

Use:

```text
Read `AI/05_CHECKPOINT.md`, inspect the related implementation, and prepare the implementation plan for this feature:

[feature details]
```

### Before accepting a feature

Use:

```text
Read `AI/06_FEATURE_REVIEW.md`.
Review the complete implementation, fix all deviations, run the available checks, and report:
1. Files changed
2. Standards verified
3. Checks executed
4. Remaining risks or incomplete items
```

## Reading order

Always read:

1. `02_STACK.md`
2. `03_RULES.md`
3. `04_AI_RULES.md`

Then read the standards relevant to the task.

| Task | Read |
|---|---|
| Backend feature | `architecture/coding.md`, `architecture/errors.md`, `architecture/api.md` |
| Database change | `architecture/database.md`, `reference/migration-example.md`, `reference/sqlc-example.md` |
| Authentication or permissions | `architecture/security.md`, `building-blocks/auth-session.md`, `architecture/tenant-isolation.md` |
| Frontend feature | `architecture/ui.md`, relevant building block |
| Deployment or infrastructure | `architecture/deployment.md`, `architecture/observability.md` |
| Approval flow | `building-blocks/approval-workflow.md` |
| CRUD/master screen | `building-blocks/crud.md` |
| Import/export | `building-blocks/import-export.md` |
| File upload | `building-blocks/file-upload.md` |
| Number generation | `architecture/database.md`, `building-blocks/number-series.md` |

## Core principle

The AI agent and developer should spend minimal time deciding *how* to work. The standards define the approved engineering approach; project requirements define *what* to build.
