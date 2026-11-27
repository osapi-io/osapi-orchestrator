# InterfaceGet

Get a network interface of the host.

Reads the interface as list does.

## Usage

```go
step := o.InterfaceGet("web-01", name)
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `name`    | `string` | The name of the item.                                         |

## Result

Each host answers with one `osapi.InterfaceEntry` row. A list answers with one
row for each item.

```go
var row osapi.InterfaceEntry
err := report.Decode("get-interface", &row)
```

| Field       | Type              | Description                                                                                                                                                                                                              |
| ----------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Addresses` | `[]string`        | Each of its addresses, in CIDR notation.                                                                                                                                                                                 |
| `DHCP4`     | `bool`            | Whether DHCP gave it an IPv4 address: a route of it has proto dhcp.                                                                                                                                                      |
| `Error`     | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname`  | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `IPv4`      | `string`          | Its first IPv4 address that is not link-local.                                                                                                                                                                           |
| `IPv6`      | `string`          | Its first IPv6 address that is not link-local.                                                                                                                                                                           |
| `MAC`       | `string`          | Its MAC address; empty for one without, such as lo.                                                                                                                                                                      |
| `MTU`       | `int`             | Its MTU.                                                                                                                                                                                                                 |
| `Name`      | `string`          | The interface.                                                                                                                                                                                                           |
| `Primary`   | `bool`            | Whether it has the default route.                                                                                                                                                                                        |
| `State`     | `string`          | up or down.                                                                                                                                                                                                              |
| `Status`    | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `interface:read` permission. The step sends
`GET /api/node/{hostname}/interface/{name}`.

## Example

Refer to
[`examples/operations/interface.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/interface.go).
