# NodeHostnameUpdate

Sets the system hostname on the target node.

## Usage

```go
step := o.NodeHostnameUpdate("web-01", "new-hostname")
```

## Parameters

| Parameter  | Type     | Description                                               |
| ---------- | -------- | --------------------------------------------------------- |
| `target`   | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `hostname` | `string` | The new hostname to set.                                  |

## Result Type

```go
var result osapi.HostnameUpdateResult
err := results.Decode("update-hostname-1", &result)
```

| Field      | Type     | Description                                       |
| ---------- | -------- | ------------------------------------------------- |
| `Hostname` | `string` | The node's hostname.                              |
| `Changed`  | `bool`   | Whether the hostname was actually changed.        |
| `Error`    | `string` | Error message if update failed; empty on success. |

## Idempotency

**Idempotent.** Checks the current hostname before applying the change. Returns
`Changed: true` only if the hostname was actually modified. If the hostname
already matches the desired value, the operation returns `Changed: false`.

## Permissions

Requires `node:write` permission.

## Example

See
[`examples/operations/hostname-update.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/hostname-update.go)
for a complete working example.
