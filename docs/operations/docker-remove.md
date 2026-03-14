# DockerRemove

Removes a container from the target host.

## Usage

```go
step := o.DockerRemove("web-01", "c1a2b3d4e5f6", &osapi.DockerRemoveParams{
    Force: true,
})
```

## Parameters

| Parameter | Type                      | Description                             |
| --------- | ------------------------- | --------------------------------------- |
| `target`  | `string`                  | Target host or routing value.           |
| `id`      | `string`                  | Container ID or name to remove.         |
| `params`  | `*osapi.DockerRemoveParams`| Optional remove parameters.            |

### DockerRemoveParams Fields

| Field   | Type   | Description                                              |
| ------- | ------ | -------------------------------------------------------- |
| `Force` | `bool` | Force removal of a running container (kill then remove). |

## Result Type

```go
var result osapi.DockerActionResult
err := results.Decode("docker-remove", &result)
```

| Field     | Type     | Description                          |
| --------- | -------- | ------------------------------------ |
| `ID`      | `string` | Container ID.                        |
| `Message` | `string` | Status message from the operation.   |
| `Error`   | `string` | Error if removal failed.             |

## Idempotency

**No.** Returns an error if the container does not exist. Use `OnError`
with `Continue` for pre-cleanup patterns.

## Permissions

Requires `docker:write` permission.

## Example

See
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go)
for a complete working example.
