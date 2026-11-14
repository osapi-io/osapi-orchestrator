# GroupList

Lists groups on the target host.

## Usage

```go
step := o.GroupList("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result type

```go
var result osapi.GroupInfoResult
err := results.Decode("list-group-1", &result)
```

| Field      | Type          | Description                                                  |
| ---------- | ------------- | ------------------------------------------------------------ |
| `Hostname` | `string`      | The hostname of the agent that ran the operation.            |
| `Groups`   | `[]GroupInfo` | List of groups.                                              |
| `Error`    | `string`      | The error message if the operation failed. Empty on success. |

### GroupInfo

| Field     | Type       | Description         |
| --------- | ---------- | ------------------- |
| `Name`    | `string`   | Group name.         |
| `GID`     | `int`      | Numeric group ID.   |
| `Members` | `[]string` | Group member names. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `user:read` permission.

## Example

Refer to
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete example.
