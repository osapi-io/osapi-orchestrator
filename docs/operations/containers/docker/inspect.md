# DockerInspect

Gets the details of one container on the target host.

## Usage

```go
step := o.DockerInspect("web-01", "c1a2b3d4e5f6")
```

## Parameters

| Parameter | Type     | Description                                 |
| --------- | -------- | ------------------------------------------- |
| `target`  | `string` | The target host or routing value.           |
| `id`      | `string` | The ID or name of the container to examine. |

## Result type

```go
var result osapi.DockerDetailResult
err := results.Decode("docker-inspect", &result)
```

| Field             | Type                | Description                                                  |
| ----------------- | ------------------- | ------------------------------------------------------------ |
| `ID`              | `string`            | The ID of the container.                                     |
| `Name`            | `string`            | The name of the container.                                   |
| `Image`           | `string`            | The image that Docker used to create the container.          |
| `State`           | `string`            | The current state of the container.                          |
| `Created`         | `string`            | The time when Docker created the container.                  |
| `Ports`           | `[]string`          | The published port mappings.                                 |
| `Mounts`          | `[]string`          | The volume mounts.                                           |
| `Env`             | `[]string`          | The environment variables.                                   |
| `NetworkSettings` | `map[string]string` | The details of the network configuration.                    |
| `Health`          | `string`            | The status of the health check (if the container has one).   |
| `Error`           | `string`            | The error message if the operation failed. Empty on success. |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `docker:read` permission.

## Example

Refer to
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go)
for a complete example.
