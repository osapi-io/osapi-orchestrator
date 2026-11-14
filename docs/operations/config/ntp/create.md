# NTPCreate

Creates the NTP configuration on the target host. The operation deploys a
timesyncd configuration file that has the servers you give.

## Usage

```go
step := o.NTPCreate("web-01", osapi.NtpCreateOpts{
    Servers: []string{"0.pool.ntp.org", "1.pool.ntp.org"},
})
```

## Parameters

| Parameter | Type            | Description                                               |
| --------- | --------------- | --------------------------------------------------------- |
| `target`  | `string`        | Target host: `_any`, `_all`, hostname, or label selector. |
| `opts`    | `NtpCreateOpts` | Create options (refer to the table below).                |

### NtpCreateOpts

| Field     | Type       | Required | Description                       |
| --------- | ---------- | -------- | --------------------------------- |
| `Servers` | `[]string` | Yes      | The addresses of the NTP servers. |

## Result type

```go
var result osapi.NtpMutationResult
err := results.Decode("create-ntp-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Changed` | `bool`   | `true` if the operation created the configuration.           |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Non-idempotent.** If the configuration already exists, the operation returns
an error. Use [NTPUpdate](update.md) to change a configuration that exists.

## Permissions

The token must have the `ntp:write` permission.

## Example

Refer to
[`examples/operations/ntp.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/ntp.go)
for a complete example.
