# RouteList

Lists network routes on the target host.

## Usage

```go
step := o.RouteList("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result type

```go
var result osapi.RouteListResult
err := results.Decode("list-route-1", &result)
```

| Field      | Type          | Description                                                  |
| ---------- | ------------- | ------------------------------------------------------------ |
| `Hostname` | `string`      | The hostname of the agent that ran the operation.            |
| `Routes`   | `[]RouteInfo` | List of route entries.                                       |
| `Error`    | `string`      | The error message if the operation failed. Empty on success. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `network:read` permission.

## Example

Refer to
[`examples/operations/route.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/route.go)
for a complete example.
