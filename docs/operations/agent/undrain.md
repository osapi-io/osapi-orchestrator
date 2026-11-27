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

The controller answers with no data, so the step has no result to decode. The
step is changed when the agent takes jobs again. A refusal fails the step.

## Idempotency

**Idempotent.** If the agent takes jobs already, the controller answers with
409, and the step is unchanged, not failed. Thus a retry and a second run of the
plan succeed, and an `OnlyIfChanged` step after it runs only when the state
changed.

## Permissions

The token must have the `agent:write` permission.

## Example

Refer to
[`examples/operations/agent-drain.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/agent-drain.go)
for a complete example.
