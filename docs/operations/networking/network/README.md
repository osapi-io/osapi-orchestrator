# Network management

Query and update the network configuration. The operations get and set the DNS
servers, and use ICMP ping to check connectivity.

## Operations

| Method                                                   | Description                            | Idempotent |
| -------------------------------------------------------- | -------------------------------------- | ---------- |
| [`NetworkDNSGet(target, iface)`](dns-get.md)             | Get DNS configuration for an interface | Read-only  |
| [`NetworkDNSUpdate(target, iface, opts)`](dns-update.md) | Update DNS servers for an interface    | Idempotent |
| [`NetworkDNSDelete(target, iface)`](dns-delete.md)       | Delete DNS configuration               | Idempotent |
| [`NetworkPingDo(target, address)`](ping.md)              | Ping an address from the target host   | Read-only  |

## Permissions

| Operation        | Permission      |
| ---------------- | --------------- |
| Read operations  | `network:read`  |
| Write operations | `network:write` |

## Examples

Refer to
[`examples/operations/dns-update.go`](../../../../examples/operations/dns-update.go)
for a DNS read-then-write workflow. Refer to
[`examples/operations/ping.go`](../../../../examples/operations/ping.go) for a
connectivity check.
