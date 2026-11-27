# InterfaceUpdate

Change the config osapi has of a network interface.

Writes the interface's netplan file over the one osapi has, as create writes it,
and applies it. Without type dummy, the config goes in the section that netplan
status gives, which is dummy-devices for a dummy, so a dummy stays a dummy. An
interface osapi does not configure is not found: create configures it.

## Usage

```go
step := o.InterfaceUpdate("web-01", input)
```

## Parameters

| Parameter | Type                         | Description                                                   |
| --------- | ---------------------------- | ------------------------------------------------------------- |
| `target`  | `string`                     | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.InterfaceUpdateInput` | Refer to `InterfaceUpdateInput` below.                        |

### InterfaceUpdateInput

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
err := report.Decode("update-interface", &row)
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
`PUT /api/node/{hostname}/interface/{name}`.

## Example

Refer to
[`examples/operations/interface.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/interface.go).
