# DockerCreate

Creates a new container on the target host. If you set `AutoStart`, the
operation also starts the container immediately.

## Usage

```go
autoStart := true
step := o.DockerCreate("_any", osapi.DockerCreateOpts{
    Image:     "nginx:latest",
    Name:      "web",
    Ports:     []string{"8080:80"},
    AutoStart: &autoStart,
})
```

## Parameters

| Parameter | Type                     | Description                        |
| --------- | ------------------------ | ---------------------------------- |
| `target`  | `string`                 | The target host or routing value.  |
| `opts`    | `osapi.DockerCreateOpts` | The options for the new container. |

### DockerCreateOpts fields

| Field       | Type       | Description                                                        |
| ----------- | ---------- | ------------------------------------------------------------------ |
| `Image`     | `string`   | The reference of the container image (required).                   |
| `Name`      | `string`   | The name of the container (optional).                              |
| `Command`   | `[]string` | A command to use in place of the default command of the image.     |
| `Env`       | `[]string` | The environment variables, in `KEY=VALUE` format.                  |
| `Ports`     | `[]string` | The port mappings, in `host_port:container_port` format.           |
| `Volumes`   | `[]string` | The volume mounts, in `host_path:container_path` format.           |
| `AutoStart` | `*bool`    | Start the container after the operation creates it (default true). |

## Result type

```go
var result osapi.DockerResult
err := results.Decode("docker-create", &result)
```

| Field   | Type     | Description                                                  |
| ------- | -------- | ------------------------------------------------------------ |
| `ID`    | `string` | The ID of the container.                                     |
| `Name`  | `string` | The name of the container.                                   |
| `Image` | `string` | The image that the operation used to create the container.   |
| `State` | `string` | The state of the container after the operation created it.   |
| `Error` | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Non-idempotent.** The operation creates a new container each time the step
runs. Before you create a container, use `DockerList` to find a container that
exists.

## Permissions

The token must have the `docker:write` permission.

## Example

Refer to
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go)
for a complete example.
