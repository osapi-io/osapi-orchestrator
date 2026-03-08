# NetworkDNSUpdate

Updates the DNS configuration for a specific network interface on the target
node. Sets the DNS servers and search domains to the provided values.

## Usage

```go
step := o.NetworkDNSUpdate(
    "web-01",
    "eth0",
    []string{"8.8.8.8", "8.8.4.4"},
    []string{"example.com"},
)
```

## Parameters

| Parameter       | Type       | Description                                               |
| --------------- | ---------- | --------------------------------------------------------- |
| `target`        | `string`   | Target host: `_any`, `_all`, hostname, or label selector. |
| `interfaceName` | `string`   | Network interface to configure (e.g., `eth0`, `ens33`).   |
| `servers`       | `[]string` | DNS server addresses to set.                              |
| `searchDomains` | `[]string` | DNS search domains to set.                                |

## Result Type

```go
var result orchestrator.DNSUpdateResult
err := results.Decode("network.dns.update-1", &result)
```

| Field     | Type     | Description                                       |
| --------- | -------- | ------------------------------------------------- |
| `Success` | `bool`   | Whether the update completed successfully.        |
| `Message` | `string` | Human-readable message describing the outcome.    |
| `Error`   | `string` | Error message if update failed; empty on success. |

## Idempotency

**Idempotent.** Checks the current DNS configuration before applying changes.
Returns `Changed: true` only if the servers or search domains were actually
modified. If the configuration already matches the desired state, the operation
returns `Changed: false`.

## Permissions

Requires `network:write` permission.

## Example

See
[`examples/dns-update.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/dns-update.go)
for a complete working example.
