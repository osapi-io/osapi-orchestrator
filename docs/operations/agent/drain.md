# AgentDrain

Drains an agent, preventing it from accepting new jobs. Existing in-flight jobs
continue to completion. Use [AgentUndrain](undrain.md) to resume.

## Usage

```go
step := o.AgentDrain("web-01")
```

## Parameters

| Parameter  | Type     | Description                     |
| ---------- | -------- | ------------------------------- |
| `hostname` | `string` | Hostname of the agent to drain. |

## Result Type

```go
var result osapi.MessageResponse
err := results.Decode("drain-agent-1", &result)
```

| Field     | Type     | Description                                    |
| --------- | -------- | ---------------------------------------------- |
| `Message` | `string` | Human-readable message describing the outcome. |

## Idempotency

**Idempotent.** Draining an already-drained agent is a no-op.

## Permissions

Requires `agent:write` permission.

## Example

See
[`examples/operations/agent-drain.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/agent-drain.go)
for a complete working example.
