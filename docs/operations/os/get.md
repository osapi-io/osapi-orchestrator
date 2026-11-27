# OSGet

Get the operating system of the host.

Reads the distribution and its version: the ID and VERSION_ID of /etc/os-release
on Linux, and darwin with the ProductVersion of sw_vers on macOS, such as darwin
15.5.

## Usage

```go
step := o.OSGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |

## Result

Each host answers with one `osapi.OsEntry` row. A list answers with one row for
each item.

```go
var row osapi.OsEntry
err := report.Decode("get-os", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `OSInfo`   | `Release`         | OSInfo is the property os_info.                                                                                                                                                                                          |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `os:read` permission. The step sends
`GET /api/node/{hostname}/os`.

## Example

Refer to
[`examples/operations/node-info.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/node-info.go).
