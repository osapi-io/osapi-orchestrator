# File Workflow

The orchestrator provides a complete upload, deploy, verify lifecycle for
managing files across agents.

## Upload and Deploy

```go
upload := o.FileUpload("config.yaml", "raw", data).After(health)
deploy := o.FileDeploy("_any", osapi.FileDeployOpts{
    ObjectName:  "config.yaml",
    Path:        "/etc/myapp/config.yaml",
    ContentType: "raw",
    Mode:        "0644",
}).After(upload)
o.FileStatusGet("_any", "/etc/myapp/config.yaml").After(deploy)
```

## Conditional Upload

Use `FileChanged` with `OnlyIfChanged` to skip uploads when content is
unchanged:

```go
check := o.FileChanged("config.yaml", data).After(health)
upload := o.FileUpload("config.yaml", "raw", data).
    After(check).OnlyIfChanged()
o.FileDeploy("_any", opts).After(upload).OnlyIfChanged()
```

## Force Upload

Bypass the SHA-256 pre-check with `WithForce()`:

```go
o.FileUpload("config.yaml", "raw", data, orchestrator.WithForce())
```

## Example

See
[`examples/operations/file-deploy.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/file-deploy.go)
and
[`examples/operations/file-changed.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/file-changed.go)
for complete working examples.
