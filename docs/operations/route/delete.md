# RouteDelete

Deletes route configuration for a network interface on the target node.

## Usage

```go
step := o.RouteDelete("web-01", "eth0")
```

## Parameters

| Parameter       | Type     | Description                                               |
| --------------- | -------- | --------------------------------------------------------- |
| `target`        | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `interfaceName` | `string` | Network interface to delete routes from.                  |

## Result Type

```go
var result osapi.RouteMutationResult
err := results.Decode("delete-route-1", &result)
```

| Field       | Type     | Description                                            |
| ----------- | -------- | ------------------------------------------------------ |
| `Interface` | `string` | Interface name.                                        |
| `Changed`   | `bool`   | Whether the routes were removed or already absent.     |
| `Error`     | `string` | Error message if deletion failed; empty on success.    |

## Idempotency

**Idempotent.** If no routes exist for the interface, returns `Changed: false`.

## Permissions

Requires `network:write` permission.

## Example

See
[`examples/operations/route.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/route.go)
for a complete working example.
