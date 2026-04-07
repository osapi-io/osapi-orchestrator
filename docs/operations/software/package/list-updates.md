# PackageListUpdates

Lists available package updates on the target node.

## Usage

```go
step := o.PackageListUpdates("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result Type

```go
var result osapi.PackageUpdateResult
err := results.Decode("list-package-updates-1", &result)
```

| Field      | Type           | Description                                      |
| ---------- | -------------- | ------------------------------------------------ |
| `Hostname` | `string`       | The node's hostname.                             |
| `Updates`  | `[]UpdateInfo` | List of available updates.                       |
| `Error`    | `string`       | Error message if query failed; empty on success. |

### UpdateInfo

| Field            | Type     | Description                  |
| ---------------- | -------- | ---------------------------- |
| `Name`           | `string` | Package name.                |
| `CurrentVersion` | `string` | Currently installed version. |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `package:read` permission.

## Example

See
[`examples/operations/package.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/package.go)
for a complete working example.
