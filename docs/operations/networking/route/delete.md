# RouteDelete

Deletes the route configuration for a network interface on the target host.

## Usage

```go
step := o.RouteDelete("web-01", "eth0")
```

## Parameters

| Parameter       | Type     | Description                                               |
| --------------- | -------- | --------------------------------------------------------- |
| `target`        | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `interfaceName` | `string` | The network interface to delete routes from.              |

## Result type

```go
var result osapi.RouteMutationResult
err := results.Decode("delete-route-1", &result)
```

| Field       | Type     | Description                                                  |
| ----------- | -------- | ------------------------------------------------------------ |
| `Interface` | `string` | The interface name.                                          |
| `Changed`   | `bool`   | `true` if the operation removed the routes.                  |
| `Error`     | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** If no routes exist for the interface, the operation returns
`Changed: false`. If routes exist, the operation removes them and returns
`Changed: true`.

## Permissions

The token must have the `network:write` permission.

## Example

Refer to
[`examples/operations/route.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/route.go)
for a complete example.
