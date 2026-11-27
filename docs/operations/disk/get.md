# DiskGet

Get the disk usage of the host.

Reads the size, the used space and the free space of each local filesystem, in
bytes: the mounts of /proc/mounts on Linux and of mount on macOS, with their
sizes from df. Local is ext4, xfs and btrfs on Linux, and apfs and hfs on macOS.
Network and special filesystems are left out.

## Usage

```go
step := o.DiskGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, a hostname, or a label selector. |

## Result

Each host answers with one `osapi.DiskEntry` row. A list answers with one row
for each item.

```go
var row osapi.DiskEntry
err := report.Decode("get-disk", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Disks`    | `[]Disk`          | Each local filesystem.                                                                                                                                                                                                   |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `disk:read` permission. The step sends
`GET /api/node/{hostname}/disk`.

## Example

Refer to
[`examples/features/parallel.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/parallel.go).
