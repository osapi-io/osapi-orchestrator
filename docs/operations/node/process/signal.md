# ProcessSignal

Sends a signal to a process on the target host.

## Usage

```go
step := o.ProcessSignal("web-01", 1234, osapi.ProcessSignalOpts{
    Signal: "TERM",
})
```

## Parameters

| Parameter | Type                | Description                                               |
| --------- | ------------------- | --------------------------------------------------------- |
| `target`  | `string`            | Target host: `_any`, `_all`, hostname, or label selector. |
| `pid`     | `int`               | Process ID to signal.                                     |
| `opts`    | `ProcessSignalOpts` | Signal options (refer to the table below).                |

### ProcessSignalOpts

| Field    | Type     | Required | Description                                       |
| -------- | -------- | -------- | ------------------------------------------------- |
| `Signal` | `string` | Yes      | Signal name (for example, `TERM`, `KILL`, `HUP`). |

## Result type

```go
var result osapi.ProcessSignalResult
err := results.Decode("signal-process-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `PID`     | `int`    | Process ID that received the signal.                         |
| `Signal`  | `string` | Signal that the agent sent.                                  |
| `Changed` | `bool`   | `true` if the operation sent the signal.                     |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Non-idempotent.** The operation sends the signal each time the step runs. It
does not examine the process state first.

## Permissions

The token must have the `process:execute` permission.

## Example

Refer to
[`examples/operations/process.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/process.go)
for a complete example.
