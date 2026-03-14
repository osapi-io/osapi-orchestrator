# Operations

The orchestrator provides typed constructors for every OSAPI operation. Each
method returns a `*Step` that can be chained with ordering, conditions, and
error handling.

## Operations

| Method                                        | Operation               | Idempotent | Category |
| --------------------------------------------- | ----------------------- | ---------- | -------- |
| [`HealthCheck`](health-check.md)              | Liveness probe          | Read-only  | Health   |
| [`NodeHostnameGet`](node-hostname-get.md)     | `node.hostname.get`     | Read-only  | Node     |
| [`NodeStatusGet`](node-status-get.md)         | `node.status.get`       | Read-only  | Node     |
| [`NodeUptimeGet`](node-uptime-get.md)         | `node.uptime.get`       | Read-only  | Node     |
| [`NodeDiskGet`](node-disk-get.md)             | `node.disk.get`         | Read-only  | Node     |
| [`NodeMemoryGet`](node-memory-get.md)         | `node.memory.get`       | Read-only  | Node     |
| [`NodeLoadGet`](node-load-get.md)             | `node.load.get`         | Read-only  | Node     |
| [`NetworkDNSGet`](network-dns-get.md)         | `network.dns.get`       | Read-only  | Network  |
| [`NetworkDNSUpdate`](network-dns-update.md)   | `network.dns.update`    | Yes        | Network  |
| [`NetworkPingDo`](network-ping-do.md)         | `network.ping.do`       | Read-only  | Network  |
| [`CommandExec`](command-exec.md)              | `command.exec.execute`  | No         | Command  |
| [`CommandShell`](command-shell.md)            | `command.shell.execute` | No         | Command  |
| [`FileDeploy`](file-deploy.md)                | `file.deploy.execute`   | Yes        | File     |
| [`FileStatusGet`](file-status-get.md)         | `file.status.get`       | Read-only  | File     |
| [`FileUpload`](file-upload.md)                | Upload to Object Store  | No         | File     |
| [`FileChanged`](file-changed.md)              | Check file drift        | Read-only  | File     |
| [`AgentList`](agent-list.md)                  | List active agents      | Read-only  | Agent    |
| [`AgentGet`](agent-get.md)                    | Get agent details       | Read-only  | Agent    |
| [`DockerPull`](docker-pull.md)                | `docker.pull`           | No         | Docker   |
| [`DockerCreate`](docker-create.md)            | `docker.create`         | No         | Docker   |
| [`DockerStart`](docker-start.md)              | `docker.start`          | Yes        | Docker   |
| [`DockerStop`](docker-stop.md)                | `docker.stop`           | Yes        | Docker   |
| [`DockerRemove`](docker-remove.md)            | `docker.remove`         | No         | Docker   |
| [`DockerExec`](docker-exec.md)                | `docker.exec`           | No         | Docker   |
| [`DockerInspect`](docker-inspect.md)          | `docker.inspect`        | Read-only  | Docker   |
| [`DockerList`](docker-list.md)                | `docker.list`           | Read-only  | Docker   |
| [`DockerImageRemove`](docker-image-remove.md) | `docker.image.remove`   | Yes        | Docker   |

## Idempotency

- **Read-only** operations never modify state and always return
  `Changed: false`.
- **Idempotent** write operations check current state before mutating and return
  `Changed: true` only if something actually changed.
- **Non-idempotent** operations (command exec/shell) always return
  `Changed: true`. Use guards (`When`, `OnlyIfChanged`) to control when they
  run.
