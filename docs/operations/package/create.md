# PackageCreate

Install a package.

Installs the package with apt-get install, at the version when one is given. A
package that is installed, at that version when one is given, does not change
and reports changed: false. It runs on the Debian family only, and elsewhere
reports skipped.

## Usage

```go
step := o.PackageCreate("web-01", input)
```

## Parameters

| Parameter | Type                       | Description                                                   |
| --------- | -------------------------- | ------------------------------------------------------------- |
| `target`  | `string`                   | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.PackageCreateInput` | Refer to `PackageCreateInput` below.                          |

### PackageCreateInput

| Field     | Type     | Description                                                             |
| --------- | -------- | ----------------------------------------------------------------------- |
| `Name`    | `string` | The package to install, such as curl.                                   |
| `Version` | `string` | The version to install, such as 8.5.0-2ubuntu10; the newest when empty. |

## Result

Each host answers with one `osapi.PackageResult` row. A list answers with one
row for each item.

```go
var row osapi.PackageResult
err := report.Decode("create-package", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`  | `bool`            | Whether the host changed.                                                                                                                                                                                                |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Name`     | `string`          | The package.                                                                                                                                                                                                             |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `package:write` permission. The step sends
`POST /api/node/{hostname}/package`.

## Example

Refer to
[`examples/operations/package.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/package.go).
