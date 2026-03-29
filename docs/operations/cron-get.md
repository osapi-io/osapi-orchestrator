# CronGet

Retrieves a specific cron entry by name from the target node.

## Usage

```go
step := o.CronGet("web-01", "backup")
```

## Parameters

| Parameter   | Type     | Description                                               |
| ----------- | -------- | --------------------------------------------------------- |
| `target`    | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `entryName` | `string` | Name of the cron entry to retrieve.                       |

## Result Type

```go
var result osapi.CronEntryResult
err := results.Decode("get-cron-1", &result)
```

| Field      | Type     | Description                                            |
| ---------- | -------- | ------------------------------------------------------ |
| `Hostname` | `string` | The node's hostname.                                   |
| `Name`     | `string` | Cron entry name.                                       |
| `Object`   | `string` | Object store name of the deployed script.              |
| `Schedule` | `string` | Cron expression (empty if interval-based).             |
| `Interval` | `string` | Periodic interval: hourly, daily, weekly, or monthly.  |
| `Source`   | `string` | Path to the deployed script on the filesystem.         |
| `User`     | `string` | User the cron job runs as.                             |
| `Error`    | `string` | Error message if query failed; empty on success.       |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `cron:read` permission.

## Example

See
[`examples/operations/cron.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/cron.go)
for a complete working example.
