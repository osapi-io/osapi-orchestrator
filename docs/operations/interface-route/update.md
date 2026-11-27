# InterfaceRouteUpdate

Change the gateway or the metric of a route.

Writes the route over the one osapi has to the same destination, and applies it.
A destination osapi has no route to is not found: create adds it.

## Usage

```go
step := o.InterfaceRouteUpdate("web-01", input)
```

## Parameters

| Parameter | Type                              | Description                                                   |
| --------- | --------------------------------- | ------------------------------------------------------------- |
| `target`  | `string`                          | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.InterfaceRouteUpdateInput` | Refer to `InterfaceRouteUpdateInput` below.                   |

### InterfaceRouteUpdateInput

| Field         | Type     | Description                                                                                                                  |
| ------------- | -------- | ---------------------------------------------------------------------------------------------------------------------------- |
| `Name`        | `string` | The name of the path: the interface, such as eth0, or a fact such as @fact.interface.primary. It names the item in the path. |
| `Destination` | `string` | Where the route goes, in CIDR notation, such as 10.0.0.0/8; not a default route.                                             |
| `Metric`      | `int`    | The metric, from 0 to 4294967295; none when it is 0.                                                                         |
| `Via`         | `string` | The gateway, an address of the destination's family, such as 192.168.1.1.                                                    |

## Result

Each host answers with one `osapi.InterfaceRouteResult` row. A list answers with
one row for each item.

```go
var row osapi.InterfaceRouteResult
err := report.Decode("update-interface-route", &row)
```

| Field         | Type              | Description                                                                                                                                                                                                              |
| ------------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`     | `bool`            | Whether the host changed.                                                                                                                                                                                                |
| `Destination` | `string`          | The destination of the route.                                                                                                                                                                                            |
| `Error`       | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname`    | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Status`      | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `interface:write` permission. The step sends
`PUT /api/node/{hostname}/interface/{name}/route`.

## Example

Refer to
[`examples/operations/route.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/route.go).
