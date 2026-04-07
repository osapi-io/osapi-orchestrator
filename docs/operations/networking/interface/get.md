# InterfaceGet

Retrieves a specific network interface configuration from the target node.

## Usage

```go
step := o.InterfaceGet("web-01", "eth0")
```

## Parameters

| Parameter   | Type     | Description                                               |
| ----------- | -------- | --------------------------------------------------------- |
| `target`    | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `ifaceName` | `string` | Name of the network interface.                            |

## Result Type

```go
var result osapi.InterfaceGetResult
err := results.Decode("get-interface-1", &result)
```

| Field       | Type             | Description                                      |
| ----------- | ---------------- | ------------------------------------------------ |
| `Hostname`  | `string`         | The node's hostname.                             |
| `Interface` | `*InterfaceInfo` | Interface details.                               |
| `Error`     | `string`         | Error message if query failed; empty on success. |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `network:read` permission.

## Example

See
[`examples/operations/interface.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/interface.go)
for a complete working example.
