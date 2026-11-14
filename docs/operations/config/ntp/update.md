# NTPUpdate

Updates the NTP configuration on the target host.

## Usage

```go
step := o.NTPUpdate("web-01", osapi.NtpUpdateOpts{
    Servers: []string{"time.google.com", "time.cloudflare.com"},
})
```

## Parameters

| Parameter | Type            | Description                                               |
| --------- | --------------- | --------------------------------------------------------- |
| `target`  | `string`        | Target host: `_any`, `_all`, hostname, or label selector. |
| `opts`    | `NtpUpdateOpts` | Update options (refer to the table below).                |

### NtpUpdateOpts

| Field     | Type       | Required | Description                       |
| --------- | ---------- | -------- | --------------------------------- |
| `Servers` | `[]string` | Yes      | The addresses of the NTP servers. |

## Result type

```go
var result osapi.NtpMutationResult
err := results.Decode("update-ntp-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Changed` | `bool`   | `true` if the operation changed the configuration.           |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** The operation compares the current NTP servers with the servers
you give. It returns `Changed: true` only if it changed the configuration.

## Permissions

The token must have the `ntp:write` permission.

## Example

Refer to
[`examples/operations/ntp.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/ntp.go)
for a complete example.
