# AgentList

Lists all active agents registered with the OSAPI server. Returns agent metadata
including hostname, status, architecture, OS info, memory, load averages,
labels, and network interfaces. This is a convenience wrapper that uses
`TaskFunc` to call the agent list API directly.

## Usage

```go
step := o.AgentList()
```

## Parameters

None.

## Result Type

```go
var result orchestrator.AgentListResult
err := results.Decode("list-agents", &result)
```

| Field    | Type            | Description                      |
| -------- | --------------- | -------------------------------- |
| `Agents` | `[]AgentResult` | Slice of active agent summaries. |
| `Total`  | `int`           | Total number of agents.          |

### AgentResult

| Field           | Type                | Description                        |
| --------------- | ------------------- | ---------------------------------- |
| `Hostname`      | `string`            | Agent hostname.                    |
| `Status`        | `string`            | Agent status (e.g., `Ready`).      |
| `Architecture`  | `string`            | CPU architecture (e.g., `amd64`).  |
| `KernelVersion` | `string`            | Kernel version string.             |
| `CPUCount`      | `int`               | Number of CPUs.                    |
| `FQDN`          | `string`            | Fully qualified domain name.       |
| `ServiceMgr`    | `string`            | Service manager (e.g., `systemd`). |
| `PackageMgr`    | `string`            | Package manager (e.g., `apt`).     |
| `Labels`        | `map[string]string` | Agent labels for targeting.        |
| `OSInfo`        | `*AgentOSInfo`      | OS distribution and version.       |
| `Memory`        | `*AgentMemory`      | Memory usage stats.                |
| `LoadAverage`   | `*AgentLoadAverage` | System load averages.              |
| `Interfaces`    | `[]InterfaceResult` | Network interfaces.                |
| `Uptime`        | `string`            | System uptime.                     |

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `agent:read` permission.
