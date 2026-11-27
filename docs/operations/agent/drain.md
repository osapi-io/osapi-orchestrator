# AgentDrain

Drains an agent. A drained agent does not accept new jobs. The jobs that are in
progress continue until they complete. Use [AgentUndrain](undrain.md) to accept
jobs again.

## Usage

```go
step := o.AgentDrain("web-01")
```

## Parameters

| Parameter  | Type     | Description                         |
| ---------- | -------- | ----------------------------------- |
| `hostname` | `string` | The hostname of the agent to drain. |

## Result

The controller answers with no data, so the step has no result to decode. A
refusal fails the step.

## Idempotency

**Idempotent.** If the agent is already drained, the operation does nothing.

## Permissions

The token must have the `agent:write` permission.

## Example

Refer to
[`examples/operations/agent-drain.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/agent-drain.go)
for a complete example.
