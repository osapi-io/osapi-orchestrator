# GroupUpdate

Updates a group on the target node.

## Usage

```go
step := o.GroupUpdate("web-01", "developers", osapi.GroupUpdateOpts{
    Members: []string{"alice", "bob"},
})
```

## Parameters

| Parameter   | Type              | Description                                               |
| ----------- | ----------------- | --------------------------------------------------------- |
| `target`    | `string`          | Target host: `_any`, `_all`, hostname, or label selector. |
| `groupName` | `string`          | Name of the group to update.                              |
| `opts`      | `GroupUpdateOpts` | Update options (see below).                               |

### GroupUpdateOpts

| Field     | Type       | Required | Description                                    |
| --------- | ---------- | -------- | ---------------------------------------------- |
| `Members` | `[]string` | No       | Group member usernames (replaces existing).    |

## Result Type

```go
var result osapi.GroupMutationResult
err := results.Decode("update-group-1", &result)
```

| Field     | Type     | Description                                       |
| --------- | -------- | ------------------------------------------------- |
| `Name`    | `string` | Group name.                                       |
| `Changed` | `bool`   | Whether the group was actually modified.          |
| `Error`   | `string` | Error message if update failed; empty on success. |

## Idempotency

**Idempotent.** Compares the current members against the desired list. Returns
`Changed: true` only if the group was actually modified.

## Permissions

Requires `user:write` permission.

## Example

See
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete working example.
