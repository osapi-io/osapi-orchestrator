# GroupDelete

Deletes a group from the target host.

## Usage

```go
step := o.GroupDelete("web-01", "developers")
```

## Parameters

| Parameter   | Type     | Description                                               |
| ----------- | -------- | --------------------------------------------------------- |
| `target`    | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `groupName` | `string` | Name of the group to delete.                              |

## Result type

```go
var result osapi.GroupMutationResult
err := results.Decode("delete-group-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | Group name.                                                  |
| `Changed` | `bool`   | `true` if the operation removed the group.                   |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** If the group does not exist, the operation returns
`Changed: false`. If the group exists, the operation removes it and returns
`Changed: true`.

## Permissions

The token must have the `user:write` permission.

## Example

Refer to
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete example.
