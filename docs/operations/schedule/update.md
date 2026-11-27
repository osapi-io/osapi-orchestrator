# ScheduleUpdate

Deploy a cron entry that osapi deployed again.

Writes a new object, or the entry's own, at the entry's path, with the fields
that the request gives and the entry's own for the rest. A file that has the
content does not change. An entry that osapi did not deploy is not found. The
entry stays in its directory: an interval does not change with update.

## Usage

```go
step := o.ScheduleUpdate("web-01", input)
```

## Parameters

| Parameter | Type                        | Description                                                   |
| --------- | --------------------------- | ------------------------------------------------------------- |
| `target`  | `string`                    | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.ScheduleUpdateInput` | Refer to `ScheduleUpdateInput` below.                         |

### ScheduleUpdateInput

| Field         | Type                | Description                                                                                                                                   |
| ------------- | ------------------- | --------------------------------------------------------------------------------------------------------------------------------------------- |
| `Name`        | `string`            | The name of the path: the entry, such as backup. It names the item in the path.                                                               |
| `ContentType` | `types.ContentType` | ContentType is how a deploy treats a file: raw writes its bytes as they are; template renders it with Go text/template and the agent's facts. |
| `Object`      | `string`            | The object to deploy; the entry's own when empty.                                                                                             |
| `Schedule`    | `string`            | A new cron expression; the entry's own when empty.                                                                                            |
| `User`        | `string`            | The user the entry runs as; the entry's own when empty.                                                                                       |
| `Vars`        | `map[string]string` | The vars of a template, as key=value.                                                                                                         |

## Result

Each host answers with one `osapi.ScheduleResult` row. A list answers with one
row for each item.

```go
var row osapi.ScheduleResult
err := report.Decode("update-schedule", &row)
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
`PUT /api/node/{hostname}/schedule/{name}`.

## Example

Refer to
[`examples/operations/schedule.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/schedule.go).
