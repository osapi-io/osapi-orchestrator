# NodeMemoryGet

Retrieves memory statistics for the target node, including total, free, and
cached memory.

## Usage

```go
step := o.NodeMemoryGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result Type

```go
var result orchestrator.MemoryResult
err := results.Decode("node.memory.get-1", &result)
```

| Field    | Type     | Description             |
| -------- | -------- | ----------------------- |
| `Total`  | `uint64` | Total memory in bytes.  |
| `Free`   | `uint64` | Free memory in bytes.   |
| `Cached` | `uint64` | Cached memory in bytes. |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `node:read` permission.

## Example

See
[`examples/parallel.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/parallel.go)
for a complete working example.
