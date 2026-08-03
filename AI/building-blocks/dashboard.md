# Dashboard Building Block

## Design order

1. Critical operational alerts
2. Primary KPIs
3. Trends or comparisons
4. Actionable lists
5. Supporting detail

## Rules

- Do not overload with charts.
- Every metric must have a defined calculation.
- Show the reporting period.
- Use cached or pre-aggregated data only when measured need justifies it.
- Avoid running many expensive count queries per page request.
- Scope all metrics to the actor's access.
- Empty dashboards should explain the next meaningful action.
- Numbers shown in summary and detail must reconcile.
