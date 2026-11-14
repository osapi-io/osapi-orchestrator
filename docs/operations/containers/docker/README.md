# Container management

Manage Docker containers on target hosts through their full lifecycle. Pull an
image, create a container, run commands in it, and remove it.

## Operations

| Method                                               | Description                        | Idempotent     |
| ---------------------------------------------------- | ---------------------------------- | -------------- |
| [`DockerPull(target, opts)`](pull.md)                | Pull a Docker image                | Non-idempotent |
| [`DockerCreate(target, opts)`](create.md)            | Create a container                 | Non-idempotent |
| [`DockerStart(target, name)`](start.md)              | Start a stopped container          | Idempotent     |
| [`DockerStop(target, name)`](stop.md)                | Stop a running container           | Idempotent     |
| [`DockerRemove(target, name)`](remove.md)            | Remove a container                 | Non-idempotent |
| [`DockerExec(target, opts)`](exec.md)                | Run a command in a container       | Non-idempotent |
| [`DockerInspect(target, name)`](inspect.md)          | Examine the details of a container | Read-only      |
| [`DockerList(target)`](list.md)                      | List containers on the host        | Read-only      |
| [`DockerImageRemove(target, opts)`](image-remove.md) | Remove a Docker image              | Idempotent     |

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
