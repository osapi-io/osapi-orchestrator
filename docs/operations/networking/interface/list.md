# InterfaceList

Lists network interfaces on the target node.

## Usage

```go
step := o.InterfaceList("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result Type

```go
var result osapi.InterfaceListResult
err := results.Decode("list-interface-1", &result)
```

| Field        | Type              | Description                                      |
| ------------ | ----------------- | ------------------------------------------------ |
| `Hostname`   | `string`          | The node's hostname.                             |
| `Interfaces` | `[]InterfaceInfo` | List of network interfaces.                      |
| `Error`      | `string`          | Error message if query failed; empty on success. |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `network:read` permission.

## Example

See
[`examples/operations/interface.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/interface.go)
for a complete working example.
