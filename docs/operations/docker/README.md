# Container Management

Manage Docker containers on target nodes -- full lifecycle from image
pull through container creation, execution, and cleanup.

## Operations

| Method | Description | Idempotent |
| ------ | ----------- | ---------- |
| [`DockerPull(target, opts)`](pull.md) | Pull a Docker image | Non-idempotent |
| [`DockerCreate(target, opts)`](create.md) | Create a container | Non-idempotent |
| [`DockerStart(target, name)`](start.md) | Start a stopped container | Idempotent |
| [`DockerStop(target, name)`](stop.md) | Stop a running container | Idempotent |
| [`DockerRemove(target, name)`](remove.md) | Remove a container | Non-idempotent |
| [`DockerExec(target, opts)`](exec.md) | Execute a command in a container | Non-idempotent |
| [`DockerInspect(target, name)`](inspect.md) | Inspect container details | Read-only |
| [`DockerList(target)`](list.md) | List containers on the host | Read-only |
| [`DockerImageRemove(target, opts)`](image-remove.md) | Remove a Docker image | Idempotent |

## Permissions

| Operation | Permission |
| --------- | ---------- |
| Read operations | `docker:read` |
| Write operations | `docker:write` |
| Exec operations | `docker:execute` |

## Example

See
[`examples/operations/docker.go`](../../examples/operations/docker.go)
for a complete workflow example covering all operations.

The example demonstrates:
- Full container lifecycle: pull, create, start, exec, stop, remove
- Image cleanup with DockerImageRemove
- Inspecting container state between operations
