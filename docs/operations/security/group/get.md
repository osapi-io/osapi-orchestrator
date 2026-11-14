# GroupGet

Gets a specific group by name from the target host.

## Usage

```go
step := o.GroupGet("web-01", "developers")
```

## Parameters

| Parameter   | Type     | Description                                               |
| ----------- | -------- | --------------------------------------------------------- |
| `target`    | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `groupName` | `string` | Name of the group to get.                                 |

## Result type

```go
var result osapi.GroupInfoResult
err := results.Decode("get-group-1", &result)
```

| Field      | Type          | Description                                                  |
| ---------- | ------------- | ------------------------------------------------------------ |
| `Hostname` | `string`      | The hostname of the agent that ran the operation.            |
| `Groups`   | `[]GroupInfo` | Group details (one entry).                                   |
| `Error`    | `string`      | The error message if the operation failed. Empty on success. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `user:read` permission.

## Example

Refer to
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete example.
