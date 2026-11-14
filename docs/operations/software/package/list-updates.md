# PackageListUpdates

Lists available package updates on the target host.

## Usage

```go
step := o.PackageListUpdates("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result type

```go
var result osapi.PackageUpdateResult
err := results.Decode("list-package-updates-1", &result)
```

| Field      | Type           | Description                                                  |
| ---------- | -------------- | ------------------------------------------------------------ |
| `Hostname` | `string`       | The hostname of the agent that ran the operation.            |
| `Updates`  | `[]UpdateInfo` | List of available updates.                                   |
| `Error`    | `string`       | The error message if the operation failed. Empty on success. |

### UpdateInfo

| Field            | Type     | Description                        |
| ---------------- | -------- | ---------------------------------- |
| `Name`           | `string` | The package name.                  |
| `CurrentVersion` | `string` | The version that is installed now. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `package:read` permission.

## Example

Refer to
[`examples/operations/package.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/package.go)
for a complete example.
