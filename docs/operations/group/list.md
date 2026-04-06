# GroupList

Lists groups on the target node.

## Usage

```go
step := o.GroupList("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result Type

```go
var result osapi.GroupInfoResult
err := results.Decode("list-group-1", &result)
```

| Field      | Type          | Description                                      |
| ---------- | ------------- | ------------------------------------------------ |
| `Hostname` | `string`      | The node's hostname.                             |
| `Groups`   | `[]GroupInfo` | List of groups.                                  |
| `Error`    | `string`      | Error message if query failed; empty on success. |

### GroupInfo

| Field     | Type       | Description         |
| --------- | ---------- | ------------------- |
| `Name`    | `string`   | Group name.         |
| `GID`     | `int`      | Numeric group ID.   |
| `Members` | `[]string` | Group member names. |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `user:read` permission.

## Example

See
[`examples/operations/user.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/user.go)
for a complete working example.
