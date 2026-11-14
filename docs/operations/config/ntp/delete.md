# NTPDelete

Deletes the NTP configuration from the target host.

## Usage

```go
step := o.NTPDelete("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result type

```go
var result osapi.NtpMutationResult
err := results.Decode("delete-ntp-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Changed` | `bool`   | `true` if the operation removed the configuration.           |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** If no NTP configuration exists, the operation returns
`Changed: false`. If the configuration exists, the operation removes it and
returns `Changed: true`.

## Permissions

The token must have the `ntp:write` permission.

## Example

Refer to
[`examples/operations/ntp.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/ntp.go)
for a complete example.
