# Hostname

Get and update the system hostname.

## Operations

| Method                                              | Description                        | Idempotent |
| --------------------------------------------------- | ---------------------------------- | ---------- |
| [`NodeHostnameGet(target)`](get.md)                 | Get the system hostname and labels | Read-only  |
| [`NodeHostnameUpdate(target, hostname)`](update.md) | Set the system hostname            | Idempotent |

## Permissions

| Operation        | Permission   |
| ---------------- | ------------ |
| Read operations  | `node:read`  |
| Write operations | `node:write` |

## Examples

Refer to
[`examples/operations/hostname-update.go`](../../../../examples/operations/hostname-update.go)
for a hostname workflow that reads, then writes, and uses a broadcast target.
