# DockerStop

Stops a running container on the target host.

## Usage

```go
step := o.DockerStop("web-01", "c1a2b3d4e5f6", osapi.DockerStopOpts{
    Timeout: 30,
})
```

## Parameters

| Parameter | Type                   | Description                             |
| --------- | ---------------------- | --------------------------------------- |
| `target`  | `string`               | Target host or routing value.           |
| `id`      | `string`               | Container ID or name to stop.           |
| `opts`    | `osapi.DockerStopOpts` | Stop options.                           |

### DockerStopOpts Fields

| Field     | Type  | Description                                              |
| --------- | ----- | -------------------------------------------------------- |
| `Timeout` | `int` | Seconds to wait before killing. Zero uses the default.   |

## Result Type

```go
var result osapi.DockerActionResult
err := results.Decode("docker-stop", &result)
```

| Field     | Type     | Description                          |
| --------- | -------- | ------------------------------------ |
| `ID`      | `string` | Container ID.                        |
| `Message` | `string` | Status message from the operation.   |
| `Error`   | `string` | Error if the stop failed.            |

## Idempotency

**Yes.** Stopping an already-stopped container is a no-op and returns
`Changed: false`.

## Permissions

Requires `docker:write` permission.

## Example

See
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go)
for a complete working example.
