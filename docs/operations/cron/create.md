# CronCreate

Creates a new cron entry on the target node. The script must be uploaded to the
NATS Object Store first (see [FileUpload](../file/upload.md)).

## Usage

```go
step := o.CronCreate("web-01", osapi.CronCreateOpts{
    Name:     "backup",
    Object:   "backup.sh",
    Schedule: "0 2 * * *",
    User:     "root",
})
```

With interval-based scheduling:

```go
step := o.CronCreate("web-01", osapi.CronCreateOpts{
    Name:     "cleanup",
    Object:   "cleanup.sh",
    Interval: "daily",
})
```

## Parameters

| Parameter | Type              | Description                                               |
| --------- | ----------------- | --------------------------------------------------------- |
| `target`  | `string`          | Target host: `_any`, `_all`, hostname, or label selector. |
| `opts`    | `CronCreateOpts`  | Create options (see below).                               |

### CronCreateOpts

| Field         | Type             | Required | Description                                  |
| ------------- | ---------------- | -------- | -------------------------------------------- |
| `Name`        | `string`         | Yes      | Cron entry name.                             |
| `Object`      | `string`         | Yes      | Object name in the NATS Object Store.        |
| `Schedule`    | `string`         | No       | Cron expression (mutually exclusive with     |
|               |                  |          | Interval).                                   |
| `Interval`    | `string`         | No       | Periodic interval: `hourly`, `daily`,        |
|               |                  |          | `weekly`, or `monthly`.                      |
| `User`        | `string`         | No       | User to run the command as.                  |
| `ContentType` | `string`         | No       | `"raw"` or `"template"` (default: `"raw"`).  |
| `Vars`        | `map[string]any` | No       | Template variables when ContentType is       |
|               |                  |          | `"template"`.                                |

## Result Type

```go
var result osapi.CronMutationResult
err := results.Decode("create-cron-1", &result)
```

| Field     | Type     | Description                                        |
| --------- | -------- | -------------------------------------------------- |
| `Name`    | `string` | Name of the created cron entry.                    |
| `Changed` | `bool`   | Whether the entry was created.                     |
| `Error`   | `string` | Error message if creation failed; empty on success.|

## Idempotency

**Non-idempotent.** Creating an entry that already exists returns an error. Use
[CronUpdate](update.md) to modify existing entries.

## Permissions

Requires `cron:write` permission.

## Example

See
[`examples/operations/cron.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/cron.go)
for a complete working example.
