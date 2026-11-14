# CronUpdate

Updates a cron entry on the target host. The operation changes only the fields
that you set in the options. The fields that you do not set keep their current
values.

## Usage

```go
step := o.CronUpdate("web-01", "backup", osapi.CronUpdateOpts{
    Schedule: "0 3 * * *",
})
```

## Parameters

| Parameter   | Type             | Description                                               |
| ----------- | ---------------- | --------------------------------------------------------- |
| `target`    | `string`         | Target host: `_any`, `_all`, hostname, or label selector. |
| `entryName` | `string`         | The name of the cron entry to update.                     |
| `opts`      | `CronUpdateOpts` | Update options (refer to the table below).                |

### CronUpdateOpts

| Field         | Type             | Required | Description                                |
| ------------- | ---------------- | -------- | ------------------------------------------ |
| `Object`      | `string`         | No       | The new object to deploy.                  |
| `Schedule`    | `string`         | No       | The new cron expression.                   |
| `User`        | `string`         | No       | The new user to run the command as.        |
| `ContentType` | `string`         | No       | `"raw"` or `"template"`.                   |
| `Vars`        | `map[string]any` | No       | The template variables when ContentType is |
|               |                  |          | `"template"`.                              |

## Result type

```go
var result osapi.CronMutationResult
err := results.Decode("update-cron-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | The name of the updated cron entry.                          |
| `Changed` | `bool`   | `true` if the operation changed the entry.                   |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** The operation compares the current entry with the values you
give. It returns `Changed: true` only if it changed the entry. If the entry
already matches, the operation returns `Changed: false`.

## Permissions

The token must have the `cron:write` permission.

## Example

Refer to
[`examples/operations/cron.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/cron.go)
for a complete example.
