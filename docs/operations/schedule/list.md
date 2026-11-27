# ScheduleList

List the cron entries that osapi deployed.

Lists each file in /etc/cron.d and in /etc/cron.hourly, daily, weekly and
monthly that osapi deployed, with its object, its cron expression or interval,
and its user. A cron file that osapi did not deploy is not listed.

## Usage

```go
step := o.ScheduleList("web-01")
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |

## Result

Each host answers with one `osapi.ScheduleEntry` row. A list answers with one
row for each item.

```go
var row osapi.ScheduleEntry
err := report.Decode("list-schedule", &row)
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
`GET /api/node/{hostname}/schedule`.

## Example

Refer to
[`examples/operations/schedule.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/schedule.go).
