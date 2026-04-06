# SysctlDelete

Deletes a sysctl parameter from the target node.

## Usage

```go
step := o.SysctlDelete("web-01", "net.ipv4.ip_forward")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `key`     | `string` | Sysctl parameter key to delete.                           |

## Result Type

```go
var result osapi.SysctlMutationResult
err := results.Decode("delete-sysctl-1", &result)
```

| Field     | Type     | Description                                          |
| --------- | -------- | ---------------------------------------------------- |
| `Key`     | `string` | Sysctl parameter key.                                |
| `Changed` | `bool`   | Whether the parameter was removed or already absent. |
| `Error`   | `string` | Error message if deletion failed; empty on success.  |

## Idempotency

**Idempotent.** If the parameter does not exist, the operation returns
`Changed: false`. If the parameter exists and is removed, returns
`Changed: true`.

## Permissions

Requires `sysctl:write` permission.

## Example

See
[`examples/operations/sysctl.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/sysctl.go)
for a complete working example.
