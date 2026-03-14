# DockerCreate

Creates a new container on the target host. Optionally starts it immediately
when `AutoStart` is set.

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

| Parameter | Type                     | Description                   |
| --------- | ------------------------ | ----------------------------- |
| `target`  | `string`                 | Target host or routing value. |
| `opts`    | `osapi.DockerCreateOpts` | Container creation options.   |

### DockerCreateOpts Fields

| Field       | Type      | Description                                            |
| ----------- | --------- | ------------------------------------------------------ |
| `Image`     | `string`  | Container image reference (required).                  |
| `Name`      | `string`  | Optional container name.                               |
| `Command`   | `[]string`| Overrides the image's default command.                 |
| `Env`       | `[]string`| Environment variables in `KEY=VALUE` format.           |
| `Ports`     | `[]string`| Port mappings in `host_port:container_port` format.    |
| `Volumes`   | `[]string`| Volume mounts in `host_path:container_path` format.    |
| `AutoStart` | `*bool`   | Start the container after creation (default true).     |

## Result Type

```go
var result osapi.DockerResult
err := results.Decode("docker-create", &result)
```

| Field     | Type     | Description                        |
| --------- | -------- | ---------------------------------- |
| `ID`      | `string` | Container ID.                      |
| `Name`    | `string` | Container name.                    |
| `Image`   | `string` | Image used to create the container.|
| `State`   | `string` | Container state after creation.    |
| `Error`   | `string` | Error if creation failed.          |

## Idempotency

**No.** Creates a new container each time the step runs. Use
`DockerList` to check for an existing container before creating.

## Permissions

Requires `docker:write` permission.

## Example

See
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go)
for a complete working example.
