# GroupDelete

Deletes a group from the target node.

## Usage

```go
step := o.GroupDelete("web-01", "developers")
```

## Parameters

| Parameter   | Type     | Description                                               |
| ----------- | -------- | --------------------------------------------------------- |
| `target`    | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `groupName` | `string` | Name of the group to delete.                              |

## Result Type

```go
var result osapi.GroupMutationResult
err := results.Decode("delete-group-1", &result)
```

| Field     | Type     | Description                                          |
| --------- | -------- | ---------------------------------------------------- |
| `Name`    | `string` | Group name.                                          |
| `Changed` | `bool`   | Whether the group was removed or already absent.     |
| `Error`   | `string` | Error message if deletion failed; empty on success.  |

## Idempotency

**Idempotent.** If the group does not exist, returns `Changed: false`.

## Permissions

Requires `user:write` permission.

## Example

See
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete working example.
