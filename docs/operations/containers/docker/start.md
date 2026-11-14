# DockerStart

Starts a stopped container on the target host.

## Usage

```go
step := o.DockerStart("web-01", "c1a2b3d4e5f6")
```

## Parameters

| Parameter | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `target`  | `string` | Target host: `_any`, `_all`, hostname, or label selector. |
| `id`      | `string` | The ID or name of the container to start.                 |

## Result type

```go
var result osapi.DockerActionResult
err := results.Decode("docker-start", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `ID`      | `string` | The ID of the container.                                     |
| `Message` | `string` | The status message of the operation.                         |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** If the container is already in the `running` state, the
operation does nothing and returns `Changed: false`.

## Permissions

The token must have the `docker:write` permission.

## Example

Refer to
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go)
for a complete example.
