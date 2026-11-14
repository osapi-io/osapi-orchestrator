# UserCreate

Creates a user account on the target host.

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
| `opts`    | `UserCreateOpts` | Create options (refer to the table below).                |

### UserCreateOpts

| Field   | Type     | Required | Description                                                         |
| ------- | -------- | -------- | ------------------------------------------------------------------- |
| `Name`  | `string` | Yes      | Username.                                                           |
| `UID`   | `int`    | No       | Numeric user ID. If zero, the system assigns one.                   |
| `GID`   | `int`    | No       | Primary group ID. If zero, the system creates a group for the user. |
| `Home`  | `string` | No       | Home directory path.                                                |
| `Shell` | `string` | No       | Login shell path.                                                   |

## Result type

```go
var result osapi.UserMutationResult
err := results.Decode("create-user-1", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `Name`    | `string` | Username.                                                    |
| `Changed` | `bool`   | `true` if the operation created the user.                    |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Non-idempotent.** If the user already exists, the operation returns an error.

## Permissions

The token must have the `user:write` permission.

## Example

Refer to
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete example.
