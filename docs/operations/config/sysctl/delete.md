# SysctlDelete

Deletes a sysctl parameter from the target host.

## Usage

```go
step := o.SysctlDelete("web-01", "net.ipv4.ip_forward")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `key`     | `string` | The key of the sysctl parameter to delete.                |

## Result type

```go
var result osapi.SysctlMutationResult
err := results.Decode("delete-sysctl-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Key`     | `string` | The key of the sysctl parameter.                             |
| `Changed` | `bool`   | `true` if the operation removed the parameter.               |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** If no parameter exists, the operation returns `Changed: false`.
If the parameter exists, the operation removes it and returns `Changed: true`.

## Permissions

The token must have the `sysctl:write` permission.

## Example

Refer to
[`examples/operations/sysctl.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/sysctl.go)
for a complete example.
