# CronDelete

Deletes a cron entry from the target host.

## Usage

```go
step := o.CronDelete("web-01", "backup")
```

## Parameters

| Parameter   | Type     | Description                                               |
| ----------- | -------- | --------------------------------------------------------- |
| `target`    | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `entryName` | `string` | Name of the cron entry to delete.                         |

## Result type

```go
var result osapi.CronMutationResult
err := results.Decode("delete-cron-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | Name of the deleted cron entry.                              |
| `Changed` | `bool`   | `true` if the operation removed the entry.                   |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** If the entry does not exist, the operation returns
`Changed: false`. If the entry exists, the operation removes it and returns
`Changed: true`.

## Permissions

The token must have the `cron:write` permission.

## Example

Refer to
[`examples/operations/cron.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/cron.go)
for a complete example.
