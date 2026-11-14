# Container management

Manage Docker containers on target hosts through their full lifecycle. Pull an
image, create a container, run commands in it, and remove it.

## Operations

| Method                                                            | Description                        | Idempotent     |
| ----------------------------------------------------------------- | ---------------------------------- | -------------- |
| [`DockerPull(target, opts)`](pull.md)                             | Pull a Docker image                | Non-idempotent |
| [`DockerCreate(target, opts)`](create.md)                         | Create a container                 | Non-idempotent |
| [`DockerStart(target, id)`](start.md)                             | Start a stopped container          | Idempotent     |
| [`DockerStop(target, id, opts)`](stop.md)                         | Stop a running container           | Idempotent     |
| [`DockerRemove(target, id, params)`](remove.md)                   | Remove a container                 | Non-idempotent |
| [`DockerExec(target, id, opts)`](exec.md)                         | Run a command in a container       | Non-idempotent |
| [`DockerInspect(target, id)`](inspect.md)                         | Examine the details of a container | Read-only      |
| [`DockerList(target, params)`](list.md)                           | List containers on the host        | Read-only      |
| [`DockerImageRemove(target, imageName, params)`](image-remove.md) | Remove a Docker image              | Idempotent     |

## Permissions

| Operation        | Permission       |
| ---------------- | ---------------- |
| Read operations  | `docker:read`    |
| Write operations | `docker:write`   |
| Exec operations  | `docker:execute` |

## Example

Refer to
[`examples/operations/docker.go`](../../../../examples/operations/docker.go) for
a complete workflow example that uses all the operations.

The example shows how to:

- Do the full container lifecycle: pull, create, start, exec, stop, remove
- Remove an image with DockerImageRemove
- Examine the container state between operations
