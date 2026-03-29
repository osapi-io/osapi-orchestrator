# Network Management

Query and update network configuration -- DNS server settings and ICMP ping
connectivity checks.

## Operations

| Method                                                   | Description                            | Idempotent |
| -------------------------------------------------------- | -------------------------------------- | ---------- |
| [`NetworkDNSGet(target, iface)`](dns-get.md)             | Get DNS configuration for an interface | Read-only  |
| [`NetworkDNSUpdate(target, iface, opts)`](dns-update.md) | Update DNS servers for an interface    | Idempotent |
| [`NetworkPingDo(target, address)`](ping.md)              | Ping a host from the target node       | Read-only  |

## Permissions

| Operation        | Permission      |
| ---------------- | --------------- |
| Read operations  | `network:read`  |
| Write operations | `network:write` |

## Examples

See
[`examples/operations/dns-update.go`](../../examples/operations/dns-update.go)
for DNS read-then-write and
[`examples/operations/ping.go`](../../examples/operations/ping.go) for
connectivity verification.
