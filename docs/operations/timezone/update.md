# TimezoneUpdate

Sets the system timezone on the target node.

## Usage

```go
step := o.TimezoneUpdate("web-01", osapi.TimezoneUpdateOpts{
    Timezone: "America/New_York",
})
```

## Parameters

| Parameter | Type                 | Description                                               |
| --------- | -------------------- | --------------------------------------------------------- |
| `target`  | `string`             | Target host: `_any`, `_all`, hostname, or label selector. |
| `opts`    | `TimezoneUpdateOpts` | Update options (see below).                               |

### TimezoneUpdateOpts

| Field      | Type     | Required | Description                                        |
| ---------- | -------- | -------- | -------------------------------------------------- |
| `Timezone` | `string` | Yes      | IANA timezone name (e.g., `America/New_York`, `UTC`). |

## Result Type

```go
var result osapi.TimezoneMutationResult
err := results.Decode("update-timezone-1", &result)
```

| Field      | Type     | Description                                       |
| ---------- | -------- | ------------------------------------------------- |
| `Timezone` | `string` | The timezone that was set.                        |
| `Changed`  | `bool`   | Whether the timezone was actually modified.       |
| `Error`    | `string` | Error message if update failed; empty on success. |

## Idempotency

**Idempotent.** Compares the current timezone against the desired value.
Returns `Changed: true` only if the timezone was actually modified.

## Permissions

Requires `timezone:write` permission.

## Example

See
[`examples/operations/timezone.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/timezone.go)
for a complete working example.
