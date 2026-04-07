# Node

Query node-level runtime information -- status, load averages, uptime, OS
details, and manage power and processes.

## Operations

| Method                                        | Description                     | Idempotent |
| --------------------------------------------- | ------------------------------- | ---------- |
| [`NodeStatusGet(target)`](status-get.md)      | Get node registration status    | Read-only  |
| [`NodeLoadGet(target)`](load-get.md)          | Get load averages               | Read-only  |
| [`NodeUptimeGet(target)`](uptime-get.md)      | Get system uptime               | Read-only  |
| [`NodeOSGet(target)`](os-get.md)              | Get OS distribution and version | Read-only  |

See also: [Power](power/), [Process](process/), [Log](log/)

## Permissions

| Operation        | Permission   |
| ---------------- | ------------ |
| Read operations  | `node:read`  |

## Examples

See [`examples/operations/node-info.go`](../../examples/operations/node-info.go)
for a read-only workflow covering status, load, uptime, and OS queries.
