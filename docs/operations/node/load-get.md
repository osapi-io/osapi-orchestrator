# NodeLoadGet

Gets the system load averages for the target host over 1-minute, 5-minute, and
15-minute intervals.

## Usage

```go
step := o.NodeLoadGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result type

```go
var result osapi.LoadResult
err := results.Decode("node.load.get-1", &result)
```

| Field    | Type      | Description             |
| -------- | --------- | ----------------------- |
| `Load1`  | `float32` | 1-minute load average.  |
| `Load5`  | `float32` | 5-minute load average.  |
| `Load15` | `float32` | 15-minute load average. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `node:read` permission.

## Example

Refer to
[`examples/operations/node-info.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/node-info.go)
for a complete example.
