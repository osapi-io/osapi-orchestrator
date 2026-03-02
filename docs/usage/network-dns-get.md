# NetworkDNSGet

Retrieves the DNS configuration for a specific network interface on the target
node, including configured DNS servers and search domains.

## Usage

```go
step := o.NetworkDNSGet("web-01", "eth0")
```

## Parameters

| Parameter       | Type     | Description                                               |
| --------------- | -------- | --------------------------------------------------------- |
| `target`        | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `interfaceName` | `string` | Network interface to query (e.g., `eth0`, `ens33`).       |

## Result Type

```go
var result orchestrator.DNSConfigResult
err := results.Decode("network.dns.get-1", &result)
```

| Field           | Type       | Description                              |
| --------------- | ---------- | ---------------------------------------- |
| `DNSServers`    | `[]string` | List of configured DNS server addresses. |
| `SearchDomains` | `[]string` | List of DNS search domains.              |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `network:read` permission.
