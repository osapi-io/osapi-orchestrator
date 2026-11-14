# NodeStatusGet

Gets the current status of the target host. The result has the registration
details and the agent metadata.

## Usage

```go
step := o.NodeStatusGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result type

This operation has no typed result struct yet. The raw result data is in the
`Data` field of the SDK result, as `map[string]any`.

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `node:read` permission.

## Example

Refer to
[`examples/operations/node-info.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/node-info.go)
for a complete example.
