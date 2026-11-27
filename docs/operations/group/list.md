# GroupList

List the groups of the host.

Lists each group with its ID and its members, by name, from /etc/group through
gohai, on every platform. The members are the users that /etc/group names; a
user whose primary group it is is not listed again.

## Usage

```go
step := o.GroupList("web-01")
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |

## Result

Each host answers with one `osapi.GroupAccount` row. A list answers with one row
for each item.

```go
var row osapi.GroupAccount
err := report.Decode("list-group", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `GID`      | `int`             | Its numeric group ID.                                                                                                                                                                                                    |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Members`  | `[]string`        | The users it has as members, besides those whose primary group it is.                                                                                                                                                    |
| `Name`     | `string`          | The group.                                                                                                                                                                                                               |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `group:read` permission. The step sends
`GET /api/node/{hostname}/group`.

## Example

Refer to
[`examples/operations/group.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/group.go).
