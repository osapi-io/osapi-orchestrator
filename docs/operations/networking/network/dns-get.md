# NetworkDNSGet

Gets the DNS configuration for a specific network interface on the target host.
The result has the configured DNS servers and search domains.

## Usage

```go
step := o.NetworkDNSGet("web-01", "eth0")
```

## Parameters

| Parameter       | Type     | Description                                                |
| --------------- | -------- | ---------------------------------------------------------- |
| `target`        | `string` | Target host: `_any`, `_all`, hostname, or label selector.  |
| `interfaceName` | `string` | Network interface to query (for example, `eth0`, `ens33`). |

## Result type

```go
var result osapi.DNSConfig
err := results.Decode("network.dns.get-1", &result)
```

| Field           | Type       | Description                              |
| --------------- | ---------- | ---------------------------------------- |
| `DNSServers`    | `[]string` | List of configured DNS server addresses. |
| `SearchDomains` | `[]string` | List of DNS search domains.              |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `network:read` permission.

## Example

Refer to
[`examples/operations/dns-update.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/dns-update.go)
for a complete example.
