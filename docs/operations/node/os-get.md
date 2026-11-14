# NodeOSGet

Gets the OS information from the target host. The result has the distribution
and the version.

## Usage

```go
step := o.NodeOSGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result type

```go
var result osapi.OSInfoResult
err := results.Decode("get-os-1", &result)
```

| Field      | Type      | Description                                                  |
| ---------- | --------- | ------------------------------------------------------------ |
| `Hostname` | `string`  | The hostname of the agent that ran the operation.            |
| `OSInfo`   | `*OSInfo` | OS distribution and version (refer to the table below).      |
| `Error`    | `string`  | The error message if the operation failed. Empty on success. |

### OSInfo

| Field          | Type     | Description                              |
| -------------- | -------- | ---------------------------------------- |
| `Distribution` | `string` | OS distribution (for example, `Ubuntu`). |
| `Version`      | `string` | OS version (for example, `24.04`).       |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `node:read` permission.

## Example

Refer to
[`examples/operations/node-info.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/node-info.go)
for a complete example.
