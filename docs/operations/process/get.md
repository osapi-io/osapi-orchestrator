# ProcessGet

Get one process.

Reads one process by its ID, without the others. A process that the host does
not have is not found.

## Usage

```go
step := o.ProcessGet("web-01", pid)
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `pid`     | `int`    | The process ID.                                               |

## Result

Each host answers with one `osapi.ProcessEntry` row. A list answers with one row
for each item.

```go
var row osapi.ProcessEntry
err := report.Decode("get-process", &row)
```

| Field        | Type              | Description                                                                                                                                                                                                              |
| ------------ | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Command`    | `string`          | Its full command line.                                                                                                                                                                                                   |
| `CPUPercent` | `float64`         | Its CPU time as a percent of the time since it started.                                                                                                                                                                  |
| `Error`      | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname`   | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `MemPercent` | `float64`         | Its resident memory as a percent of the host's memory.                                                                                                                                                                   |
| `MemRSS`     | `int64`           | Its resident memory in bytes.                                                                                                                                                                                            |
| `Name`       | `string`          | The name of the process.                                                                                                                                                                                                 |
| `PID`        | `int`             | The process ID.                                                                                                                                                                                                          |
| `StartTime`  | `string`          | When it started, in RFC 3339.                                                                                                                                                                                            |
| `State`      | `string`          | Its state, such as running or sleep.                                                                                                                                                                                     |
| `Status`     | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |
| `User`       | `string`          | The user that runs it.                                                                                                                                                                                                   |

## Permissions

The token must have the `process:read` permission. The step sends
`GET /api/node/{hostname}/process/{pid}`.

## Example

Refer to
[`examples/operations/process.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/process.go).
