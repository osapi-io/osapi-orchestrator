# DNSUpdate

Set the DNS resolvers of an interface.

Writes a netplan file of the interface's own, such as
/etc/netplan/osapi-dns-eth0.yaml for eth0, under the netplan definition that
matches the interface, and applies it. A definition that matches several
interfaces, such as cloud-init's all for en\*, gives each of them the resolvers.
An empty list keeps the interface's servers or domains, and an interface with no
resolvers yet starts from none. override_dhcp drops the servers DHCP gives. It
runs on the Debian family, and elsewhere, and in a container, reports skipped.

## Usage

```go
step := o.DNSUpdate("web-01", input)
```

## Parameters

| Parameter | Type                   | Description                                                   |
| --------- | ---------------------- | ------------------------------------------------------------- |
| `target`  | `string`               | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.DNSUpdateInput` | Refer to `DNSUpdateInput` below.                              |

### DNSUpdateInput

| Field           | Type       | Description                                                                                                                            |
| --------------- | ---------- | -------------------------------------------------------------------------------------------------------------------------------------- |
| `InterfaceName` | `string`   | The interface_name of the path: the interface, such as eth0, or a fact such as @fact.interface.primary. It names the item in the path. |
| `OverrideDHCP`  | `bool`     | Use only these servers: drop the ones DHCP gives.                                                                                      |
| `SearchDomains` | `[]string` | The search domains; the interface's own when empty.                                                                                    |
| `Servers`       | `[]string` | The DNS servers, IPv4 or IPv6; the interface's own when empty.                                                                         |

## Result

Each host answers with one `osapi.DnsResult` row. A list answers with one row
for each item.

```go
var row osapi.DnsResult
err := report.Decode("update-dns", &row)
```

| Field           | Type              | Description                                                                                                                                                                                                              |
| --------------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`       | `bool`            | Whether the host changed.                                                                                                                                                                                                |
| `Error`         | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname`      | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `InterfaceName` | `string`          | The interface.                                                                                                                                                                                                           |
| `Status`        | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `dns:write` permission. The step sends
`PUT /api/node/{hostname}/dns/{interface_name}`.

## Example

Refer to
[`examples/operations/dns-update.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/dns-update.go).
