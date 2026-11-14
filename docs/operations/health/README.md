# Health

Check the connection to the OSAPI API server. A small liveness probe makes sure
that the server is reachable. The probe does not use the job system.

## Operations

| Method                      | Description                      | Idempotent |
| --------------------------- | -------------------------------- | ---------- |
| [`HealthCheck()`](check.md) | Liveness probe on the API server | Read-only  |

## Permissions

| Operation    | Permission    |
| ------------ | ------------- |
| Health check | `health:read` |

## Example

Refer to [`examples/features/basic.go`](../../../examples/features/basic.go) for
an example that uses `HealthCheck` as a gate before other operations run.
