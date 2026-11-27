# TimezoneUpdate

Set the timezone of the host.

Sets the timezone with timedatectl set-timezone. A host that already has the
timezone does not change and reports changed: false. It runs on the Debian
family only, and elsewhere reports skipped.

## Usage

```go
step := o.TimezoneUpdate("web-01", input)
```

## Parameters

| Parameter | Type                        | Description                                                   |
| --------- | --------------------------- | ------------------------------------------------------------- |
| `target`  | `string`                    | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.TimezoneUpdateInput` | Refer to `TimezoneUpdateInput` below.                         |

### TimezoneUpdateInput

| Field      | Type     | Description                                                     |
| ---------- | -------- | --------------------------------------------------------------- |
| `Timezone` | `string` | The IANA name of the timezone, such as UTC or America/New_York. |

## Result

Each host answers with one `osapi.TimezoneResult` row. A list answers with one
row for each item.

```go
var row osapi.TimezoneResult
err := report.Decode("update-timezone", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`  | `bool`            | Whether the host changed.                                                                                                                                                                                                |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |
| `Timezone` | `string`          | The timezone of the host.                                                                                                                                                                                                |

## Permissions

The token must have the `timezone:write` permission. The step sends
`PUT /api/node/{hostname}/timezone`.

## Example

Refer to
[`examples/operations/timezone.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/timezone.go).
