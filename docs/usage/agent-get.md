# AgentGet

Retrieves detailed information about a specific agent by hostname. Returns
the same rich metadata as `AgentList` but for a single agent. This is a
convenience wrapper that uses `TaskFunc` to call the agent get API directly.

## Usage

```go
step := o.AgentGet("web-01")
```

## Parameters

| Parameter  | Type     | Description                        |
| ---------- | -------- | ---------------------------------- |
| `hostname` | `string` | Hostname of the agent to retrieve. |

## Result Type

```go
var result orchestrator.AgentResult
err := results.Decode("get-agent", &result)
```

See [AgentList](agent-list.md) for the full `AgentResult` field reference.

## Idempotency

**Read-only.** Never modifies state. Always returns `Changed: false`.

## Permissions

Requires `agent:read` permission.
