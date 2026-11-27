# SysctlGet

Get a kernel parameter's current value, and whether osapi manages it.

## Usage

```go
step := o.SysctlGet("web-01", key)
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `key`     | `string` | The name of the kernel parameter.                             |

## Result

Each host answers with one `osapi.SysctlEntry` row. A list answers with one row
for each item.

```go
var row osapi.SysctlEntry
err := report.Decode("get-sysctl", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Key`      | `string`          | Parameter name.                                                                                                                                                                                                          |
| `Managed`  | `bool`            | Whether osapi manages it: a file of its own sets it at boot.                                                                                                                                                             |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |
| `Value`    | `string`          | Current runtime value.                                                                                                                                                                                                   |

## Permissions

The token must have the `sysctl:read` permission. The step sends
`GET /api/node/{hostname}/sysctl/{key}`.

## Example

Refer to
[`examples/operations/sysctl.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/sysctl.go).
