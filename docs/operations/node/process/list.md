# ProcessList

Lists the processes on the target host.

## Usage

```go
step := o.ProcessList("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result type

```go
var result osapi.ProcessInfoResult
err := results.Decode("list-process-1", &result)
```

| Field       | Type            | Description                                                  |
| ----------- | --------------- | ------------------------------------------------------------ |
| `Hostname`  | `string`        | The hostname of the agent that ran the operation.            |
| `Processes` | `[]ProcessInfo` | List of processes.                                           |
| `Error`     | `string`        | The error message if the operation failed. Empty on success. |

### ProcessInfo

| Field        | Type      | Description                     |
| ------------ | --------- | ------------------------------- |
| `PID`        | `int`     | The process ID.                 |
| `Name`       | `string`  | The process name.               |
| `User`       | `string`  | The user that runs the process. |
| `State`      | `string`  | The process state.              |
| `CPUPercent` | `float64` | The CPU usage percentage.       |
| `MemPercent` | `float32` | The memory usage percentage.    |
| `MemRSS`     | `int64`   | The resident set size (bytes).  |
| `Command`    | `string`  | The full command line.          |
| `StartTime`  | `string`  | The process start time.         |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `process:read` permission.

## Example

Refer to
[`examples/operations/process.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/process.go)
for a complete example.
