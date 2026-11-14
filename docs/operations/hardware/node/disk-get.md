# NodeDiskGet

Gets the disk use of all the mounted filesystems on the target host.

## Usage

```go
step := o.NodeDiskGet("web-01")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |

## Result type

```go
var result osapi.DiskResult
err := results.Decode("node.disk.get-1", &result)
```

### DiskResult

| Field   | Type          | Description                  |
| ------- | ------------- | ---------------------------- |
| `Disks` | `[]DiskUsage` | The mounted disk partitions. |

### DiskUsage

| Field   | Type     | Description                            |
| ------- | -------- | -------------------------------------- |
| `Name`  | `string` | The name of the device or mount point. |
| `Total` | `uint64` | The total capacity, in bytes.          |
| `Used`  | `uint64` | The used space, in bytes.              |
| `Free`  | `uint64` | The free space, in bytes.              |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `node:read` permission.

## Example

Refer to
[`examples/operations/node-info.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/node-info.go)
for a complete example.
