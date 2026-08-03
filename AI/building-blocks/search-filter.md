# Search, Filter, Sort, and Pagination Building Block

## Search

- Define searchable fields explicitly.
- Normalise input where the domain requires it.
- Use database indexes appropriate to actual search.
- Do not use unbounded `%term%` scans on large tables without review.

## Filters

- Whitelist filter names.
- Validate values.
- Keep tenant/project scope mandatory.
- Use date ranges with clear inclusive/exclusive behaviour.

## Sorting

- Whitelist sort fields.
- Define default sort.
- Add stable ID tiebreaker.
- Never interpolate raw user input into `ORDER BY`.

## Pagination

- Enforce maximum `per_page`.
- Return total counts only when needed; avoid expensive counts on very large datasets.
- Use cursor pagination for large, high-change streams.

## Export

Exports must apply the same permission and filter rules as the list endpoint.
