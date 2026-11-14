# DockerPull

Pulls a Docker image to the target host.

## Usage

```go
step := o.DockerPull("_any", osapi.DockerPullOpts{
    Image: "nginx:latest",
})
```

## Parameters

| Parameter | Type                   | Description                       |
| --------- | ---------------------- | --------------------------------- |
| `target`  | `string`               | The target host or routing value. |
| `opts`    | `osapi.DockerPullOpts` | The pull options (Image).         |

## Result type

```go
var result osapi.DockerPullResult
err := results.Decode("docker-pull", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `ImageID` | `string` | The digest ID of the pulled image.                           |
| `Tag`     | `string` | The tag of the pulled image.                                 |
| `Size`    | `int64`  | The size of the image, in bytes.                             |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Non-idempotent.** The operation always pulls the latest version of the image
from the registry. It does this also if the image is already on the host.

## Permissions

The token must have the `docker:write` permission.

## Example

Refer to
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go)
for a complete example.
