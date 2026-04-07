# UserCreate

Creates a user account on the target node.

## Usage

```go
step := o.UserCreate("web-01", osapi.UserCreateOpts{
    Name:  "deploy",
    Shell: "/bin/bash",
    Home:  "/home/deploy",
})
```

## Parameters

| Parameter | Type             | Description                                               |
| --------- | ---------------- | --------------------------------------------------------- |
| `target`  | `string`         | Target host: `_any`, `_all`, hostname, or label selector. |
| `opts`    | `UserCreateOpts` | Create options (see below).                               |

### UserCreateOpts

| Field   | Type     | Required | Description                                        |
| ------- | -------- | -------- | -------------------------------------------------- |
| `Name`  | `string` | Yes      | Username.                                          |
| `UID`   | `int`    | No       | Numeric user ID (system assigns if zero).          |
| `GID`   | `int`    | No       | Primary group ID (creates matching group if zero). |
| `Home`  | `string` | No       | Home directory path.                               |
| `Shell` | `string` | No       | Login shell path.                                  |

## Result Type

```go
var result osapi.UserMutationResult
err := results.Decode("create-user-1", &result)
```

| Field     | Type     | Description                                         |
| --------- | -------- | --------------------------------------------------- |
| `Name`    | `string` | Username.                                           |
| `Changed` | `bool`   | Whether the user was created.                       |
| `Error`   | `string` | Error message if creation failed; empty on success. |

## Idempotency

**Non-idempotent.** Creating a user that already exists returns an error.

## Permissions

Requires `user:write` permission.

## Example

See
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete working example.
