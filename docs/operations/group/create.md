# GroupCreate

Add a group.

Adds the group with groupadd, with the ID when one is given. A group that the
host has already does not change and reports changed: false, whatever its ID. It
runs on the Debian family only, and elsewhere reports skipped.

## Usage

```go
step := o.GroupCreate("web-01", input)
```

## Parameters

| Parameter | Type                     | Description                                                   |
| --------- | ------------------------ | ------------------------------------------------------------- |
| `target`  | `string`                 | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.GroupCreateInput` | Refer to `GroupCreateInput` below.                            |

### GroupCreateInput

| Field    | Type     | Description                                                 |
| -------- | -------- | ----------------------------------------------------------- |
| `GID`    | `int`    | Its numeric group ID; the system assigns one when it is 0.  |
| `Name`   | `string` | The group to add, such as developers.                       |
| `System` | `bool`   | Add it as a system group, with an ID from the system range. |

## Result

Each host answers with one `osapi.GroupResult` row. A list answers with one row
for each item.

```go
var row osapi.GroupResult
err := report.Decode("create-group", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`  | `bool`            | Whether the host changed.                                                                                                                                                                                                |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Name`     | `string`          | The group.                                                                                                                                                                                                               |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `group:write` permission. The step sends
`POST /api/node/{hostname}/group`.

## Example

Refer to
[`examples/operations/group.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/group.go).
