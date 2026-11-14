# DockerImageRemove

Removes a Docker image from the target host.

## Usage

```go
step := o.DockerImageRemove("_any", "nginx:latest",
    &osapi.DockerRemoveImageParams{Force: true},
)
```

## Parameters

| Parameter   | Type                             | Description                                |
| ----------- | -------------------------------- | ------------------------------------------ |
| `target`    | `string`                         | The target host or routing value.          |
| `imageName` | `string`                         | The name or ID of the image to remove.     |
| `params`    | `*osapi.DockerRemoveImageParams` | The parameters for the removal (optional). |

### DockerRemoveImageParams

| Field   | Type   | Description                                    |
| ------- | ------ | ---------------------------------------------- |
| `Force` | `bool` | Remove the image, also if a container uses it. |

## Result type

```go
var result osapi.DockerActionResult
err := results.Decode("docker-image-remove", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `ID`      | `string` | The name or ID of the image.                                 |
| `Message` | `string` | The status message of the operation.                         |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** If the image does not exist, the operation does nothing.

## Permissions

The token must have the `docker:write` permission.

## Example

Refer to
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go)
for a complete example.
