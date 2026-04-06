# GroupCreate

Creates a group on the target node.

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
| `opts`    | `GroupCreateOpts` | Create options (see below).                               |

### GroupCreateOpts

| Field    | Type     | Required | Description                                |
| -------- | -------- | -------- | ------------------------------------------ |
| `Name`   | `string` | Yes      | Group name.                                |
| `GID`    | `int`    | No       | Numeric group ID (system assigns if zero). |
| `System` | `bool`   | No       | Create a system group.                     |

## Result Type

```go
var result osapi.GroupMutationResult
err := results.Decode("create-group-1", &result)
```

| Field     | Type     | Description                                         |
| --------- | -------- | --------------------------------------------------- |
| `Name`    | `string` | Group name.                                         |
| `Changed` | `bool`   | Whether the group was created.                      |
| `Error`   | `string` | Error message if creation failed; empty on success. |

## Idempotency

**Non-idempotent.** Creating a group that already exists returns an error.

## Permissions

Requires `user:write` permission.

## Example

See
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete working example.
