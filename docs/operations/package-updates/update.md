# PackageUpdatesUpdate

Refresh the package index.

Runs apt-get update. It reports changed: true when the Release file of a
repository under /var/lib/apt/lists changed, and changed: false when each
repository was unchanged.

## Usage

```go
step := o.PackageUpdatesUpdate("web-01")
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |

## Result

Each host answers with one `osapi.PackageUpdatesResult` row. A list answers with
one row for each item.

```go
var row osapi.PackageUpdatesResult
err := report.Decode("update-package-updates", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`  | `bool`            | Whether the refresh got new index data.                                                                                                                                                                                  |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `package:write` permission. The step sends
`PUT /api/node/{hostname}/package/updates`.

## Example

Refer to
[`examples/operations/package.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/package.go).
