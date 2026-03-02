# NodeDiskGet

Retrieves disk usage information for all mounted filesystems on the target node.

## Usage

```go
step := o.NodeDiskGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result Type

```go
var result orchestrator.DiskResult
err := results.Decode("node.disk.get-1", &result)
```

### DiskResult

| Field   | Type          | Description                      |
| ------- | ------------- | -------------------------------- |
| `Disks` | `[]DiskUsage` | List of mounted disk partitions. |

### DiskUsage

| Field   | Type     | Description                 |
| ------- | -------- | --------------------------- |
| `Name`  | `string` | Device or mount point name. |
| `Total` | `uint64` | Total capacity in bytes.    |
| `Used`  | `uint64` | Used space in bytes.        |
| `Free`  | `uint64` | Free space in bytes.        |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `node:read` permission.
