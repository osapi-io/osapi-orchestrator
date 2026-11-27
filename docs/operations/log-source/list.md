# LogSourceList

List the sources of the journal.

Runs journalctl --field=SYSLOG_IDENTIFIER: each syslog identifier that has an
entry, once, sorted. A name here is what log list takes as unit when the source
is a systemd unit.

## Usage

```go
step := o.LogSourceList("web-01")
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |

## Result

Each host answers with one `osapi.LogSourceSource` row. A list answers with one
row for each item.

```go
var row osapi.LogSourceSource
err := report.Decode("list-log-source", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Name`     | `string`          | The syslog identifier, such as sshd.                                                                                                                                                                                     |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `log:read` permission. The step sends
`GET /api/node/{hostname}/log/source`.

## Example

Refer to
[`examples/operations/log.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/log.go).
