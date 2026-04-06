# NetworkDNSDelete

Deletes the DNS configuration for a specific network interface on the target
node.

## Usage

```go
step := o.NetworkDNSDelete("web-01", "eth0")
```

## Parameters

| Parameter       | Type     | Description                                               |
| --------------- | -------- | --------------------------------------------------------- |
| `target`        | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `interfaceName` | `string` | Network interface to delete DNS configuration from.       |

## Result Type

```go
var result osapi.DNSDeleteResult
err := results.Decode("delete-dns-1", &result)
```

| Field     | Type     | Description                                          |
| --------- | -------- | ---------------------------------------------------- |
| `Changed` | `bool`   | Whether the DNS config was removed or already absent. |
| `Error`   | `string` | Error message if deletion failed; empty on success.  |

## Idempotency

**Idempotent.** If no DNS configuration exists for the interface, returns
`Changed: false`. If configuration exists and is removed, returns
`Changed: true`.

## Permissions

Requires `network:write` permission.

## Example

See
[`examples/operations/dns-update.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/dns-update.go)
for a complete working example.
