# ProcessSignal

Sends a signal to a process on the target node.

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
| `opts`    | `ProcessSignalOpts` | Signal options (see below).                               |

### ProcessSignalOpts

| Field    | Type     | Required | Description                                |
| -------- | -------- | -------- | ------------------------------------------ |
| `Signal` | `string` | Yes      | Signal name (e.g., `TERM`, `KILL`, `HUP`). |

## Result Type

```go
var result osapi.ProcessSignalResult
err := results.Decode("signal-process-1", &result)
```

| Field     | Type     | Description                                       |
| --------- | -------- | ------------------------------------------------- |
| `PID`     | `int`    | Process ID that was signaled.                     |
| `Signal`  | `string` | Signal that was sent.                             |
| `Changed` | `bool`   | Whether the signal was delivered.                 |
| `Error`   | `string` | Error message if signal failed; empty on success. |

## Idempotency

**Non-idempotent.** Always delivers the signal regardless of current state.

## Permissions

Requires `process:execute` permission.

## Example

See
[`examples/operations/process.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/process.go)
for a complete working example.
