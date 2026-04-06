# PackageGet

Retrieves information about a specific installed package.

## Usage

```go
step := o.PackageGet("web-01", "nginx")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `pkgName` | `string` | Name of the package to retrieve.                          |

## Result Type

```go
var result osapi.PackageInfoResult
err := results.Decode("get-package-1", &result)
```

| Field      | Type            | Description                                      |
| ---------- | --------------- | ------------------------------------------------ |
| `Hostname` | `string`        | The node's hostname.                             |
| `Packages` | `[]PackageInfo` | Package details (single entry).                  |
| `Error`    | `string`        | Error message if query failed; empty on success. |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `package:read` permission.

## Example

See
[`examples/operations/package.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/package.go)
for a complete working example.
