# SysctlCreate

Creates a new sysctl parameter on the target node.

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
| `opts`    | `SysctlCreateOpts` | Create options (see below).                               |

### SysctlCreateOpts

| Field   | Type     | Required | Description                                       |
| ------- | -------- | -------- | ------------------------------------------------- |
| `Key`   | `string` | Yes      | Sysctl parameter key (e.g., `net.ipv4.ip_forward`). |
| `Value` | `string` | Yes      | Value to set for the parameter.                   |

## Result Type

```go
var result osapi.SysctlMutationResult
err := results.Decode("create-sysctl-1", &result)
```

| Field     | Type     | Description                                         |
| --------- | -------- | --------------------------------------------------- |
| `Key`     | `string` | Sysctl parameter key.                               |
| `Changed` | `bool`   | Whether the parameter was created.                  |
| `Error`   | `string` | Error message if creation failed; empty on success. |

## Idempotency

**Non-idempotent.** Creating a parameter that already exists returns an error.
Use [SysctlUpdate](update.md) to modify existing parameters.

## Permissions

Requires `sysctl:write` permission.

## Example

See
[`examples/operations/sysctl.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/sysctl.go)
for a complete working example.
