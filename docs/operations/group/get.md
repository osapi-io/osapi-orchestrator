# GroupGet

Retrieves a specific group by name from the target node.

## Usage

```go
step := o.GroupGet("web-01", "developers")
```

## Parameters

| Parameter   | Type     | Description                                               |
| ----------- | -------- | --------------------------------------------------------- |
| `target`    | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `groupName` | `string` | Name of the group to retrieve.                            |

## Result Type

```go
var result osapi.GroupInfoResult
err := results.Decode("get-group-1", &result)
```

| Field      | Type          | Description                                      |
| ---------- | ------------- | ------------------------------------------------ |
| `Hostname` | `string`      | The node's hostname.                             |
| `Groups`   | `[]GroupInfo` | Group details (single entry).                    |
| `Error`    | `string`      | Error message if query failed; empty on success. |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `user:read` permission.

## Example

See
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete working example.
