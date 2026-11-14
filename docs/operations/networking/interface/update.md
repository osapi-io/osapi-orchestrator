# InterfaceUpdate

Updates a network interface configuration on the target host.

## Usage

```go
mtu := 9000
step := o.InterfaceUpdate("web-01", "eth0", osapi.InterfaceConfigOpts{
    MTU: &mtu,
})
```

## Parameters

| Parameter   | Type                  | Description                                                    |
| ----------- | --------------------- | -------------------------------------------------------------- |
| `target`    | `string`              | Target host: `_any`, `_all`, hostname, or label selector.      |
| `ifaceName` | `string`              | Name of the network interface to update.                       |
| `opts`      | `InterfaceConfigOpts` | Configuration options (refer to [InterfaceCreate](create.md)). |

## Result type

```go
var result osapi.InterfaceMutationResult
err := results.Decode("update-interface-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | Interface name.                                              |
| `Changed` | `bool`   | `true` if the operation changed the configuration.           |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** The operation compares the current configuration with the
configuration you give. It returns `Changed: true` only if it changed the
configuration.

## Permissions

The token must have the `network:write` permission.

## Example

Refer to
[`examples/operations/interface.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/interface.go)
for a complete example.
