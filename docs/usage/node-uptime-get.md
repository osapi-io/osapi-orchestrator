# NodeUptimeGet

Retrieves the system uptime of the target node.

## Usage

```go
step := o.NodeUptimeGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result Type

There is no typed result struct for this operation yet. The raw result data is
available through the SDK result's `Data` field as `map[string]any`.

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `node:read` permission.
