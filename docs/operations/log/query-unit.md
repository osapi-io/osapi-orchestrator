# LogQueryUnit

Queries journal log entries for a specific systemd unit on the target node.

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
| `unit`    | `string`       | Systemd unit name to query logs for.                      |
| `opts`    | `LogQueryOpts` | Query options (see [LogQuery](query.md)).                  |

## Result Type

```go
var result osapi.LogEntryResult
err := results.Decode("query-log-unit-1", &result)
```

| Field      | Type         | Description                                      |
| ---------- | ------------ | ------------------------------------------------ |
| `Hostname` | `string`     | The node's hostname.                             |
| `Entries`  | `[]LogEntry` | List of journal entries for the unit.            |
| `Error`    | `string`     | Error message if query failed; empty on success. |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `log:read` permission.

## Example

See
[`examples/operations/log.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/log.go)
for a complete working example.
