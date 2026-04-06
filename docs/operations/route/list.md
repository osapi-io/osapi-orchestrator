# RouteList

Lists network routes on the target node.

## Usage

```go
step := o.RouteList("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result Type

```go
var result osapi.RouteListResult
err := results.Decode("list-route-1", &result)
```

| Field      | Type          | Description                                      |
| ---------- | ------------- | ------------------------------------------------ |
| `Hostname` | `string`      | The node's hostname.                             |
| `Routes`   | `[]RouteInfo` | List of route entries.                           |
| `Error`    | `string`      | Error message if query failed; empty on success. |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `network:read` permission.

## Example

See
[`examples/operations/route.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/route.go)
for a complete working example.
