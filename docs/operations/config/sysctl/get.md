# SysctlGet

Gets one sysctl parameter by its key from the target host.

## Usage

```go
step := o.SysctlGet("web-01", "net.ipv4.ip_forward")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `key`     | `string` | The key of the sysctl parameter to get.                   |

## Result type

```go
var result osapi.SysctlEntryResult
err := results.Decode("get-sysctl-1", &result)
```

| Field      | Type     | Description                                                  |
| ---------- | -------- | ------------------------------------------------------------ |
| `Hostname` | `string` | The hostname of the agent that ran the operation.            |
| `Key`      | `string` | The key of the sysctl parameter.                             |
| `Value`    | `string` | The current value of the parameter.                          |
| `Error`    | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `sysctl:read` permission.

## Example

Refer to
[`examples/operations/sysctl.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/sysctl.go)
for a complete example.
