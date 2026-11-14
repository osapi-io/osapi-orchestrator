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

## Result type

```go
var result osapi.MessageResponse
err := results.Decode("drain-agent-1", &result)
```

| Field     | Type     | Description                                            |
| --------- | -------- | ------------------------------------------------------ |
| `Message` | `string` | A message that tells the result, for a person to read. |

## Idempotency

**Idempotent.** If the agent is already drained, the operation does nothing.

## Permissions

The token must have the `agent:write` permission.

## Example

Refer to
[`examples/operations/agent-drain.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/agent-drain.go)
for a complete example.
