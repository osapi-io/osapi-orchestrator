# InterfaceGet

Gets the configuration of a specific network interface from the target host.

## Usage

```go
step := o.InterfaceGet("web-01", "eth0")
```

## Parameters

| Parameter   | Type     | Description                                               |
| ----------- | -------- | --------------------------------------------------------- |
| `target`    | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `ifaceName` | `string` | Name of the network interface.                            |

## Result type

```go
var result osapi.InterfaceGetResult
err := results.Decode("get-interface-1", &result)
```

| Field       | Type             | Description                                                  |
| ----------- | ---------------- | ------------------------------------------------------------ |
| `Hostname`  | `string`         | The hostname of the agent that ran the operation.            |
| `Interface` | `*InterfaceInfo` | Interface details.                                           |
| `Error`     | `string`         | The error message if the operation failed. Empty on success. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `network:read` permission.

## Example

Refer to
[`examples/operations/interface.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/interface.go)
for a complete example.
