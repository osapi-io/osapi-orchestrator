# DockerInspect

Retrieves detailed information about a specific container on the target host.

## Usage

```go
step := o.DockerInspect("web-01", "c1a2b3d4e5f6")
```

## Parameters

| Parameter | Type     | Description                      |
| --------- | -------- | -------------------------------- |
| `target`  | `string` | Target host or routing value.    |
| `id`      | `string` | Container ID or name to inspect. |

## Result Type

```go
var result osapi.DockerDetailResult
err := results.Decode("docker-inspect", &result)
```

| Field             | Type                | Description                           |
| ----------------- | ------------------- | ------------------------------------- |
| `ID`              | `string`            | Container ID.                         |
| `Name`            | `string`            | Container name.                       |
| `Image`           | `string`            | Image the container was created from. |
| `State`           | `string`            | Current container state.              |
| `Created`         | `string`            | Container creation timestamp.         |
| `Ports`           | `[]string`          | Published port mappings.              |
| `Mounts`          | `[]string`          | Volume mounts.                        |
| `Env`             | `[]string`          | Environment variables.                |
| `NetworkSettings` | `map[string]string` | Network configuration details.        |
| `Health`          | `string`            | Health check status (if configured).  |
| `Error`           | `string`            | Error if inspect failed.              |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `docker:read` permission.

## Example

See
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go)
for a complete working example.
