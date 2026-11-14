# NetworkDNSUpdate

Updates the DNS configuration for a specific network interface on the target
host. The operation sets the DNS servers and search domains to the values you
give.

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

| Parameter       | Type       | Description                                                    |
| --------------- | ---------- | -------------------------------------------------------------- |
| `target`        | `string`   | Target host: `_any`, `_all`, hostname, or label selector.      |
| `interfaceName` | `string`   | Network interface to configure (for example, `eth0`, `ens33`). |
| `servers`       | `[]string` | DNS server addresses to set.                                   |
| `searchDomains` | `[]string` | DNS search domains to set.                                     |

## Result type

```go
var result osapi.DNSUpdateResult
err := results.Decode("network.dns.update-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Success` | `bool`   | `true` if the update succeeded.                              |
| `Message` | `string` | A message about the result, for a person to read.            |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** The operation compares the current DNS configuration with the
configuration you give. It returns `Changed: true` only if it changed the
servers or search domains. If the configuration already matches, the operation
returns `Changed: false`.

## Permissions

The token must have the `network:write` permission.

## Example

Refer to
[`examples/operations/dns-update.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/dns-update.go)
for a complete example.
