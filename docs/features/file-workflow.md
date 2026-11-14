# File workflow

The orchestrator can upload a file, deploy it to hosts, and verify it. Use these
steps to manage files on many hosts.

## Upload and deploy

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

## Conditional upload

Use `FileChanged` with `OnlyIfChanged` to skip the upload if the content did not
change:

```go
check := o.FileChanged("config.yaml", data).After(health)
upload := o.FileUpload("config.yaml", "raw", data).
    After(check).OnlyIfChanged()
o.FileDeploy("_any", opts).After(upload).OnlyIfChanged()
```

## Force upload

Use `WithForce()` to skip the SHA-256 check before the upload:

```go
o.FileUpload("config.yaml", "raw", data, orchestrator.WithForce())
```

## Example

Refer to
[`examples/operations/file-deploy.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/file-deploy.go)
and
[`examples/operations/file-changed.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/file-changed.go)
for complete examples.
