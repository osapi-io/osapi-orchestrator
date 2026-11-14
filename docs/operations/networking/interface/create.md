# InterfaceCreate

Creates a network interface configuration on the target host through Netplan.

## Usage

```go
dhcp4 := true
step := o.InterfaceCreate("web-01", "eth1", osapi.InterfaceConfigOpts{
    DHCP4: &dhcp4,
})
```

## Parameters

| Parameter   | Type                  | Description                                               |
| ----------- | --------------------- | --------------------------------------------------------- |
| `target`    | `string`              | Target host: `_any`, `_all`, hostname, or label selector. |
| `ifaceName` | `string`              | The name of the network interface to configure.           |
| `opts`      | `InterfaceConfigOpts` | The configuration options (refer to the table below).     |

### InterfaceConfigOpts

| Field       | Type       | Required | Description                        |
| ----------- | ---------- | -------- | ---------------------------------- |
| `DHCP4`     | `*bool`    | No       | Enable or disable DHCPv4.          |
| `DHCP6`     | `*bool`    | No       | Enable or disable DHCPv6.          |
| `Addresses` | `[]string` | No       | The IP addresses in CIDR notation. |
| `Gateway4`  | `string`   | No       | The IPv4 gateway address.          |
| `Gateway6`  | `string`   | No       | The IPv6 gateway address.          |
| `MTU`       | `*int`     | No       | The maximum transmission unit.     |

## Result type

```go
var result osapi.InterfaceMutationResult
err := results.Decode("create-interface-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | The interface name.                                          |
| `Changed` | `bool`   | `true` if the operation created the configuration.           |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Non-idempotent.** If the interface configuration already exists, the operation
returns an error. Use [InterfaceUpdate](update.md) to change a configuration
that exists.

## Permissions

The token must have the `network:write` permission.

## Example

Refer to
[`examples/operations/interface.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/interface.go)
for a complete example.
