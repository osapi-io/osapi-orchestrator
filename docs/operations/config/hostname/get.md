# NodeHostnameGet

Gets the hostname and the labels of the target host.

## Usage

```go
step := o.NodeHostnameGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result type

```go
var result osapi.HostnameResult
err := results.Decode("node.hostname.get-1", &result)
```

| Field      | Type                | Description                        |
| ---------- | ------------------- | ---------------------------------- |
| `Hostname` | `string`            | The system hostname.               |
| `Labels`   | `map[string]string` | The key-value labels of the agent. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `node:read` permission.

## Example

Refer to
[`examples/operations/node-info.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/node-info.go)
for a complete example.
