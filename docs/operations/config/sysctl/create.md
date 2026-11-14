# SysctlCreate

Creates a new sysctl parameter on the target host.

## Usage

```go
step := o.SysctlCreate("web-01", osapi.SysctlCreateOpts{
    Key:   "net.ipv4.ip_forward",
    Value: "1",
})
```

## Parameters

| Parameter | Type               | Description                                               |
| --------- | ------------------ | --------------------------------------------------------- |
| `target`  | `string`           | Target host: `_any`, `_all`, hostname, or label selector. |
| `opts`    | `SysctlCreateOpts` | Create options (refer to the table below).                |

### SysctlCreateOpts

| Field   | Type     | Required | Description                                                           |
| ------- | -------- | -------- | --------------------------------------------------------------------- |
| `Key`   | `string` | Yes      | The key of the sysctl parameter (for example, `net.ipv4.ip_forward`). |
| `Value` | `string` | Yes      | The value of the parameter.                                           |

## Result type

```go
var result osapi.SysctlMutationResult
err := results.Decode("create-sysctl-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Key`     | `string` | The key of the sysctl parameter.                             |
| `Changed` | `bool`   | `true` if the operation created the parameter.               |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Non-idempotent.** If the parameter already exists, the operation returns an
error. Use [SysctlUpdate](update.md) to change a parameter that exists.

## Permissions

The token must have the `sysctl:write` permission.

## Example

Refer to
[`examples/operations/sysctl.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/sysctl.go)
for a complete example.
