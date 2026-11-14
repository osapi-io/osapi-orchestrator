# Agent management

Find and manage OSAPI agents. List the active agents, examine one agent, and use
drain and undrain to control which agents accept jobs.

## Operations

| Method                                 | Description                        | Idempotent |
| -------------------------------------- | ---------------------------------- | ---------- |
| [`AgentList()`](list.md)               | List all the active agents         | Read-only  |
| [`AgentGet(hostname)`](get.md)         | Get the details of one agent       | Read-only  |
| [`AgentDrain(hostname)`](drain.md)     | Drain an agent (it stops new jobs) | Idempotent |
| [`AgentUndrain(hostname)`](undrain.md) | Undrain an agent (it accepts jobs) | Idempotent |

## Permissions

| Operation        | Permission    |
| ---------------- | ------------- |
| Read operations  | `agent:read`  |
| Write operations | `agent:write` |

## Example

Refer to
[`examples/operations/agent-drain.go`](../../../examples/operations/agent-drain.go)
for a complete workflow example that uses all the operations.

The example shows how to:

- List the agents to find the fleet
- Examine the metadata of one agent
- Drain and undrain an agent for a maintenance window
