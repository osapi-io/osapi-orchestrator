# LogQueryUnit

Queries journal log entries for a specific systemd unit on the target host.

## Usage

```go
lines := 100
step := o.LogQueryUnit("web-01", "nginx", osapi.LogQueryOpts{
    Lines: &lines,
})
```

## Parameters

| Parameter | Type           | Description                                               |
| --------- | -------------- | --------------------------------------------------------- |
| `target`  | `string`       | Target host: `_any`, `_all`, hostname, or label selector. |
| `unit`    | `string`       | The systemd unit name to query logs for.                  |
| `opts`    | `LogQueryOpts` | Query options (refer to [LogQuery](query.md)).            |

## Result type

```go
var result osapi.LogEntryResult
err := results.Decode("query-log-unit-1", &result)
```

| Field      | Type         | Description                                                  |
| ---------- | ------------ | ------------------------------------------------------------ |
| `Hostname` | `string`     | The hostname of the agent that ran the operation.            |
| `Entries`  | `[]LogEntry` | List of journal entries for the unit.                        |
| `Error`    | `string`     | The error message if the operation failed. Empty on success. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `log:read` permission.

## Example

Refer to
[`examples/operations/log.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/log.go)
for a complete example.
