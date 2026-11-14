# GroupCreate

Creates a group on the target host.

## Usage

```go
step := o.GroupCreate("web-01", osapi.GroupCreateOpts{
    Name: "developers",
})
```

## Parameters

| Parameter | Type              | Description                                               |
| --------- | ----------------- | --------------------------------------------------------- |
| `target`  | `string`          | Target host: `_any`, `_all`, hostname, or label selector. |
| `opts`    | `GroupCreateOpts` | Create options (refer to the table below).                |

### GroupCreateOpts

| Field    | Type     | Required | Description                                        |
| -------- | -------- | -------- | -------------------------------------------------- |
| `Name`   | `string` | Yes      | Group name.                                        |
| `GID`    | `int`    | No       | Numeric group ID. If zero, the system assigns one. |
| `System` | `bool`   | No       | Create a system group.                             |

## Result type

```go
var result osapi.GroupMutationResult
err := results.Decode("create-group-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | Group name.                                                  |
| `Changed` | `bool`   | `true` if the operation created the group.                   |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Non-idempotent.** If the group already exists, the operation returns an error.

## Permissions

The token must have the `user:write` permission.

## Example

Refer to
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete example.
