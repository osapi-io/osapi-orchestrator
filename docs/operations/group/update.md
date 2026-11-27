# GroupUpdate

Set the members of a group.

Sets the members with gpasswd -M, in place of the members that the group has.
The same members, in any order, do not change it and report changed: false. A
group that the host does not have is not found. It runs on the Debian family
only, and elsewhere reports skipped.

## Usage

```go
step := o.GroupUpdate("web-01", input)
```

## Parameters

| Parameter | Type                     | Description                                                   |
| --------- | ------------------------ | ------------------------------------------------------------- |
| `target`  | `string`                 | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.GroupUpdateInput` | Refer to `GroupUpdateInput` below.                            |

### GroupUpdateInput

| Field     | Type       | Description                                                                                   |
| --------- | ---------- | --------------------------------------------------------------------------------------------- |
| `Name`    | `string`   | The name of the path: the group to remove, such as developers. It names the item in the path. |
| `Members` | `[]string` | Its members, which replace the members it has; none when empty.                               |

## Result

Each host answers with one `osapi.GroupResult` row. A list answers with one row
for each item.

```go
var row osapi.GroupResult
err := report.Decode("update-group", &row)
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
`PUT /api/node/{hostname}/group/{name}`.

## Example

Refer to
[`examples/operations/group.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/group.go).
