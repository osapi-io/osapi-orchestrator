# ScheduleGet

Get one cron entry that osapi deployed.

Reads the entry from /etc/cron.d, then from each interval directory. An entry
that osapi did not deploy is not found.

## Usage

```go
step := o.ScheduleGet("web-01", name)
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `name`    | `string` | The name of the item.                                         |

## Result

Each host answers with one `osapi.ScheduleEntry` row. A list answers with one
row for each item.

```go
var row osapi.ScheduleEntry
err := report.Decode("get-schedule", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Interval` | `string`          | Its interval, for a script in /etc/cron.daily and the like.                                                                                                                                                              |
| `Name`     | `string`          | The entry.                                                                                                                                                                                                               |
| `Object`   | `string`          | The object it was deployed from.                                                                                                                                                                                         |
| `Schedule` | `string`          | Its cron expression, for an entry in /etc/cron.d.                                                                                                                                                                        |
| `Source`   | `string`          | Where its file is: cron.d, or the interval.                                                                                                                                                                              |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |
| `User`     | `string`          | The user it runs as.                                                                                                                                                                                                     |

## Permissions

The token must have the `schedule:read` permission. The step sends
`GET /api/node/{hostname}/schedule/{name}`.

## Example

Refer to
[`examples/operations/schedule.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/schedule.go).
