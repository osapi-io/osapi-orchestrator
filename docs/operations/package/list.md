# PackageList

Lists installed packages on the target node.

## Usage

```go
step := o.PackageList("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result Type

```go
var result osapi.PackageInfoResult
err := results.Decode("list-package-1", &result)
```

| Field      | Type            | Description                                      |
| ---------- | --------------- | ------------------------------------------------ |
| `Hostname` | `string`        | The node's hostname.                             |
| `Packages` | `[]PackageInfo` | List of installed packages.                      |
| `Error`    | `string`        | Error message if query failed; empty on success. |

### PackageInfo

| Field         | Type     | Description               |
| ------------- | -------- | ------------------------- |
| `Name`        | `string` | Package name.             |
| `Version`     | `string` | Installed version.        |
| `Description` | `string` | Package description.      |
| `Status`      | `string` | Installation status.      |
| `Size`        | `int64`  | Installed size in bytes.  |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `package:read` permission.

## Example

See
[`examples/operations/package.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/package.go)
for a complete working example.
