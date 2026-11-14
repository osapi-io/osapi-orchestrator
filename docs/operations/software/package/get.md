# PackageGet

Gets information about a specific installed package.

## Usage

```go
step := o.PackageGet("web-01", "nginx")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `pkgName` | `string` | Name of the package to get.                               |

## Result type

```go
var result osapi.PackageInfoResult
err := results.Decode("get-package-1", &result)
```

| Field      | Type            | Description                                                  |
| ---------- | --------------- | ------------------------------------------------------------ |
| `Hostname` | `string`        | The hostname of the agent that ran the operation.            |
| `Packages` | `[]PackageInfo` | Package details (one entry).                                 |
| `Error`    | `string`        | The error message if the operation failed. Empty on success. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `package:read` permission.

## Example

Refer to
[`examples/operations/package.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/package.go)
for a complete example.
