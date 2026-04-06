# NTPUpdate

Updates the NTP configuration on the target node.

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
| `opts`    | `NtpUpdateOpts` | Update options (see below).                               |

### NtpUpdateOpts

| Field     | Type       | Required | Description                         |
| --------- | ---------- | -------- | ----------------------------------- |
| `Servers` | `[]string` | Yes      | NTP server addresses to configure.  |

## Result Type

```go
var result osapi.NtpMutationResult
err := results.Decode("update-ntp-1", &result)
```

| Field     | Type     | Description                                       |
| --------- | -------- | ------------------------------------------------- |
| `Changed` | `bool`   | Whether the configuration was actually modified.  |
| `Error`   | `string` | Error message if update failed; empty on success. |

## Idempotency

**Idempotent.** Compares the current NTP servers against the desired list.
Returns `Changed: true` only if the configuration was actually modified.

## Permissions

Requires `ntp:write` permission.

## Example

See
[`examples/operations/ntp.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/ntp.go)
for a complete working example.
