# DockerStop

Stops a running container on the target host.

## Usage

```go
step := o.DockerStop("web-01", "c1a2b3d4e5f6", osapi.DockerStopOpts{
    Timeout: 30,
})
```

## Parameters

| Parameter | Type                   | Description                              |
| --------- | ---------------------- | ---------------------------------------- |
| `target`  | `string`               | The target host or routing value.        |
| `id`      | `string`               | The ID or name of the container to stop. |
| `opts`    | `osapi.DockerStopOpts` | The stop options.                        |

### DockerStopOpts fields

| Field     | Type  | Description                                                 |
| --------- | ----- | ----------------------------------------------------------- |
| `Timeout` | `int` | The seconds to wait before the kill. Zero uses the default. |

## Result type

```go
var result osapi.DockerActionResult
err := results.Decode("docker-stop", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `ID`      | `string` | The ID of the container.                                     |
| `Message` | `string` | The status message of the operation.                         |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Idempotent.** If the container is already stopped, the operation does nothing
and returns `Changed: false`.

## Permissions

The token must have the `docker:write` permission.

## Example

Refer to
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go)
for a complete example.
