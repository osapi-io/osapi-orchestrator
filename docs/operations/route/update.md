# RouteUpdate

Updates route configuration for a network interface on the target node.

## Usage

```go
step := o.RouteUpdate("web-01", "eth0", osapi.RouteConfigOpts{
    Routes: []osapi.RouteItem{
        {To: "10.0.0.0/8", Via: "192.168.1.254"},
    },
})
```

## Parameters

| Parameter       | Type              | Description                                               |
| --------------- | ----------------- | --------------------------------------------------------- |
| `target`        | `string`          | Target host: `_any`, `_all`, hostname, or label selector. |
| `interfaceName` | `string`          | Network interface to update routes for.                   |
| `opts`          | `RouteConfigOpts` | Route configuration (see [RouteCreate](create.md)).       |

## Result Type

```go
var result osapi.RouteMutationResult
err := results.Decode("update-route-1", &result)
```

| Field       | Type     | Description                                       |
| ----------- | -------- | ------------------------------------------------- |
| `Interface` | `string` | Interface name.                                   |
| `Changed`   | `bool`   | Whether the routes were actually modified.        |
| `Error`     | `string` | Error message if update failed; empty on success. |

## Idempotency

**Idempotent.** Compares the current routes against the desired configuration.
Returns `Changed: true` only if routes were actually modified.

## Permissions

Requires `network:write` permission.

## Example

See
[`examples/operations/route.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/route.go)
for a complete working example.
