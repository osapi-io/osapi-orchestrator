# NTPCreate

Creates NTP configuration on the target node by deploying a timesyncd
configuration file with the specified servers.

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
| `opts`    | `NtpCreateOpts` | Create options (see below).                               |

### NtpCreateOpts

| Field     | Type       | Required | Description                        |
| --------- | ---------- | -------- | ---------------------------------- |
| `Servers` | `[]string` | Yes      | NTP server addresses to configure. |

## Result Type

```go
var result osapi.NtpMutationResult
err := results.Decode("create-ntp-1", &result)
```

| Field     | Type     | Description                                         |
| --------- | -------- | --------------------------------------------------- |
| `Changed` | `bool`   | Whether the configuration was created.              |
| `Error`   | `string` | Error message if creation failed; empty on success. |

## Idempotency

**Non-idempotent.** Creating configuration that already exists returns an error.
Use [NTPUpdate](update.md) to modify existing configuration.

## Permissions

Requires `ntp:write` permission.

## Example

See
[`examples/operations/ntp.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/ntp.go)
for a complete working example.
