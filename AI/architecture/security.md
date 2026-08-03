# Security Standard

Status: Required

Use current OWASP ASVS guidance as the verification baseline.

## Authentication

- Use an approved adaptive password-hashing algorithm.
- Use short-lived access tokens or secure server sessions.
- Rotate and revoke refresh tokens where refresh tokens are used.
- Rate-limit login, password reset, and OTP endpoints.
- OTPs require expiry, retry limits, and replay protection.
- Never log passwords, OTPs, access tokens, refresh tokens, or reset tokens.

## Authorization

- Enforce permissions on the backend.
- Validate tenant/project/workspace membership.
- Validate object-level ownership.
- Use deny-by-default.
- Do not trust IDs supplied by the client without scope checks.
- Sensitive actions may require stronger confirmation or re-authentication.

## Input and output

- Validate all untrusted input.
- Use parameterised SQL through SQLC.
- Encode output in the correct context.
- Prevent XSS, CSRF, SSRF, path traversal, and unsafe redirects.
- Validate uploaded file size, type, extension, and expected content.
- Do not trust MIME type supplied by the browser.

## Secrets

- Keep secrets outside source control.
- Do not bake secrets into images.
- Use environment or approved secret management.
- Rotate exposed secrets.
- Separate credentials by environment.

## Headers and transport

- Use TLS.
- Apply secure cookie attributes.
- Use a controlled CORS policy.
- Apply appropriate security headers.
- Avoid wildcard origins for authenticated APIs.

## Dependencies

- Use maintained dependencies.
- Scan for vulnerabilities.
- Review security updates promptly.
- Do not auto-merge major dependency changes without testing.

## Logging

- Mask sensitive personal and financial values.
- Include enough context for investigation.
- Do not expose internal errors to clients.
