# Broadcast

Target `_all` or a label selector to send a job to every matching agent. The
renderer automatically shows per-host results.

## Usage

```go
o.NodeHostnameGet("_all")
```

For label-based targeting:

```go
o.CommandExec("group:web", "uptime")
```

## Example

See
[`examples/features/broadcast.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/broadcast.go)
for a complete working example.
