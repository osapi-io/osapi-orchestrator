# ConnectivityPing

Ping an address from the host.

Sends three ICMP echo requests from the host and reports how many came back, and
their round-trip times. A host that does not answer is 100% loss, not an error.
It changes nothing, so it reports changed: false. On Linux the agent pings
without root when net.ipv4.ping_group_range admits its group.

## Usage

```go
step := o.ConnectivityPing("web-01", input)
```

## Parameters

| Parameter | Type                          | Description                                                   |
| --------- | ----------------------------- | ------------------------------------------------------------- |
| `target`  | `string`                      | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.ConnectivityPingInput` | Refer to `ConnectivityPingInput` below.                       |

### ConnectivityPingInput

| Field     | Type     | Description                                                                       |
| --------- | -------- | --------------------------------------------------------------------------------- |
| `Address` | `string` | The IPv4 or IPv6 address to ping, or an @fact. reference that the agent resolves. |

## Result

Each host answers with one `osapi.ConnectivityPingResult` row. A list answers
with one row for each item.

```go
var row osapi.ConnectivityPingResult
err := report.Decode("ping-connectivity", &row)
```

| Field             | Type              | Description                                                                                                                                                                                                              |
| ----------------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `AvgRTT`          | `string`          | Average round-trip time in milliseconds, with two decimals. Absent when no reply came back.                                                                                                                              |
| `Changed`         | `bool`            | Always false: a ping changes nothing on the host.                                                                                                                                                                        |
| `Error`           | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname`        | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `MaxRTT`          | `string`          | Maximum round-trip time in milliseconds, with two decimals. Absent when no reply came back.                                                                                                                              |
| `MinRTT`          | `string`          | Minimum round-trip time in milliseconds, with two decimals. Absent when no reply came back.                                                                                                                              |
| `PacketLoss`      | `float64`         | Percentage of packet loss.                                                                                                                                                                                               |
| `PacketsReceived` | `int`             | Number of packets received.                                                                                                                                                                                              |
| `PacketsSent`     | `int`             | Number of packets sent.                                                                                                                                                                                                  |
| `Status`          | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `connectivity:execute` permission. The step sends
`POST /api/node/{hostname}/connectivity/ping`.

## Example

Refer to
[`examples/operations/ping.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/ping.go).
