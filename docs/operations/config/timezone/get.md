# TimezoneGet

Gets the system timezone from the target host.

## Usage

```go
step := o.TimezoneGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result type

```go
var result osapi.TimezoneResult
err := results.Decode("get-timezone-1", &result)
```

| Field       | Type     | Description                                                  |
| ----------- | -------- | ------------------------------------------------------------ |
| `Hostname`  | `string` | The hostname of the agent that ran the operation.            |
| `Timezone`  | `string` | The IANA timezone name (for example, `America/New_York`).    |
| `UTCOffset` | `string` | The UTC offset (for example, `-0500`).                       |
| `Error`     | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `timezone:read` permission.

## Example

Refer to
[`examples/operations/timezone.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/timezone.go)
for a complete example.
