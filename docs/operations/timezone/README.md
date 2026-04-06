# Timezone Management

Manage the system timezone on target nodes -- get the current timezone and
update it.

## Operations

| Method                                          | Description              | Idempotent |
| ----------------------------------------------- | ------------------------ | ---------- |
| [`TimezoneGet(target)`](get.md)                 | Get the system timezone  | Read-only  |
| [`TimezoneUpdate(target, opts)`](update.md)     | Set the system timezone  | Idempotent |

## Permissions

| Operation        | Permission       |
| ---------------- | ---------------- |
| Read operations  | `timezone:read`  |
| Write operations | `timezone:write` |

## Example

See
[`examples/operations/timezone.go`](../../examples/operations/timezone.go) for
a complete workflow example.
