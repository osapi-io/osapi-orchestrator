# DockerImageRemove

Removes a Docker image from the target host.

## Usage

```go
step := o.DockerImageRemove("_any", "nginx:latest",
    &osapi.DockerImageRemoveParams{Force: true},
)
```

## Parameters

| Parameter   | Type                             | Description                   |
| ----------- | -------------------------------- | ----------------------------- |
| `target`    | `string`                         | Target host or routing value. |
| `imageName` | `string`                         | Image name or ID to remove.   |
| `params`    | `*osapi.DockerImageRemoveParams` | Optional remove parameters.   |

### DockerImageRemoveParams

| Field   | Type   | Description                            |
| ------- | ------ | -------------------------------------- |
| `Force` | `bool` | Force removal even if image is in use. |

## Result Type

```go
var result osapi.DockerActionResult
err := results.Decode("docker-image-remove", &result)
```

| Field     | Type     | Description                        |
| --------- | -------- | ---------------------------------- |
| `ID`      | `string` | Image name or ID.                  |
| `Message` | `string` | Status message from the operation. |
| `Error`   | `string` | Error if removal failed.           |

## Idempotency

**Yes.** Removing an already-absent image is a no-op.

## Permissions

Requires `docker:write` permission.

## Example

See
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go)
for a complete working example.
