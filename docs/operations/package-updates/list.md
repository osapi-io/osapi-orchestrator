# PackageUpdatesList

List the packages the index has newer versions of.

Runs apt list --upgradable: each package with its installed version and the
version the index has. Refresh the index first for what the repositories have
now.

## Usage

```go
step := o.PackageUpdatesList("web-01")
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |

## Result

Each host answers with one `osapi.PackageUpdatesAvailable` row. A list answers
with one row for each item.

```go
var row osapi.PackageUpdatesAvailable
err := report.Decode("list-package-updates", &row)
```

| Field            | Type              | Description                                                                                                                                                                                                              |
| ---------------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `CurrentVersion` | `string`          | The installed version.                                                                                                                                                                                                   |
| `Error`          | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname`       | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Name`           | `string`          | The package.                                                                                                                                                                                                             |
| `NewVersion`     | `string`          | The version the index has.                                                                                                                                                                                               |
| `Status`         | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `package:read` permission. The step sends
`GET /api/node/{hostname}/package/updates`.

## Example

Refer to
[`examples/operations/package.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/package.go).
