# SysctlList

List the kernel's parameters and which osapi manages.

Lists every parameter the kernel shows under /proc/sys, with its current value,
as sysctl -a does, and marks the ones osapi manages. A prefix, such as net.ipv4,
limits the list to that part of the tree. The kernel can lack a parameter that
osapi manages, for example when its module is not loaded. Then the list shows
the value that osapi wrote.

## Usage

```go
step := o.SysctlList("web-01", params)
```

## Parameters

| Parameter | Type                     | Description                                                   |
| --------- | ------------------------ | ------------------------------------------------------------- |
| `target`  | `string`                 | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `params`  | `osapi.SysctlListParams` | Refer to `SysctlListParams` below.                            |

### SysctlListParams

| Field    | Type     | Description                                                                                                                   |
| -------- | -------- | ----------------------------------------------------------------------------------------------------------------------------- |
| `Prefix` | `string` | Prefix is the prefix of the query: only this parameter and the ones under it, such as net.ipv4. The zero value leaves it out. |

## Result

Each host answers with one `osapi.SysctlEntry` row. A list answers with one row
for each item.

```go
var row osapi.SysctlEntry
err := report.Decode("list-sysctl", &row)
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
`GET /api/node/{hostname}/sysctl`.

## Example

Refer to
[`examples/operations/sysctl.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/sysctl.go).
