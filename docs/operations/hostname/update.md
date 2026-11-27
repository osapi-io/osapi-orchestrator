# HostnameUpdate

Set the hostname of the host.

Sets the static hostname with hostnamectl set-hostname, which keeps it after a
reboot. A host that already has the name does not change and reports changed:
false. It runs on the Debian family only, and elsewhere reports skipped. In a
container the runtime owns the hostname, so there it also reports skipped.

## Usage

```go
step := o.HostnameUpdate("web-01", input)
```

## Parameters

| Parameter | Type                        | Description                                                   |
| --------- | --------------------------- | ------------------------------------------------------------- |
| `target`  | `string`                    | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.HostnameUpdateInput` | Refer to `HostnameUpdateInput` below.                         |

### HostnameUpdateInput

| Field  | Type     | Description                          |
| ------ | -------- | ------------------------------------ |
| `Name` | `string` | The hostname to set, such as web-01. |

## Result

Each host answers with one `osapi.HostnameResult` row. A list answers with one
row for each item.

```go
var row osapi.HostnameResult
err := report.Decode("update-hostname", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`  | `bool`            | Whether the host changed.                                                                                                                                                                                                |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Name`     | `string`          | The hostname of the host.                                                                                                                                                                                                |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `hostname:write` permission. The step sends
`PUT /api/node/{hostname}/hostname`.

## Example

Refer to
[`examples/features/host-status.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/host-status.go).
