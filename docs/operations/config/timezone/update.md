# TimezoneUpdate

Sets the system timezone on the target host.

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
| `opts`    | `TimezoneUpdateOpts` | Update options (refer to the table below).                |

### TimezoneUpdateOpts

| Field      | Type     | Required | Description                                                      |
| ---------- | -------- | -------- | ---------------------------------------------------------------- |
| `Timezone` | `string` | Yes      | The IANA timezone name (for example, `America/New_York`, `UTC`). |

## Result type

```go
var result osapi.TimezoneMutationResult
err := results.Decode("update-timezone-1", &result)
```

| Field      | Type     | Description                                                  |
| ---------- | -------- | ------------------------------------------------------------ |
| `Timezone` | `string` | The timezone that the operation set.                         |
| `Changed`  | `bool`   | `true` if the operation changed the timezone.                |
| `Error`    | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** The operation compares the current timezone with the timezone
you give. It returns `Changed: true` only if it changed the timezone.

## Permissions

The token must have the `timezone:write` permission.

## Example

Refer to
[`examples/operations/timezone.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/timezone.go)
for a complete example.
