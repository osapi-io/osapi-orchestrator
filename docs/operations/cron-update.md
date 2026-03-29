# CronUpdate

Updates an existing cron entry on the target node. Only the fields specified in
the options are changed; unset fields retain their current values.

## Usage

```go
step := o.CronUpdate("web-01", "backup", osapi.CronUpdateOpts{
    Schedule: "0 3 * * *",
})
```

## Parameters

| Parameter   | Type              | Description                                               |
| ----------- | ----------------- | --------------------------------------------------------- |
| `target`    | `string`          | Target host: `_any`, `_all`, hostname, or label selector. |
| `entryName` | `string`          | Name of the cron entry to update.                         |
| `opts`      | `CronUpdateOpts`  | Update options (see below).                               |

### CronUpdateOpts

| Field         | Type             | Required | Description                                 |
| ------------- | ---------------- | -------- | ------------------------------------------- |
| `Object`      | `string`         | No       | New object to deploy.                       |
| `Schedule`    | `string`         | No       | New cron expression.                        |
| `User`        | `string`         | No       | New user to run the command as.             |
| `ContentType` | `string`         | No       | `"raw"` or `"template"`.                    |
| `Vars`        | `map[string]any` | No       | Template variables when ContentType is      |
|               |                  |          | `"template"`.                               |

## Result Type

```go
var result osapi.CronMutationResult
err := results.Decode("update-cron-1", &result)
```

| Field     | Type     | Description                                       |
| --------- | -------- | ------------------------------------------------- |
| `Name`    | `string` | Name of the updated cron entry.                   |
| `Changed` | `bool`   | Whether the entry was actually modified.          |
| `Error`   | `string` | Error message if update failed; empty on success. |

## Idempotency

**Idempotent.** Compares the current entry state against the desired state.
Returns `Changed: true` only if the entry was actually modified. If the entry
already matches, returns `Changed: false`.

## Permissions

Requires `cron:write` permission.

## Example

```go
plan := o.Plan("update-schedule")
o.CronUpdate("web-01", "backup", osapi.CronUpdateOpts{
    Schedule: "0 3 * * *",
})
report := plan.Execute(ctx)
```
