# Hardware

Query hardware resource information -- disk usage and memory statistics.

## Operations

| Method                                   | Description               | Idempotent |
| ---------------------------------------- | ------------------------- | ---------- |
| [`NodeDiskGet(target)`](disk-get.md)     | Get disk usage statistics | Read-only  |
| [`NodeMemoryGet(target)`](memory-get.md) | Get memory statistics     | Read-only  |

## Permissions

| Operation       | Permission  |
| --------------- | ----------- |
| Read operations | `node:read` |

## Examples

See [`examples/operations/node-info.go`](../../../examples/operations/node-info.go)
for a read-only workflow covering disk and memory queries.
