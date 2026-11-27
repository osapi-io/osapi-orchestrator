# ScheduleCreate

Deploy a cron entry.

Gets the object from the object store and writes it as /etc/cron.d/backup for
backup, mode 0644, for a cron expression, or as /etc/cron.daily/backup, mode
0755, for a daily interval. An object of content type template, or content_type
template, is rendered with the agent's facts and the vars. A file that has the
content does not change. An entry of the name in a different directory does not
change: move it with delete and create. A file at the path that osapi did not
deploy is refused. An entry runs as its user, root included, so only admin can
deploy one.

## Usage

```go
step := o.ScheduleCreate("web-01", input)
```

## Parameters

| Parameter | Type                        | Description                                                   |
| --------- | --------------------------- | ------------------------------------------------------------- |
| `target`  | `string`                    | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.ScheduleCreateInput` | Refer to `ScheduleCreateInput` below.                         |

### ScheduleCreateInput

| Field         | Type                | Description                                                                                                                                   |
| ------------- | ------------------- | --------------------------------------------------------------------------------------------------------------------------------------------- |
| `ContentType` | `types.ContentType` | ContentType is how a deploy treats a file: raw writes its bytes as they are; template renders it with Go text/template and the agent's facts. |
| `Interval`    | `types.Interval`    | Interval is how often a scheduled job runs, in place of a cron schedule.                                                                      |
| `Name`        | `string`            | The entry, and its file name, such as backup.                                                                                                 |
| `Object`      | `string`            | The object that holds the entry, from osapi client file upload.                                                                               |
| `Schedule`    | `string`            | A cron expression, such as \*/5 * * * \*; give it or interval.                                                                                |
| `User`        | `string`            | The user the entry runs as.                                                                                                                   |
| `Vars`        | `map[string]string` | The vars of a template, as key=value.                                                                                                         |

## Result

Each host answers with one `osapi.ScheduleResult` row. A list answers with one
row for each item.

```go
var row osapi.ScheduleResult
err := report.Decode("create-schedule", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`  | `bool`            | Whether the host changed.                                                                                                                                                                                                |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Name`     | `string`          | The entry.                                                                                                                                                                                                               |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `schedule:write` permission. The step sends
`POST /api/node/{hostname}/schedule`.

## Example

Refer to
[`examples/operations/schedule.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/schedule.go).
