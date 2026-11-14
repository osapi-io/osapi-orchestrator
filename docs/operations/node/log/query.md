# LogQuery

Queries journal log entries on the target host.

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
| `opts`    | `LogQueryOpts` | Query options (refer to the table below).                 |

### LogQueryOpts

| Field      | Type      | Required | Description                                               |
| ---------- | --------- | -------- | --------------------------------------------------------- |
| `Lines`    | `*int`    | No       | The maximum number of log lines to return.                |
| `Since`    | `*string` | No       | Return only entries since this time (for example, `1h`).  |
| `Priority` | `*string` | No       | Filter by priority level (for example, `err`, `warning`). |

## Result type

```go
var result osapi.LogEntryResult
err := results.Decode("query-log-1", &result)
```

| Field      | Type         | Description                                                  |
| ---------- | ------------ | ------------------------------------------------------------ |
| `Hostname` | `string`     | The hostname of the agent that ran the operation.            |
| `Entries`  | `[]LogEntry` | List of journal entries.                                     |
| `Error`    | `string`     | The error message if the operation failed. Empty on success. |

### LogEntry

| Field       | Type     | Description             |
| ----------- | -------- | ----------------------- |
| `Timestamp` | `string` | The entry timestamp.    |
| `Unit`      | `string` | The systemd unit name.  |
| `Priority`  | `string` | The log priority level. |
| `Message`   | `string` | The log message.        |
| `PID`       | `int`    | The process ID.         |
| `Hostname`  | `string` | The source hostname.    |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `log:read` permission.

## Example

Refer to
[`examples/operations/log.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/log.go)
for a complete example.
