# NTPGet

Retrieves the NTP status and configuration from the target node.

## Usage

```go
step := o.NTPGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result Type

```go
var result osapi.NtpStatusResult
err := results.Decode("get-ntp-1", &result)
```

| Field           | Type       | Description                                      |
| --------------- | ---------- | ------------------------------------------------ |
| `Hostname`      | `string`   | The node's hostname.                             |
| `Synchronized`  | `bool`     | Whether NTP is synchronized.                     |
| `Stratum`       | `int`      | NTP stratum level.                               |
| `Offset`        | `string`   | Clock offset from the NTP source.                |
| `CurrentSource` | `string`   | The currently selected NTP source.               |
| `Servers`       | `[]string` | Configured NTP server addresses.                 |
| `Error`         | `string`   | Error message if query failed; empty on success. |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `ntp:read` permission.

## Example

See
[`examples/operations/ntp.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/ntp.go)
for a complete working example.
