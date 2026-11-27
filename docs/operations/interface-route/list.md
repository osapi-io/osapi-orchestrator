# InterfaceRouteList

List the routes of an interface.

Reads the routes of the main table that go out of the interface: ip route on
Linux, netstat on macOS. proto says where a route came from, such as dhcp,
kernel or static. An interface with no routes has none; one the host does not
have is not found.

## Usage

```go
step := o.InterfaceRouteList("web-01", name)
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `name`    | `string` | The name of the item.                                         |

## Result

Each host answers with one `osapi.InterfaceRouteEntry` row. A list answers with
one row for each item.

```go
var row osapi.InterfaceRouteEntry
err := report.Decode("list-interface-route", &row)
```

| Field         | Type              | Description                                                                                                                                                                                                              |
| ------------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Destination` | `string`          | The destination in CIDR notation, or default.                                                                                                                                                                            |
| `Error`       | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Family`      | `string`          | inet or inet6.                                                                                                                                                                                                           |
| `Gateway`     | `string`          | The gateway; empty for a route on the link.                                                                                                                                                                              |
| `Hostname`    | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Interface`   | `string`          | The interface the route goes out of.                                                                                                                                                                                     |
| `Metric`      | `int`             | The metric; 0 when the kernel gives none.                                                                                                                                                                                |
| `Proto`       | `string`          | Where the route came from, such as kernel, dhcp or static; empty on macOS.                                                                                                                                               |
| `Scope`       | `string`          | The scope, such as link; empty for a global route.                                                                                                                                                                       |
| `Status`      | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `interface:read` permission. The step sends
`GET /api/node/{hostname}/interface/{name}/route`.

## Example

Refer to
[`examples/operations/route.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/route.go).
