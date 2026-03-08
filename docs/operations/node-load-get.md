# NodeLoadGet

Retrieves system load averages for the target node over 1-minute, 5-minute, and
15-minute intervals.

## Usage

```go
step := o.NodeLoadGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result Type

```go
var result orchestrator.LoadResult
err := results.Decode("node.load.get-1", &result)
```

| Field    | Type      | Description             |
| -------- | --------- | ----------------------- |
| `Load1`  | `float32` | 1-minute load average.  |
| `Load5`  | `float32` | 5-minute load average.  |
| `Load15` | `float32` | 15-minute load average. |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `node:read` permission.

## Example

See
[`examples/features/parallel.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/parallel.go)
for a complete working example.
