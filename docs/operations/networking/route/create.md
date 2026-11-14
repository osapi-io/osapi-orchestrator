# RouteCreate

Creates the route configuration for a network interface on the target host.

## Usage

```go
step := o.RouteCreate("web-01", "eth0", osapi.RouteConfigOpts{
    Routes: []osapi.RouteItem{
        {To: "10.0.0.0/8", Via: "192.168.1.1"},
    },
})
```

## Parameters

| Parameter       | Type              | Description                                               |
| --------------- | ----------------- | --------------------------------------------------------- |
| `target`        | `string`          | Target host: `_any`, `_all`, hostname, or label selector. |
| `interfaceName` | `string`          | Network interface to add routes to.                       |
| `opts`          | `RouteConfigOpts` | Route configuration (refer to the table below).           |

### RouteConfigOpts

| Field    | Type          | Required | Description            |
| -------- | ------------- | -------- | ---------------------- |
| `Routes` | `[]RouteItem` | Yes      | List of route entries. |

### RouteItem

| Field | Type     | Required | Description                   |
| ----- | -------- | -------- | ----------------------------- |
| `To`  | `string` | Yes      | Destination in CIDR notation. |
| `Via` | `string` | Yes      | Gateway IP address.           |

## Result type

```go
var result osapi.RouteMutationResult
err := results.Decode("create-route-1", &result)
```

| Field       | Type     | Description                                                  |
| ----------- | -------- | ------------------------------------------------------------ |
| `Interface` | `string` | Interface name.                                              |
| `Changed`   | `bool`   | `true` if the operation created the routes.                  |
| `Error`     | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Non-idempotent.** If the interface already has routes, the operation returns
an error. Use [RouteUpdate](update.md) to change routes that exist.

## Permissions

The token must have the `network:write` permission.

## Example

Refer to
[`examples/operations/route.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/route.go)
for a complete example.
