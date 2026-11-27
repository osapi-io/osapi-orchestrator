# SysctlCreate

Manage a kernel parameter and apply it.

Writes osapi-net.ipv4.ip_forward.conf for net.ipv4.ip_forward to
behaviors.sysctl.dir, /etc/sysctl.d by default, and applies it with sysctl -p.
It fails when sysctl -p fails, because the value is not set. A second run with
the same value changes nothing and reports changed: false. If the kernel has a
different value, for example after a failed apply, it applies the value again.

## Usage

```go
step := o.SysctlCreate("web-01", input)
```

## Parameters

| Parameter | Type                      | Description                                                   |
| --------- | ------------------------- | ------------------------------------------------------------- |
| `target`  | `string`                  | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.SysctlCreateInput` | Refer to `SysctlCreateInput` below.                           |

### SysctlCreateInput

| Field   | Type     | Description                                  |
| ------- | -------- | -------------------------------------------- |
| `Key`   | `string` | Parameter name, such as net.ipv4.ip_forward. |
| `Value` | `string` | Value to set.                                |

## Result

Each host answers with one `osapi.SysctlMutationResult` row. A list answers with
one row for each item.

```go
var row osapi.SysctlMutationResult
err := report.Decode("create-sysctl", &row)
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
`POST /api/node/{hostname}/sysctl`.

## Example

Refer to
[`examples/operations/sysctl.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/sysctl.go).
