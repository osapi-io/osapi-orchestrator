# DNSGet

Get the DNS resolvers of an interface.

Reads the servers and the search domains of the interface: resolvectl on Linux,
with resolv.conf on a host that has no resolvers per interface, such as a
container, and scutil on macOS.

## Usage

```go
step := o.DNSGet("web-01", interfaceName)
```

## Parameters

| Parameter       | Type     | Description                                                   |
| --------------- | -------- | ------------------------------------------------------------- |
| `target`        | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `interfaceName` | `string` | The name of the network interface.                            |

## Result

Each host answers with one `osapi.DnsEntry` row. A list answers with one row for
each item.

```go
var row osapi.DnsEntry
err := report.Decode("get-dns", &row)
```

| Field           | Type              | Description                                                                                                                                                                                                              |
| --------------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Error`         | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname`      | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `SearchDomains` | `[]string`        | The search domains.                                                                                                                                                                                                      |
| `Servers`       | `[]string`        | The DNS servers.                                                                                                                                                                                                         |
| `Status`        | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `dns:read` permission. The step sends
`GET /api/node/{hostname}/dns/{interface_name}`.

## Example

Refer to
[`examples/operations/dns-update.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/dns-update.go).
