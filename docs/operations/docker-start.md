# DockerStart

Starts a stopped container on the target host.

## Usage

```go
step := o.DockerStart("web-01", "c1a2b3d4e5f6")
```

## Parameters

| Parameter | Type     | Description                             |
| --------- | -------- | --------------------------------------- |
| `target`  | `string` | Target host or routing value.           |
| `id`      | `string` | Container ID or name to start.          |

## Result Type

```go
var result osapi.DockerActionResult
err := results.Decode("docker-start", &result)
```

| Field     | Type     | Description                          |
| --------- | -------- | ------------------------------------ |
| `ID`      | `string` | Container ID.                        |
| `Message` | `string` | Status message from the operation.   |
| `Error`   | `string` | Error if the start failed.           |

## Idempotency

**Yes.** Starting an already-running container is a no-op and returns
`Changed: false`.

## Permissions

Requires `docker:write` permission.

## Example

See
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go)
for a complete working example.
