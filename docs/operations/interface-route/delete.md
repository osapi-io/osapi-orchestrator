# InterfaceRouteDelete

Drop a route that osapi wrote.

Drops the route from the interface's netplan file and applies what is left; the
last route removes the file. A route osapi does not have does not change.

## Usage

```go
step := o.InterfaceRouteDelete("web-01", input)
```

## Parameters

| Parameter | Type                              | Description                                                   |
| --------- | --------------------------------- | ------------------------------------------------------------- |
| `target`  | `string`                          | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.InterfaceRouteDeleteInput` | Refer to `InterfaceRouteDeleteInput` below.                   |

### InterfaceRouteDeleteInput

| Field         | Type     | Description                                                                                                                  |
| ------------- | -------- | ---------------------------------------------------------------------------------------------------------------------------- |
| `Name`        | `string` | The name of the path: the interface, such as eth0, or a fact such as @fact.interface.primary. It names the item in the path. |
| `Destination` | `string` | The destination of the route, in CIDR notation, such as 10.0.0.0/8.                                                          |

## Result

Each host answers with one `osapi.InterfaceRouteResult` row. A list answers with
one row for each item.

```go
var row osapi.InterfaceRouteResult
err := report.Decode("delete-interface-route", &row)
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
`DELETE /api/node/{hostname}/interface/{name}/route`.

## Example

Refer to
[`examples/operations/route.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/route.go).
