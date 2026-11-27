# LogList

List the entries of the journal.

Runs journalctl --output=json, the newest entries last: the last 100 unless
lines says, of one unit when unit says, since a time when since says, and at a
priority or above when priority says. A priority is its name, such as err or
info. The hand-written log domain read the Debian family only; this reads each
Linux host that has journald.

## Usage

```go
step := o.LogList("web-01", params)
```

## Parameters

| Parameter | Type                  | Description                                                   |
| --------- | --------------------- | ------------------------------------------------------------- |
| `target`  | `string`              | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `params`  | `osapi.LogListParams` | Refer to `LogListParams` below.                               |

### LogListParams

| Field      | Type       | Description                                                                                                                                                     |
| ---------- | ---------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Unit`     | `string`   | Unit is the unit of the query: only the entries of this systemd unit, such as nginx.service. The zero value leaves it out.                                      |
| `Lines`    | `int`      | Lines is the lines of the query: how many of the newest entries; 100 when it is not set. The zero value leaves it out.                                          |
| `Since`    | `string`   | Since is the since of the query: only the entries since this time, in any form that journalctl --since takes, such as 1 hour ago. The zero value leaves it out. |
| `Priority` | `Priority` | Priority is the priority of the query: only the entries at this priority or above. The zero value leaves it out.                                                |

## Result

Each host answers with one `osapi.LogEntry` row. A list answers with one row for
each item.

```go
var row osapi.LogEntry
err := report.Decode("list-log", &row)
```

| Field       | Type              | Description                                                                                                                                                                                                              |
| ----------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Error`     | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname`  | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Message`   | `string`          | The message.                                                                                                                                                                                                             |
| `PID`       | `int`             | The process that wrote it.                                                                                                                                                                                               |
| `Priority`  | `Priority`        | Priority is the property priority.                                                                                                                                                                                       |
| `Status`    | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |
| `Timestamp` | `string`          | When journald recorded it, in RFC 3339.                                                                                                                                                                                  |
| `Unit`      | `string`          | The syslog identifier that wrote it.                                                                                                                                                                                     |

## Permissions

The token must have the `log:read` permission. The step sends
`GET /api/node/{hostname}/log`.

## Example

Refer to
[`examples/operations/log.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/log.go).
