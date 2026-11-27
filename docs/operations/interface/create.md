# InterfaceCreate

Configure a network interface.

Writes a netplan file of the interface's own, such as
/etc/netplan/osapi-osapi0.yaml for osapi0, under a netplan definition that is
the interface's own, and applies it. type is ethernet, the default, or dummy. An
ethernet is an interface the host has, in the netplan section that netplan
status gives for it: one the host does not have is not found. A dummy goes under
dummy-devices, and netplan apply makes it, so the host need not have it.
dummy-devices needs netplan 0.107 or later. dhcp4 and dhcp6 are false unless
given, so the interface has only the addresses given. An interface that a match
defines, such as cloud-init's all for en\*, is refused, because the config would
change each interface the match takes, now and later. An interface osapi already
configures does not change: update replaces its config. It runs on the Debian
family, and elsewhere, and in a container, reports skipped.

## Usage

```go
step := o.InterfaceCreate("web-01", input)
```

## Parameters

| Parameter | Type                         | Description                                                   |
| --------- | ---------------------------- | ------------------------------------------------------------- |
| `target`  | `string`                     | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.InterfaceCreateInput` | Refer to `InterfaceCreateInput` below.                        |

### InterfaceCreateInput

| Field        | Type                  | Description                                                                                                                  |
| ------------ | --------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| `Name`       | `string`              | The name of the path: the interface, such as eth0, or a fact such as @fact.interface.primary. It names the item in the path. |
| `Addresses`  | `[]string`            | The static addresses, in CIDR notation, such as 10.0.0.5/24.                                                                 |
| `DHCP4`      | `bool`                | Take an IPv4 address from DHCP.                                                                                              |
| `DHCP6`      | `bool`                | Take an IPv6 address from DHCP.                                                                                              |
| `Gateway4`   | `string`              | The IPv4 default gateway.                                                                                                    |
| `Gateway6`   | `string`              | The IPv6 default gateway.                                                                                                    |
| `MACAddress` | `string`              | The MAC address to set.                                                                                                      |
| `MTU`        | `int`                 | The MTU, from 68 to 9000.                                                                                                    |
| `Type`       | `types.InterfaceType` | InterfaceType is the type of a network interface that osapi manages.                                                         |
| `Wakeonlan`  | `bool`                | Turn on wake-on-LAN.                                                                                                         |

## Result

Each host answers with one `osapi.InterfaceResult` row. A list answers with one
row for each item.

```go
var row osapi.InterfaceResult
err := report.Decode("create-interface", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`  | `bool`            | Whether the host changed.                                                                                                                                                                                                |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Name`     | `string`          | The interface.                                                                                                                                                                                                           |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `interface:write` permission. The step sends
`POST /api/node/{hostname}/interface/{name}`.

## Example

Refer to
[`examples/operations/interface.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/interface.go).
