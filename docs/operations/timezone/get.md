# TimezoneGet

Get the timezone of the host.

Reads the IANA name of the timezone and its offset from UTC now: the zone that
/etc/localtime links to, or /etc/timezone on Linux.

## Usage

```go
step := o.TimezoneGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |

## Result

Each host answers with one `osapi.TimezoneEntry` row. A list answers with one
row for each item.

```go
var row osapi.TimezoneEntry
err := report.Decode("get-timezone", &row)
```

| Field       | Type              | Description                                                                                                                                                                                                              |
| ----------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Error`     | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname`  | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Status`    | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |
| `Timezone`  | `string`          | The IANA name of the timezone, such as America/New_York.                                                                                                                                                                 |
| `UTCOffset` | `string`          | The offset from UTC now, such as -05:00.                                                                                                                                                                                 |

## Permissions

The token must have the `timezone:read` permission. The step sends
`GET /api/node/{hostname}/timezone`.

## Example

Refer to
[`examples/operations/timezone.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/timezone.go).
