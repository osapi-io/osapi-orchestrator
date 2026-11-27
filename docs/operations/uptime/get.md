# UptimeGet

Get the uptime of the host.

Reads how long the host has run since it booted: /proc/uptime on Linux, and the
boot time that sysctl gives on macOS. The uptime is in days, hours and minutes,
and in seconds.

## Usage

```go
step := o.UptimeGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |

## Result

Each host answers with one `osapi.UptimeEntry` row. A list answers with one row
for each item.

```go
var row osapi.UptimeEntry
err := report.Decode("get-uptime", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Seconds`  | `int`             | How long the host has run, in seconds.                                                                                                                                                                                   |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |
| `Uptime`   | `string`          | How long the host has run, in days, hours and minutes.                                                                                                                                                                   |

## Permissions

The token must have the `uptime:read` permission. The step sends
`GET /api/node/{hostname}/uptime`.

## Example

Refer to
[`examples/features/parallel.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/parallel.go).
