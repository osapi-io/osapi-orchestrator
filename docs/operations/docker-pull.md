# DockerPull

Pulls a Docker image on the target host.

## Usage

```go
step := o.DockerPull("_any", osapi.DockerPullOpts{
    Image: "nginx:latest",
})
```

## Parameters

| Parameter | Type                   | Description                   |
| --------- | ---------------------- | ----------------------------- |
| `target`  | `string`               | Target host or routing value. |
| `opts`    | `osapi.DockerPullOpts` | Pull options (Image).         |

## Result Type

```go
var result osapi.DockerPullResult
err := results.Decode("docker-pull", &result)
```

| Field     | Type     | Description             |
| --------- | -------- | ----------------------- |
| `ImageID` | `string` | Pulled image digest ID. |
| `Tag`     | `string` | Image tag pulled.       |
| `Size`    | `int64`  | Image size in bytes.    |
| `Error`   | `string` | Error if pull failed.   |

## Idempotency

**No.** Always pulls the latest version of the image from the registry,
even if the image already exists locally.

## Permissions

Requires `docker:write` permission.

## Example

See
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go)
for a complete working example.
