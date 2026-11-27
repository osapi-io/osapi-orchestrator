# SysctlDelete

Stop the management of a kernel parameter.

Removes the file that osapi wrote for it and reloads with sysctl --system. Then
a setting for the key in a different file takes effect. If there is no such
setting, the parameter keeps its current value until the next reboot. A delete
of a parameter that osapi does not manage reports changed: false.

## Usage

```go
step := o.SysctlDelete("web-01", key)
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `key`     | `string` | The name of the kernel parameter.                             |

## Result

Each host answers with one `osapi.SysctlMutationResult` row. A list answers with
one row for each item.

```go
var row osapi.SysctlMutationResult
err := report.Decode("delete-sysctl", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`  | `bool`            | Whether the host changed.                                                                                                                                                                                                |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Key`      | `string`          | Parameter name.                                                                                                                                                                                                          |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `sysctl:write` permission. The step sends
`DELETE /api/node/{hostname}/sysctl/{key}`.

## Example

Refer to
[`examples/operations/sysctl.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/sysctl.go).
