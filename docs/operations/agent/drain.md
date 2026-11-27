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

The controller answers with no data, so the step has no result to decode. The
step is changed when the agent drains. A refusal fails the step.

## Idempotency

**Idempotent.** If the agent is already drained, the controller answers with
409, and the step is unchanged, not failed. Thus a retry and a second run of the
plan succeed, and an `OnlyIfChanged` step after it runs only when the state
changed.

## Permissions

The token must have the `agent:write` permission.

## Example

Refer to
[`examples/operations/agent-drain.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/agent-drain.go)
for a complete example.
