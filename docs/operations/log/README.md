# Log Management

Query systemd journal logs on target nodes -- query all entries, list log
sources, and query by unit name.

## Operations

| Method                                              | Description                | Idempotent |
| --------------------------------------------------- | -------------------------- | ---------- |
| [`LogQuery(target, opts)`](query.md)                | Query journal log entries  | Read-only  |
| [`LogSources(target)`](sources.md)                  | List available log sources | Read-only  |
| [`LogQueryUnit(target, unit, opts)`](query-unit.md) | Query logs for a unit      | Read-only  |

## Permissions

| Operation      | Permission |
| -------------- | ---------- |
| All operations | `log:read` |

## Example

See [`examples/operations/log.go`](../../examples/operations/log.go) for a
complete workflow example covering all operations.
