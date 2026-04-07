# TimezoneGet

Retrieves the system timezone from the target node.

## Usage

```go
step := o.TimezoneGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result Type

```go
var result osapi.TimezoneResult
err := results.Decode("get-timezone-1", &result)
```

| Field       | Type     | Description                                      |
| ----------- | -------- | ------------------------------------------------ |
| `Hostname`  | `string` | The node's hostname.                             |
| `Timezone`  | `string` | IANA timezone name (e.g., `America/New_York`).   |
| `UTCOffset` | `string` | UTC offset string (e.g., `-0500`).               |
| `Error`     | `string` | Error message if query failed; empty on success. |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `timezone:read` permission.

## Example

See
[`examples/operations/timezone.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/timezone.go)
for a complete working example.
