# NetworkPingDo

Pings an address from the target host and returns packet loss statistics. Use it
to make sure that hosts can reach each other on the network.

## Usage

```go
step := o.NetworkPingDo("web-01", "8.8.8.8")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `address` | `string` | The hostname or IP address to ping.                       |

## Result type

```go
var result osapi.PingResult
err := results.Decode("network.ping.do-1", &result)
```

| Field             | Type      | Description                                                  |
| ----------------- | --------- | ------------------------------------------------------------ |
| `PacketsSent`     | `int`     | The number of ICMP packets that the host sent.               |
| `PacketsReceived` | `int`     | The number of ICMP packets that the host received.           |
| `PacketLoss`      | `float64` | The percentage of packets that the host lost (0-100).        |
| `Error`           | `string`  | The error message if the operation failed. Empty on success. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `network:read` permission.

## Example

Refer to
[`examples/operations/ping.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/ping.go)
for a complete example.
