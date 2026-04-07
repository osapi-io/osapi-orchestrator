# LogQuery

Queries journal log entries on the target node.

## Usage

```go
lines := 50
step := o.LogQuery("web-01", osapi.LogQueryOpts{
    Lines: &lines,
})
```

## Parameters

| Parameter | Type           | Description                                               |
| --------- | -------------- | --------------------------------------------------------- |
| `target`  | `string`       | Target host: `_any`, `_all`, hostname, or label selector. |
| `opts`    | `LogQueryOpts` | Query options (see below).                                |

### LogQueryOpts

| Field      | Type      | Required | Description                                        |
| ---------- | --------- | -------- | -------------------------------------------------- |
| `Lines`    | `*int`    | No       | Maximum number of log lines to return.             |
| `Since`    | `*string` | No       | Filter entries since this time (e.g., `1h`).       |
| `Priority` | `*string` | No       | Filter by priority level (e.g., `err`, `warning`). |

## Result Type

```go
var result osapi.LogEntryResult
err := results.Decode("query-log-1", &result)
```

| Field      | Type         | Description                                      |
| ---------- | ------------ | ------------------------------------------------ |
| `Hostname` | `string`     | The node's hostname.                             |
| `Entries`  | `[]LogEntry` | List of journal entries.                         |
| `Error`    | `string`     | Error message if query failed; empty on success. |

### LogEntry

| Field       | Type     | Description         |
| ----------- | -------- | ------------------- |
| `Timestamp` | `string` | Entry timestamp.    |
| `Unit`      | `string` | Systemd unit name.  |
| `Priority`  | `string` | Log priority level. |
| `Message`   | `string` | Log message.        |
| `PID`       | `int`    | Process ID.         |
| `Hostname`  | `string` | Source hostname.    |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `log:read` permission.

## Example

See
[`examples/operations/log.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/log.go)
for a complete working example.
