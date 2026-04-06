# NTPDelete

Deletes the NTP configuration from the target node.

## Usage

```go
step := o.NTPDelete("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result Type

```go
var result osapi.NtpMutationResult
err := results.Decode("delete-ntp-1", &result)
```

| Field     | Type     | Description                                           |
| --------- | -------- | ----------------------------------------------------- |
| `Changed` | `bool`   | Whether the configuration was removed or already absent. |
| `Error`   | `string` | Error message if deletion failed; empty on success.   |

## Idempotency

**Idempotent.** If no NTP configuration exists, the operation returns
`Changed: false`. If configuration exists and is removed, returns
`Changed: true`.

## Permissions

Requires `ntp:write` permission.

## Example

See
[`examples/operations/ntp.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/ntp.go)
for a complete working example.
