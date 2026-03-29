# Agent Management

Discover and manage OSAPI agents -- list active agents, inspect
individual agents, and control job acceptance with drain/undrain.

## Operations

| Method | Description | Idempotent |
| ------ | ----------- | ---------- |
| [`AgentList()`](list.md) | List all active agents | Read-only |
| [`AgentGet(hostname)`](get.md) | Get details for a specific agent | Read-only |
| [`AgentDrain(hostname)`](drain.md) | Drain agent (stop accepting jobs) | Idempotent |
| [`AgentUndrain(hostname)`](undrain.md) | Undrain agent (resume accepting jobs) | Idempotent |

## Permissions

| Operation | Permission |
| --------- | ---------- |
| Read operations | `agent:read` |
| Write operations | `agent:write` |

## Example

See
[`examples/operations/agent-drain.go`](../../examples/operations/agent-drain.go)
for a complete workflow example covering all operations.

The example demonstrates:
- Listing agents to discover the fleet
- Inspecting a specific agent's metadata
- Draining and undraining for maintenance windows
