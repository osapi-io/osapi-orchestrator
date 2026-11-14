# GroupUpdate

Updates a group on the target host.

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
| `opts`      | `GroupUpdateOpts` | Update options (refer to the table below).                |

### GroupUpdateOpts

| Field     | Type       | Required | Description                                                     |
| --------- | ---------- | -------- | --------------------------------------------------------------- |
| `Members` | `[]string` | No       | Group member usernames. This list replaces the current members. |

## Result type

```go
var result osapi.GroupMutationResult
err := results.Decode("update-group-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | Group name.                                                  |
| `Changed` | `bool`   | `true` if the operation changed the group.                   |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** The operation compares the current members with the list you
give. It returns `Changed: true` only if it changed the group.

## Permissions

The token must have the `user:write` permission.

## Example

Refer to
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete example.
