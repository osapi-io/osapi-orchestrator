# RouteUpdate

Updates the route configuration for a network interface on the target host.

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
| `opts`          | `RouteConfigOpts` | Route configuration (refer to [RouteCreate](create.md)).  |

## Result type

```go
var result osapi.RouteMutationResult
err := results.Decode("update-route-1", &result)
```

| Field       | Type     | Description                                                  |
| ----------- | -------- | ------------------------------------------------------------ |
| `Interface` | `string` | Interface name.                                              |
| `Changed`   | `bool`   | `true` if the operation changed the routes.                  |
| `Error`     | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** The operation compares the current routes with the routes you
give. It returns `Changed: true` only if it changed the routes.

## Permissions

The token must have the `network:write` permission.

## Example

Refer to
[`examples/operations/route.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/route.go)
for a complete example.
