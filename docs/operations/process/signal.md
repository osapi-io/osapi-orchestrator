# ProcessSignal

Send a signal to a process.

Sends TERM, KILL, HUP, INT, USR1 or USR2 to one process with kill, and reports
changed: true. A process that the host does not have is not found. A process ID
below 1 is refused, because kill sends to a group of processes for 0 and for a
negative ID.

## Usage

```go
step := o.ProcessSignal("web-01", input)
```

## Parameters

| Parameter | Type                       | Description                                                   |
| --------- | -------------------------- | ------------------------------------------------------------- |
| `target`  | `string`                   | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.ProcessSignalInput` | Refer to `ProcessSignalInput` below.                          |

### ProcessSignalInput

| Field    | Type           | Description                                                                       |
| -------- | -------------- | --------------------------------------------------------------------------------- |
| `PID`    | `int`          | The pid of the path: the process ID, such as 1234. It names the item in the path. |
| `Signal` | `types.Signal` | Signal is a signal that osapi sends to a process, by its name without SIG.        |

## Result

Each host answers with one `osapi.ProcessResult` row. A list answers with one
row for each item.

```go
var row osapi.ProcessResult
err := report.Decode("signal-process", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`  | `bool`            | Whether the signal was sent.                                                                                                                                                                                             |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `PID`      | `int`             | The process ID.                                                                                                                                                                                                          |
| `Signal`   | `types.Signal`    | Signal is a signal that osapi sends to a process, by its name without SIG.                                                                                                                                               |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Retries

A second run is not the same as one run. So `Retry` does not run this step again
unless the step adds `RetryUnsafe`.

## Permissions

The token must have the `process:execute` permission. The step sends
`POST /api/node/{hostname}/process/{pid}/signal`.

## Example

Refer to
[`examples/operations/process.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/process.go).
