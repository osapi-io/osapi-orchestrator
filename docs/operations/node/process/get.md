# ProcessGet

Gets a specific process by PID from the target host.

## Usage

```go
step := o.ProcessGet("web-01", 1234)
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `pid`     | `int`    | The process ID to get.                                    |

## Result type

```go
var result osapi.ProcessInfoResult
err := results.Decode("get-process-1", &result)
```

| Field       | Type            | Description                                                  |
| ----------- | --------------- | ------------------------------------------------------------ |
| `Hostname`  | `string`        | The hostname of the agent that ran the operation.            |
| `Processes` | `[]ProcessInfo` | The process details (one entry).                             |
| `Error`     | `string`        | The error message if the operation failed. Empty on success. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `process:read` permission.

## Example

Refer to
[`examples/operations/process.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/process.go)
for a complete example.
