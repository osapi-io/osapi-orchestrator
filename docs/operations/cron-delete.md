# CronDelete

Deletes a cron entry from the target node.

## Usage

```go
step := o.CronDelete("web-01", "backup")
```

## Parameters

| Parameter   | Type     | Description                                               |
| ----------- | -------- | --------------------------------------------------------- |
| `target`    | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `entryName` | `string` | Name of the cron entry to delete.                         |

## Result Type

```go
var result osapi.CronMutationResult
err := results.Decode("delete-cron-1", &result)
```

| Field     | Type     | Description                                         |
| --------- | -------- | --------------------------------------------------- |
| `Name`    | `string` | Name of the deleted cron entry.                     |
| `Changed` | `bool`   | Whether the entry was removed or already absent.    |
| `Error`   | `string` | Error message if deletion failed; empty on success. |

## Idempotency

**Idempotent.** If the entry does not exist, the operation returns
`Changed: false`. If the entry exists and is removed, returns `Changed: true`.

## Permissions

Requires `cron:write` permission.

## Example

```go
plan := o.Plan("remove-cron")
o.CronDelete("web-01", "backup")
report := plan.Execute(ctx)
```
