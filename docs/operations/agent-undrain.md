# AgentUndrain

Undrains an agent, allowing it to accept new jobs again. This reverses the
effect of [AgentDrain](agent-drain.md).

## Usage

```go
step := o.AgentUndrain("web-01")
```

## Parameters

| Parameter  | Type     | Description                          |
| ---------- | -------- | ------------------------------------ |
| `hostname` | `string` | Hostname of the agent to undrain.    |

## Result Type

```go
var result osapi.MessageResponse
err := results.Decode("undrain-agent-1", &result)
```

| Field     | Type     | Description                                    |
| --------- | -------- | ---------------------------------------------- |
| `Message` | `string` | Human-readable message describing the outcome. |

## Idempotency

**Idempotent.** Undraining an already-active agent is a no-op.

## Permissions

Requires `agent:write` permission.

## Example

See
[`examples/operations/agent-drain.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/agent-drain.go)
for a complete working example.
