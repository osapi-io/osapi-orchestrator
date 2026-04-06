# SysctlUpdate

Updates an existing sysctl parameter on the target node.

## Usage

```go
step := o.SysctlUpdate("web-01", "net.ipv4.ip_forward", osapi.SysctlUpdateOpts{
    Value: "0",
})
```

## Parameters

| Parameter | Type               | Description                                               |
| --------- | ------------------ | --------------------------------------------------------- |
| `target`  | `string`           | Target host: `_any`, `_all`, hostname, or label selector. |
| `key`     | `string`           | Sysctl parameter key to update.                           |
| `opts`    | `SysctlUpdateOpts` | Update options (see below).                               |

### SysctlUpdateOpts

| Field   | Type     | Required | Description               |
| ------- | -------- | -------- | ------------------------- |
| `Value` | `string` | Yes      | New value for the parameter. |

## Result Type

```go
var result osapi.SysctlMutationResult
err := results.Decode("update-sysctl-1", &result)
```

| Field     | Type     | Description                                       |
| --------- | -------- | ------------------------------------------------- |
| `Key`     | `string` | Sysctl parameter key.                             |
| `Changed` | `bool`   | Whether the parameter was actually modified.      |
| `Error`   | `string` | Error message if update failed; empty on success. |

## Idempotency

**Idempotent.** Compares the current value against the desired value. Returns
`Changed: true` only if the parameter was actually modified. If the value
already matches, returns `Changed: false`.

## Permissions

Requires `sysctl:write` permission.

## Example

See
[`examples/operations/sysctl.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/sysctl.go)
for a complete working example.
