# ProcessGet

Retrieves a specific process by PID from the target node.

## Usage

```go
step := o.ProcessGet("web-01", 1234)
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `pid`     | `int`    | Process ID to retrieve.                                   |

## Result Type

```go
var result osapi.ProcessInfoResult
err := results.Decode("get-process-1", &result)
```

| Field       | Type            | Description                                      |
| ----------- | --------------- | ------------------------------------------------ |
| `Hostname`  | `string`        | The node's hostname.                             |
| `Processes` | `[]ProcessInfo` | Process details (single entry).                  |
| `Error`     | `string`        | Error message if query failed; empty on success. |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `process:read` permission.

## Example

See
[`examples/operations/process.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/process.go)
for a complete working example.
