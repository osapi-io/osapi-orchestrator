# AgentList

Lists all the active agents in the registry of the OSAPI server. The result has
the metadata of each agent: hostname, status, architecture, OS information,
memory, load averages, labels, and network interfaces. The operation is a
wrapper that uses `TaskFunc` to call the agent list API directly.

## Usage

```go
step := o.AgentList()
```

## Parameters

None.

## Result type

```go
var result osapi.ListAgentsResponse
err := results.Decode("list-agents", &result)
```

| Field    | Type            | Description                     |
| -------- | --------------- | ------------------------------- |
| `Agents` | `[]AgentResult` | A summary of each active agent. |
| `Total`  | `int`           | The total number of agents.     |

### AgentResult

| Field           | Type                | Description                                     |
| --------------- | ------------------- | ----------------------------------------------- |
| `Hostname`      | `string`            | The hostname of the agent.                      |
| `Status`        | `string`            | The status of the agent (for example, `Ready`). |
| `Architecture`  | `string`            | The CPU architecture (for example, `amd64`).    |
| `KernelVersion` | `string`            | The kernel version.                             |
| `CPUCount`      | `int`               | The number of CPUs.                             |
| `FQDN`          | `string`            | The fully qualified domain name.                |
| `ServiceMgr`    | `string`            | The service manager (for example, `systemd`).   |
| `PackageMgr`    | `string`            | The package manager (for example, `apt`).       |
| `Labels`        | `map[string]string` | The labels that select the agent as a target.   |
| `OSInfo`        | `*AgentOSInfo`      | The OS distribution and its version.            |
| `Memory`        | `*AgentMemory`      | The memory use statistics.                      |
| `LoadAverage`   | `*AgentLoadAverage` | The system load averages.                       |
| `Interfaces`    | `[]InterfaceResult` | The network interfaces.                         |
| `Conditions`    | `[]ConditionResult` | The conditions of the host.                     |
| `Uptime`        | `string`            | The system uptime.                              |
| `StartedAt`     | `time.Time`         | The time when the agent started.                |
| `RegisteredAt`  | `time.Time`         | The time when the agent registered.             |
| `State`         | `string`            | The state of the agent (for example, `active`). |
| `Facts`         | `map[string]any`    | The facts of the agent, for templates to use.   |

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `agent:read` permission.

## Example

Refer to
[`examples/operations/agent-drain.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/agent-drain.go)
for a complete example.
