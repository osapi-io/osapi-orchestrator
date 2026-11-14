# NetworkDNSDelete

Deletes the DNS configuration for a specific network interface on the target
host.

## Usage

```go
step := o.NetworkDNSDelete("web-01", "eth0")
```

## Parameters

| Parameter       | Type     | Description                                               |
| --------------- | -------- | --------------------------------------------------------- |
| `target`        | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `interfaceName` | `string` | Network interface to delete DNS configuration from.       |

## Result type

```go
var result osapi.DNSDeleteResult
err := results.Decode("delete-dns-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Changed` | `bool`   | `true` if the operation removed the DNS configuration.       |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** If no DNS configuration exists for the interface, the operation
returns `Changed: false`. If the configuration exists, the operation removes it
and returns `Changed: true`.

## Permissions

The token must have the `network:write` permission.

## Example

Refer to
[`examples/operations/dns-update.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/dns-update.go)
for a complete example.
