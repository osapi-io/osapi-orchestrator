# PackageList

List the packages installed on the host.

Lists each installed package with its version: dpkg on the Debian family, rpm on
the other Linux families, and Homebrew on macOS.

## Usage

```go
step := o.PackageList("web-01")
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |

## Result

Each host answers with one `osapi.PackageInstalled` row. A list answers with one
row for each item.

```go
var row osapi.PackageInstalled
err := report.Decode("list-package", &row)
```

| Field          | Type              | Description                                                                                                                                                                                                              |
| -------------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Architecture` | `string`          | The architecture it is built for, such as amd64.                                                                                                                                                                         |
| `Description`  | `string`          | Its one-line summary.                                                                                                                                                                                                    |
| `Error`        | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname`     | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Manager`      | `string`          | The package database that has it: dpkg, rpm or brew.                                                                                                                                                                     |
| `Name`         | `string`          | The package.                                                                                                                                                                                                             |
| `Size`         | `int64`           | Its installed size in bytes.                                                                                                                                                                                             |
| `Status`       | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |
| `Version`      | `string`          | The installed version.                                                                                                                                                                                                   |

## Permissions

The token must have the `package:read` permission. The step sends
`GET /api/node/{hostname}/package`.

## Example

Refer to
[`examples/operations/package.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/package.go).
