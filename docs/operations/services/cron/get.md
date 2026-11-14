# CronGet

Gets a specific cron entry by name from the target host.

## Usage

```go
step := o.CronGet("web-01", "backup")
```

## Parameters

| Parameter   | Type     | Description                                               |
| ----------- | -------- | --------------------------------------------------------- |
| `target`    | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `entryName` | `string` | Name of the cron entry to get.                            |

## Result type

```go
var result osapi.CronEntryResult
err := results.Decode("get-cron-1", &result)
```

| Field      | Type     | Description                                                  |
| ---------- | -------- | ------------------------------------------------------------ |
| `Hostname` | `string` | The hostname of the agent that ran the operation.            |
| `Name`     | `string` | Cron entry name.                                             |
| `Object`   | `string` | Object store name of the deployed script.                    |
| `Schedule` | `string` | Cron expression. Empty if the entry uses an interval.        |
| `Interval` | `string` | Periodic interval: hourly, daily, weekly, or monthly.        |
| `Source`   | `string` | Path to the deployed script on the filesystem.               |
| `User`     | `string` | User the cron job runs as.                                   |
| `Error`    | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `cron:read` permission.

## Example

Refer to
[`examples/operations/cron.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/cron.go)
for a complete example.
