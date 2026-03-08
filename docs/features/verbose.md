# Verbose Output

Enable detailed output with `WithVerbose()`. When enabled, the renderer shows
stdout, stderr, and full response data for every task.

## Usage

```go
o := orchestrator.New(url, token, orchestrator.WithVerbose())
```

## Example

See
[`examples/verbose.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/verbose.go)
for a complete working example.
