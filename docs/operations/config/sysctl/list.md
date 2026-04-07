# SysctlList

Lists all managed sysctl parameters on the target node.

## Usage

```go
step := o.SysctlList("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result Type

```go
var result osapi.SysctlEntryResult
err := results.Decode("list-sysctl-1", &result)
```

| Field      | Type     | Description                                      |
| ---------- | -------- | ------------------------------------------------ |
| `Hostname` | `string` | The node's hostname.                             |
| `Key`      | `string` | Sysctl parameter key.                            |
| `Value`    | `string` | Current value of the parameter.                  |
| `Error`    | `string` | Error message if query failed; empty on success. |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `sysctl:read` permission.

## Example

See
[`examples/operations/sysctl.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/sysctl.go)
for a complete working example.
