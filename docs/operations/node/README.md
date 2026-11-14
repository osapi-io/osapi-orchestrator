# Node

Query runtime information about target hosts: status, load averages, uptime, and
OS details. The related domains manage power, processes, and logs.

## Operations

| Method                                   | Description                     | Idempotent |
| ---------------------------------------- | ------------------------------- | ---------- |
| [`NodeStatusGet(target)`](status-get.md) | Get host registration status    | Read-only  |
| [`NodeLoadGet(target)`](load-get.md)     | Get load averages               | Read-only  |
| [`NodeUptimeGet(target)`](uptime-get.md) | Get system uptime               | Read-only  |
| [`NodeOSGet(target)`](os-get.md)         | Get OS distribution and version | Read-only  |

Related domains: [Power](power/), [Process](process/), [Log](log/)

## Permissions

| Operation       | Permission  |
| --------------- | ----------- |
| Read operations | `node:read` |

## Examples

Refer to
[`examples/operations/node-info.go`](../../../examples/operations/node-info.go)
for a read-only workflow that queries status, load, uptime, and OS.
