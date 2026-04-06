# InterfaceDelete

Deletes a network interface configuration from the target node.

## Usage

```go
step := o.InterfaceDelete("web-01", "eth1")
```

## Parameters

| Parameter   | Type     | Description                                               |
| ----------- | -------- | --------------------------------------------------------- |
| `target`    | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `ifaceName` | `string` | Name of the network interface to delete.                  |

## Result Type

```go
var result osapi.InterfaceMutationResult
err := results.Decode("delete-interface-1", &result)
```

| Field     | Type     | Description                                              |
| --------- | -------- | -------------------------------------------------------- |
| `Name`    | `string` | Interface name.                                          |
| `Changed` | `bool`   | Whether the configuration was removed or already absent. |
| `Error`   | `string` | Error message if deletion failed; empty on success.      |

## Idempotency

**Idempotent.** If the configuration does not exist, returns `Changed: false`.

## Permissions

Requires `network:write` permission.

## Example

See
[`examples/operations/interface.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/interface.go)
for a complete working example.
