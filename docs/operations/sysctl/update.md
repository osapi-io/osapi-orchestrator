# SysctlUpdate

Set a kernel parameter, manage it, and apply it.

The same as create, addressed by the parameter: osapi manages the parameter from
then on whether it did before or not. Setting the value it already has reports
changed: false.

## Usage

```go
step := o.SysctlUpdate("web-01", input)
```

## Parameters

| Parameter | Type                      | Description                                                   |
| --------- | ------------------------- | ------------------------------------------------------------- |
| `target`  | `string`                  | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.SysctlUpdateInput` | Refer to `SysctlUpdateInput` below.                           |

### SysctlUpdateInput

| Field   | Type     | Description                                                                                      |
| ------- | -------- | ------------------------------------------------------------------------------------------------ |
| `Key`   | `string` | The key of the path: parameter name, such as net.ipv4.ip_forward. It names the item in the path. |
| `Value` | `string` | Value to set.                                                                                    |

## Result

Each host answers with one `osapi.SysctlMutationResult` row. A list answers with
one row for each item.

```go
var row osapi.SysctlMutationResult
err := report.Decode("update-sysctl", &row)
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
`PUT /api/node/{hostname}/sysctl/{key}`.

## Example

Refer to
[`examples/operations/sysctl.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/sysctl.go).
