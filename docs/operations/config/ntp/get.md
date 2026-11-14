# NTPGet

Gets the NTP status and configuration from the target host.

## Usage

```go
step := o.NTPGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result type

```go
var result osapi.NtpStatusResult
err := results.Decode("get-ntp-1", &result)
```

| Field           | Type       | Description                                                  |
| --------------- | ---------- | ------------------------------------------------------------ |
| `Hostname`      | `string`   | The hostname of the agent that ran the operation.            |
| `Synchronized`  | `bool`     | `true` if NTP is synchronized.                               |
| `Stratum`       | `int`      | The NTP stratum level.                                       |
| `Offset`        | `string`   | The clock offset from the NTP source.                        |
| `CurrentSource` | `string`   | The NTP source that NTP uses now.                            |
| `Servers`       | `[]string` | The addresses of the configured NTP servers.                 |
| `Error`         | `string`   | The error message if the operation failed. Empty on success. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `ntp:read` permission.

## Example

Refer to
[`examples/operations/ntp.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/ntp.go)
for a complete example.
