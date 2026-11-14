# SysctlUpdate

Updates a sysctl parameter that exists on the target host.

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
| `key`     | `string`           | The key of the sysctl parameter to update.                |
| `opts`    | `SysctlUpdateOpts` | Update options (refer to the table below).                |

### SysctlUpdateOpts

| Field   | Type     | Required | Description                     |
| ------- | -------- | -------- | ------------------------------- |
| `Value` | `string` | Yes      | The new value of the parameter. |

## Result type

```go
var result osapi.SysctlMutationResult
err := results.Decode("update-sysctl-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Key`     | `string` | The key of the sysctl parameter.                             |
| `Changed` | `bool`   | `true` if the operation changed the parameter.               |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** The operation compares the current value with the value you
give. It returns `Changed: true` only if it changed the parameter. If the value
is already the same, the operation returns `Changed: false`.

## Permissions

The token must have the `sysctl:write` permission.

## Example

Refer to
[`examples/operations/sysctl.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/sysctl.go)
for a complete example.
