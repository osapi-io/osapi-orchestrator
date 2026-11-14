# Hardware

Get information about hardware resources: disk use and memory statistics.

## Operations

| Method                                   | Description                 | Idempotent |
| ---------------------------------------- | --------------------------- | ---------- |
| [`NodeDiskGet(target)`](disk-get.md)     | Get the disk use statistics | Read-only  |
| [`NodeMemoryGet(target)`](memory-get.md) | Get the memory statistics   | Read-only  |

## Permissions

| Operation       | Permission  |
| --------------- | ----------- |
| Read operations | `node:read` |

## Examples

Refer to
[`examples/operations/node-info.go`](../../../../examples/operations/node-info.go)
for a read-only workflow that gets the disk and memory statistics.
