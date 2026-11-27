# AgentGet

Gets the details of one agent by its hostname. The result has the same metadata
as `AgentList`, but for one agent only. The operation is a wrapper that uses
`TaskFunc` to call the agent get API directly.

## Usage

```go
step := o.AgentGet("web-01")
```

## Parameters

| Parameter  | Type     | Description                       |
| ---------- | -------- | --------------------------------- |
| `hostname` | `string` | The hostname of the agent to get. |

## Result type

```go
var result osapi.AgentInfo
err := results.Decode("get-agent", &result)
```

Refer to [AgentList](list.md) for all the fields of `AgentResult`.

## Idempotency

**Read-only.** The operation does not change the host. It always returns
`Changed: false`.

## Permissions

The token must have the `agent:read` permission.

## Example

Refer to
[`examples/operations/agent-drain.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/agent-drain.go)
for a complete example.
