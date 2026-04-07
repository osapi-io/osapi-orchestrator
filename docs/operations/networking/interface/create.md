# InterfaceCreate

Creates a network interface configuration on the target node via Netplan.

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
| `ifaceName` | `string`              | Name of the network interface to configure.               |
| `opts`      | `InterfaceConfigOpts` | Configuration options (see below).                        |

### InterfaceConfigOpts

| Field       | Type       | Required | Description                    |
| ----------- | ---------- | -------- | ------------------------------ |
| `DHCP4`     | `*bool`    | No       | Enable or disable DHCPv4.      |
| `DHCP6`     | `*bool`    | No       | Enable or disable DHCPv6.      |
| `Addresses` | `[]string` | No       | IP addresses in CIDR notation. |
| `Gateway4`  | `string`   | No       | IPv4 gateway address.          |
| `Gateway6`  | `string`   | No       | IPv6 gateway address.          |
| `MTU`       | `*int`     | No       | Maximum transmission unit.     |

## Result Type

```go
var result osapi.InterfaceMutationResult
err := results.Decode("create-interface-1", &result)
```

| Field     | Type     | Description                                         |
| --------- | -------- | --------------------------------------------------- |
| `Name`    | `string` | Interface name.                                     |
| `Changed` | `bool`   | Whether the configuration was created.              |
| `Error`   | `string` | Error message if creation failed; empty on success. |

## Idempotency

**Non-idempotent.** Creating an interface configuration that already exists
returns an error. Use [InterfaceUpdate](update.md) to modify existing configs.

## Permissions

Requires `network:write` permission.

## Example

See
[`examples/operations/interface.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/interface.go)
for a complete working example.
