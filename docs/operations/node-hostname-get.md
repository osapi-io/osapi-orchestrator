# NodeHostnameGet

Retrieves the hostname and labels of the target node.

## Usage

```go
step := o.NodeHostnameGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result Type

```go
var result orchestrator.HostnameResult
err := results.Decode("node.hostname.get-1", &result)
```

| Field      | Type                | Description                             |
| ---------- | ------------------- | --------------------------------------- |
| `Hostname` | `string`            | The node's hostname.                    |
| `Labels`   | `map[string]string` | Key-value labels assigned to the agent. |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `node:read` permission.

## Example

See
[`examples/basic.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/basic.go)
for a complete working example.
