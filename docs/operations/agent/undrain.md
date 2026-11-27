# AgentUndrain

Undrains an agent. The agent then accepts new jobs again. This operation
reverses [AgentDrain](drain.md).

## Usage

```go
step := o.AgentUndrain("web-01")
```

## Parameters

| Parameter  | Type     | Description                           |
| ---------- | -------- | ------------------------------------- |
| `hostname` | `string` | The hostname of the agent to undrain. |

## Result

The controller answers with no data, so the step has no result to decode. A
refusal fails the step.

## Idempotency

**Idempotent.** If the agent is already active, the operation does nothing.

## Permissions

The token must have the `agent:write` permission.

## Example

Refer to
[`examples/operations/agent-drain.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/agent-drain.go)
for a complete example.
