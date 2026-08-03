# SvelteKit, TypeScript, and UI Standard

Status: Required

## TypeScript

- Enable strict mode.
- Avoid `any`.
- Type component props, form data, API requests, and responses.
- Document unavoidable unsafe casts.
- Keep server-only types and secrets out of browser bundles.

## SvelteKit responsibilities

Use:

- `+page.server.ts` for authenticated server-side loading and actions where appropriate
- Form actions for normal server-validated forms when they improve reliability
- Client API calls for interactions that genuinely require them
- Hooks for approved cross-cutting request/session behaviour

Do not recreate legacy SPA patterns when SvelteKit provides a supported server pattern.

## API access

- Use one typed API client.
- Centralise authentication headers, request IDs, error parsing, and retry policy.
- Do not call `fetch` ad hoc throughout components.
- Do not duplicate the API error contract.

## Components

- Reuse shared components.
- Keep generic components free of product-specific business rules.
- Keep page components focused.
- Extract sections when complexity improves readability, not simply to reduce line count.
- Avoid excessive global stores.

## Forms

Every form must have:

- Persistent labels
- Client-side usability validation
- Authoritative server validation
- Field-level error display
- Submission loading state
- Duplicate-submit protection
- Preserved input after recoverable failure
- Clear success feedback

## Tables

- Use server-side pagination for large data.
- Whitelist filters and sort fields.
- Keep actions predictable.
- Avoid excessive columns.
- Provide mobile alternatives where tables become unusable.
- Export only when required.

## UI states

Every asynchronous page or component must consider:

- Initial loading
- Empty data
- Validation error
- Server error
- Permission denied
- Success
- Disabled/submitting state

## Permissions

Frontend permission checks control visibility and UX only. Backend authorization remains mandatory.

## Accessibility

- Use semantic HTML.
- Associate labels with fields.
- Support keyboard navigation.
- Use visible focus states.
- Do not convey meaning through colour alone.
- Provide accessible names for icon-only controls.

## Responsive behaviour

- Design mobile behaviour intentionally.
- Do not merely shrink desktop layouts.
- Keep destructive and high-risk actions safe on small screens.
