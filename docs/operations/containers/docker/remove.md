# DockerRemove

Removes a container from the target host.

## Usage

```go
step := o.DockerRemove("web-01", "c1a2b3d4e5f6", &osapi.DockerRemoveParams{
    Force: true,
})
```

## Parameters

| Parameter | Type                        | Description                                |
| --------- | --------------------------- | ------------------------------------------ |
| `target`  | `string`                    | The target host or routing value.          |
| `id`      | `string`                    | The ID or name of the container to remove. |
| `params`  | `*osapi.DockerRemoveParams` | The parameters for the removal (optional). |

### DockerRemoveParams fields

| Field   | Type   | Description                                           |
| ------- | ------ | ----------------------------------------------------- |
| `Force` | `bool` | Remove a running container (kill it, then remove it). |

## Result type

```go
var result osapi.DockerActionResult
err := results.Decode("docker-remove", &result)
```

| Field     | Type     | Description                                                  |
| --------- | -------- | ------------------------------------------------------------ |
| `ID`      | `string` | The ID of the container.                                     |
| `Message` | `string` | The status message of the operation.                         |
| `Error`   | `string` | The error message if the operation failed. Empty on success. |

## Idempotency

**Non-idempotent.** If the container does not exist, the operation returns an
error. To remove a container before a plan starts, use `OnError` with
`Continue`.

## Permissions

The token must have the `docker:write` permission.

## Example

Refer to
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go)
for a complete example.
