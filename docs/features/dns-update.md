# DNS update

This page shows the pattern that reads, then writes. Read the current DNS
configuration. Then update it with new servers.

## Usage

```go
getDNS := o.DNSGet("_any", "eth0").After(health)

o.DNSUpdate("_any", osapi.DNSUpdateInput{
    InterfaceName: "eth0",
    Servers:       []string{"8.8.8.8", "8.8.4.4"},
    SearchDomains: []string{"example.com"},
}).After(getDNS)
```

The update step depends on the read step. Thus the plan reads the DNS
configuration before it changes it.

## Example

Refer to
[`examples/operations/dns-update.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/dns-update.go)
for a complete example.
