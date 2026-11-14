# PackageList

Lists installed packages on the target host.

## Usage

```go
step := o.PackageList("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result type

```go
var result osapi.PackageInfoResult
err := results.Decode("list-package-1", &result)
```

| Field      | Type            | Description                                                  |
| ---------- | --------------- | ------------------------------------------------------------ |
| `Hostname` | `string`        | The hostname of the agent that ran the operation.            |
| `Packages` | `[]PackageInfo` | List of installed packages.                                  |
| `Error`    | `string`        | The error message if the operation failed. Empty on success. |

### PackageInfo

| Field         | Type     | Description                  |
| ------------- | -------- | ---------------------------- |
| `Name`        | `string` | The package name.            |
| `Version`     | `string` | The installed version.       |
| `Description` | `string` | The package description.     |
| `Status`      | `string` | The installation status.     |
| `Size`        | `int64`  | The installed size in bytes. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `package:read` permission.

## Example

Refer to
[`examples/operations/package.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/package.go)
for a complete example.
