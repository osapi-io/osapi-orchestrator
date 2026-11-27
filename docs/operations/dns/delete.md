# DNSDelete

Drop the DNS resolvers osapi set on an interface.

Removes the interface's netplan file and applies what is left. An interface
osapi set nothing on does not change.

## Usage

```go
step := o.DNSDelete("web-01", interfaceName)
```

## Parameters

| Parameter       | Type     | Description                                                   |
| --------------- | -------- | ------------------------------------------------------------- |
| `target`        | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `interfaceName` | `string` | The name of the network interface.                            |

## Result

Each host answers with one `osapi.DnsResult` row. A list answers with one row
for each item.

```go
var row osapi.DnsResult
err := report.Decode("delete-dns", &row)
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
`DELETE /api/node/{hostname}/dns/{interface_name}`.

## Example

Refer to
[`examples/operations/dns-update.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/dns-update.go).
