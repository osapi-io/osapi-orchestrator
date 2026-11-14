# Verbose output

Use `WithVerbose()` to get more output. With this option, the renderer shows
stdout, stderr, and all the response data for each task.

## Usage

```go
o := orchestrator.New(url, token, orchestrator.WithVerbose())
```

## Example

Refer to
[`examples/features/verbose.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/features/verbose.go)
for a complete example.
