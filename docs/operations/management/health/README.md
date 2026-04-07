# Health

Check OSAPI API server connectivity -- a lightweight liveness probe that
confirms the server is reachable without going through the job system.

## Operations

| Method                      | Description                           | Idempotent |
| --------------------------- | ------------------------------------- | ---------- |
| [`HealthCheck()`](check.md) | Liveness probe against the API server | Read-only  |

## Permissions

| Operation    | Permission    |
| ------------ | ------------- |
| Health check | `health:read` |

## Example

See [`examples/features/basic.go`](../../examples/features/basic.go) for a
working example that uses `HealthCheck` as a gate before running other
operations.
