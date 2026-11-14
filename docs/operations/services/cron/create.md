# CronCreate

Creates a new cron entry on the target host. Upload the script to the NATS
Object Store first. Refer to [FileUpload](../../files/file/upload.md).

## Usage

```go
step := o.CronCreate("web-01", osapi.CronCreateOpts{
    Name:     "backup",
    Object:   "backup.sh",
    Schedule: "0 2 * * *",
    User:     "root",
})
```

To schedule with an interval:

```go
step := o.CronCreate("web-01", osapi.CronCreateOpts{
    Name:     "cleanup",
    Object:   "cleanup.sh",
    Interval: "daily",
})
```

## Parameters

| Parameter | Type             | Description                                               |
| --------- | ---------------- | --------------------------------------------------------- |
| `target`  | `string`         | Target host: `_any`, `_all`, hostname, or label selector. |
| `opts`    | `CronCreateOpts` | Create options (refer to the table below).                |

### CronCreateOpts

| Field         | Type             | Required | Description                                 |
| ------------- | ---------------- | -------- | ------------------------------------------- |
| `Name`        | `string`         | Yes      | The cron entry name.                        |
| `Object`      | `string`         | Yes      | The object name in the NATS Object Store.   |
| `Schedule`    | `string`         | No       | The cron expression. Do not use with        |
|               |                  |          | `Interval`.                                 |
| `Interval`    | `string`         | No       | The periodic interval: `hourly`, `daily`,   |
|               |                  |          | `weekly`, or `monthly`.                     |
| `User`        | `string`         | No       | The user to run the command as.             |
| `ContentType` | `string`         | No       | `"raw"` or `"template"` (default: `"raw"`). |
| `Vars`        | `map[string]any` | No       | The template variables when ContentType is  |
|               |                  |          | `"template"`.                               |

## Result type

```go
var result osapi.CronMutationResult
err := results.Decode("create-cron-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | The name of the created cron entry.                          |
| `Changed` | `bool`   | `true` if the operation created the entry.                   |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Non-idempotent.** If the entry already exists, the operation returns an error.
Use [CronUpdate](update.md) to change an entry that exists.

## Permissions

The token must have the `cron:write` permission.

## Example

Refer to
[`examples/operations/cron.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/cron.go)
for a complete example.
