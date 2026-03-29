# NetworkPingDo

Pings a host from the target node and returns packet loss statistics. Useful for
verifying network connectivity between nodes.

## Usage

```go
step := o.NetworkPingDo("web-01", "8.8.8.8")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `address` | `string` | Hostname or IP address to ping.                           |

## Result Type

```go
var result osapi.PingResult
err := results.Decode("network.ping.do-1", &result)
```

| Field             | Type      | Description                                     |
| ----------------- | --------- | ----------------------------------------------- |
| `PacketsSent`     | `int`     | Number of ICMP packets sent.                    |
| `PacketsReceived` | `int`     | Number of ICMP packets received.                |
| `PacketLoss`      | `float64` | Percentage of packets lost (0-100).             |
| `Error`           | `string`  | Error message if ping failed; empty on success. |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `network:read` permission.

## Example

See
[`examples/operations/dns-update.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/dns-update.go)
for a complete working example.
