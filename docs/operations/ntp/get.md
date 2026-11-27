# NTPGet

Get the sync status of chrony and its sources.

Reads chronyc tracking, for whether chrony is synchronized, its stratum, its
offset from NTP time and the source that it follows, and chronyc sources, for
each source that it has.

## Usage

```go
step := o.NTPGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |

## Result

Each host answers with one `osapi.NtpStatus` row. A list answers with one row
for each item.

```go
var row osapi.NtpStatus
err := report.Decode("get-ntp", &row)
```

| Field           | Type              | Description                                                                                                                                                                                                              |
| --------------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `CurrentSource` | `string`          | The source that chrony follows now.                                                                                                                                                                                      |
| `Error`         | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname`      | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Offset`        | `string`          | The offset of the system clock from NTP time, such as +0.000003422s.                                                                                                                                                     |
| `Servers`       | `[]string`        | The sources that chrony has.                                                                                                                                                                                             |
| `Status`        | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |
| `Stratum`       | `int`             | The NTP stratum of the host.                                                                                                                                                                                             |
| `Synchronized`  | `bool`            | Whether chrony is synchronized to a source.                                                                                                                                                                              |

## Permissions

The token must have the `ntp:read` permission. The step sends
`GET /api/node/{hostname}/ntp`.

## Example

Refer to
[`examples/operations/ntp.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/ntp.go).
