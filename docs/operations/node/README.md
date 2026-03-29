# Node Management

Query and update host-level system information -- hostname, disk usage,
memory statistics, load averages, uptime, OS details, and overall node
status.

## Operations

| Method | Description | Idempotent |
| ------ | ----------- | ---------- |
| [`NodeHostnameGet(target)`](hostname-get.md) | Get system hostname and labels | Read-only |
| [`NodeHostnameUpdate(target, hostname)`](hostname-update.md) | Set system hostname | Idempotent |
| [`NodeStatusGet(target)`](status-get.md) | Get node registration status | Read-only |
| [`NodeDiskGet(target)`](disk-get.md) | Get disk usage statistics | Read-only |
| [`NodeMemoryGet(target)`](memory-get.md) | Get memory statistics | Read-only |
| [`NodeLoadGet(target)`](load-get.md) | Get load averages | Read-only |
| [`NodeUptimeGet(target)`](uptime-get.md) | Get system uptime | Read-only |
| [`NodeOSGet(target)`](os-get.md) | Get OS distribution and version | Read-only |

## Permissions

| Operation | Permission |
| --------- | ---------- |
| Read operations | `node:read` |
| Write operations | `node:write` |

## Examples

See
[`examples/operations/node-info.go`](../../examples/operations/node-info.go)
for a read-only workflow covering hostname, disk, memory, load, uptime,
and OS queries.

See
[`examples/operations/hostname-update.go`](../../examples/operations/hostname-update.go)
for a read-then-write hostname workflow with broadcast targeting.
