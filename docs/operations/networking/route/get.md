# RouteGet

Retrieves routes for a specific network interface from the target node.

## Usage

```go
step := o.RouteGet("web-01", "eth0")
```

## Parameters

| Parameter       | Type     | Description                                               |
| --------------- | -------- | --------------------------------------------------------- |
| `target`        | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `interfaceName` | `string` | Network interface to retrieve routes for.                 |

## Result Type

```go
var result osapi.RouteGetResult
err := results.Decode("get-route-1", &result)
```

| Field      | Type          | Description                                      |
| ---------- | ------------- | ------------------------------------------------ |
| `Hostname` | `string`      | The node's hostname.                             |
| `Routes`   | `[]RouteInfo` | List of routes for the interface.                |
| `Error`    | `string`      | Error message if query failed; empty on success. |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `network:read` permission.

## Example

See
[`examples/operations/route.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/route.go)
for a complete working example.
