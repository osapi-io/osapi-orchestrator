# InterfaceDelete

Deletes a network interface configuration from the target host.

## Usage

```go
step := o.InterfaceDelete("web-01", "eth1")
```

## Parameters

| Parameter   | Type     | Description                                               |
| ----------- | -------- | --------------------------------------------------------- |
| `target`    | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `ifaceName` | `string` | The name of the network interface to delete.              |

## Result type

```go
var result osapi.InterfaceMutationResult
err := results.Decode("delete-interface-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | The interface name.                                          |
| `Changed` | `bool`   | `true` if the operation removed the configuration.           |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** If the configuration does not exist, the operation returns
`Changed: false`. If the configuration exists, the operation removes it and
returns `Changed: true`.

## Permissions

The token must have the `network:write` permission.

## Example

Refer to
[`examples/operations/interface.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/interface.go)
for a complete example.
