# RouteCreate

Creates route configuration for a network interface on the target node.

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
| `opts`          | `RouteConfigOpts` | Route configuration (see below).                          |

### RouteConfigOpts

| Field    | Type          | Required | Description            |
| -------- | ------------- | -------- | ---------------------- |
| `Routes` | `[]RouteItem` | Yes      | List of route entries. |

### RouteItem

| Field | Type     | Required | Description                   |
| ----- | -------- | -------- | ----------------------------- |
| `To`  | `string` | Yes      | Destination in CIDR notation. |
| `Via` | `string` | Yes      | Gateway IP address.           |

## Result Type

```go
var result osapi.RouteMutationResult
err := results.Decode("create-route-1", &result)
```

| Field       | Type     | Description                                         |
| ----------- | -------- | --------------------------------------------------- |
| `Interface` | `string` | Interface name.                                     |
| `Changed`   | `bool`   | Whether the routes were created.                    |
| `Error`     | `string` | Error message if creation failed; empty on success. |

## Idempotency

**Non-idempotent.** Creating routes on an interface that already has routes
returns an error. Use [RouteUpdate](update.md) to modify existing routes.

## Permissions

Requires `network:write` permission.

## Example

See
[`examples/operations/route.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/route.go)
for a complete working example.
