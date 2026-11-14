# NodeHostnameUpdate

Sets the system hostname on the target host.

## Usage

```go
step := o.NodeHostnameUpdate("web-01", "new-hostname")
```

## Parameters

| Parameter  | Type     | Description                                               |
| ---------- | -------- | --------------------------------------------------------- |
| `target`   | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `hostname` | `string` | The new hostname.                                         |

## Result type

```go
var result osapi.HostnameUpdateResult
err := results.Decode("update-hostname-1", &result)
```

| Field      | Type     | Description                                                  |
| ---------- | -------- | ------------------------------------------------------------ |
| `Hostname` | `string` | The system hostname.                                         |
| `Changed`  | `bool`   | `true` if the operation changed the hostname.                |
| `Error`    | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** The operation compares the current hostname with the hostname
you give. It returns `Changed: true` only if it changed the hostname. If the
hostname is already the same, the operation returns `Changed: false`.

## Permissions

The token must have the `node:write` permission.

## Example

Refer to
[`examples/operations/hostname-update.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/hostname-update.go)
for a complete example.
