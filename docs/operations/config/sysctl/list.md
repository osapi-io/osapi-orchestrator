# SysctlList

Lists all the managed sysctl parameters on the target host.

## Usage

```go
step := o.SysctlList("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result type

```go
var result osapi.SysctlEntryResult
err := results.Decode("list-sysctl-1", &result)
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
