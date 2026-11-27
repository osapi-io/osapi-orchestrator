# HostnameGet

Get the hostname of the host.

Reads the host's name as the hostname command prints it, on Linux and on macOS.
It makes no DNS lookup.

## Usage

```go
step := o.HostnameGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |

## Result

Each host answers with one `osapi.HostnameEntry` row. A list answers with one
row for each item.

```go
var row osapi.HostnameEntry
err := report.Decode("get-hostname", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Name`     | `string`          | The hostname of the host.                                                                                                                                                                                                |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `hostname:read` permission. The step sends
`GET /api/node/{hostname}/hostname`.

## Example

Refer to
[`examples/features/guards.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/guards.go).
