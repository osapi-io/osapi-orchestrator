# ProcessList

Lists running processes on the target node.

## Usage

```go
step := o.ProcessList("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result Type

```go
var result osapi.ProcessInfoResult
err := results.Decode("list-process-1", &result)
```

| Field       | Type            | Description                                      |
| ----------- | --------------- | ------------------------------------------------ |
| `Hostname`  | `string`        | The node's hostname.                             |
| `Processes` | `[]ProcessInfo` | List of running processes.                       |
| `Error`     | `string`        | Error message if query failed; empty on success. |

### ProcessInfo

| Field        | Type      | Description                |
| ------------ | --------- | -------------------------- |
| `PID`        | `int`     | Process ID.                |
| `Name`       | `string`  | Process name.              |
| `User`       | `string`  | User running the process.  |
| `State`      | `string`  | Process state.             |
| `CPUPercent` | `float64` | CPU usage percentage.      |
| `MemPercent` | `float32` | Memory usage percentage.   |
| `MemRSS`     | `int64`   | Resident set size (bytes). |
| `Command`    | `string`  | Full command line.         |
| `StartTime`  | `string`  | Process start time.        |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `process:read` permission.

## Example

See
[`examples/operations/process.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/process.go)
for a complete working example.
