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

## Example

See
[`examples/operations/dns-update.go`](../../examples/operations/dns-update.go)
for a complete workflow example covering all operations.

The example demonstrates:

- Read-then-write pattern with DNS configuration
- Broadcasting DNS updates across multiple hosts
- Verifying changes with a follow-up read
