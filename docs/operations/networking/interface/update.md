# InterfaceUpdate

Updates a network interface configuration on the target node.

## Usage

```go
mtu := 9000
step := o.InterfaceUpdate("web-01", "eth0", osapi.InterfaceConfigOpts{
    MTU: &mtu,
})
```

## Parameters

| Parameter   | Type                  | Description                                               |
| ----------- | --------------------- | --------------------------------------------------------- |
| `target`    | `string`              | Target host: `_any`, `_all`, hostname, or label selector. |
| `ifaceName` | `string`              | Name of the network interface to update.                  |
| `opts`      | `InterfaceConfigOpts` | Configuration options (see [InterfaceCreate](create.md)). |

## Result Type

```go
var result osapi.InterfaceMutationResult
err := results.Decode("update-interface-1", &result)
```

| Field     | Type     | Description                                       |
| --------- | -------- | ------------------------------------------------- |
| `Name`    | `string` | Interface name.                                   |
| `Changed` | `bool`   | Whether the configuration was actually modified.  |
| `Error`   | `string` | Error message if update failed; empty on success. |

## Idempotency

**Idempotent.** Compares the current configuration against the desired state.
Returns `Changed: true` only if the configuration was actually modified.

## Permissions

Requires `network:write` permission.

## Example

See
[`examples/operations/interface.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/interface.go)
for a complete working example.
