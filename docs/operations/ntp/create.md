# NTPCreate

Set the NTP servers that chrony uses.

Writes a server line for each server to /etc/chrony/sources.d/osapi-ntp.sources,
and reloads the sources with chronyc reload sources. A file that has the same
servers does not change, and reports changed: false.

## Usage

```go
step := o.NTPCreate("web-01", input)
```

## Parameters

| Parameter | Type                   | Description                                                   |
| --------- | ---------------------- | ------------------------------------------------------------- |
| `target`  | `string`               | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.NTPCreateInput` | Refer to `NTPCreateInput` below.                              |

### NTPCreateInput

| Field     | Type       | Description                                 |
| --------- | ---------- | ------------------------------------------- |
| `Servers` | `[]string` | The NTP servers, hostnames or IP addresses. |

## Result

Each host answers with one `osapi.NtpResult` row. A list answers with one row
for each item.

```go
var row osapi.NtpResult
err := report.Decode("create-ntp", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`  | `bool`            | Whether the host changed.                                                                                                                                                                                                |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `ntp:write` permission. The step sends
`POST /api/node/{hostname}/ntp`.

## Example

Refer to
[`examples/operations/ntp.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/ntp.go).
