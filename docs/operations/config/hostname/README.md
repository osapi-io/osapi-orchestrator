# Hostname

Query and update the system hostname.

## Operations

| Method                                                       | Description                    | Idempotent |
| ------------------------------------------------------------ | ------------------------------ | ---------- |
| [`NodeHostnameGet(target)`](get.md)                          | Get system hostname and labels | Read-only  |
| [`NodeHostnameUpdate(target, hostname)`](update.md)          | Set system hostname            | Idempotent |

## Permissions

| Operation        | Permission   |
| ---------------- | ------------ |
| Read operations  | `node:read`  |
| Write operations | `node:write` |

## Examples

See
[`examples/operations/hostname-update.go`](../../../examples/operations/hostname-update.go)
for a read-then-write hostname workflow with broadcast targeting.
