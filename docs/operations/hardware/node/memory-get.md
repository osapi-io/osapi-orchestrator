# NodeMemoryGet

Gets the memory statistics of the target host: total, free, and cached memory.

## Usage

```go
step := o.NodeMemoryGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result type

```go
var result osapi.MemoryResult
err := results.Decode("node.memory.get-1", &result)
```

| Field    | Type     | Description                  |
| -------- | -------- | ---------------------------- |
| `Total`  | `uint64` | The total memory, in bytes.  |
| `Free`   | `uint64` | The free memory, in bytes.   |
| `Cached` | `uint64` | The cached memory, in bytes. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `node:read` permission.

## Example

Refer to
[`examples/operations/node-info.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/node-info.go)
for a complete example.
