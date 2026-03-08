# DNS Update

The read-then-write pattern: read current DNS configuration, then update it with
new servers.

## Usage

```go
getDNS := o.NetworkDNSGet("_any", "eth0").After(health)

o.NetworkDNSUpdate(
    "_any",
    "eth0",
    []string{"8.8.8.8", "8.8.4.4"},
    []string{"example.com"},
).After(getDNS)
```

The update step depends on the read step, ensuring DNS is queried before it is
modified.

## Example

See
[`examples/operations/dns-update.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/dns-update.go)
for a complete working example.
