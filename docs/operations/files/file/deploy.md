# FileDeploy

Deploys a file from the NATS Object Store to the filesystem of the target host.
Upload the object first (refer to [FileUpload](upload.md)). ContentType selects
how the agent writes the content. With `raw`, the agent writes the content as it
is. With `template`, the agent renders the content as a Go template with the
agent facts and the variables that you give.

## Usage

```go
step := o.FileDeploy("web-01", osapi.FileDeployOpts{
    ObjectName:  "app-config.yaml",
    Path:        "/etc/myapp/config.yaml",
    ContentType: "raw",
    Mode:        "0644",
    Owner:       "root",
    Group:       "root",
})
```

With a template:

```go
step := o.FileDeploy("web-01", osapi.FileDeployOpts{
    ObjectName:  "nginx.conf.tmpl",
    Path:        "/etc/nginx/nginx.conf",
    ContentType: "template",
    Mode:        "0644",
    Owner:       "root",
    Group:       "root",
    Vars:        map[string]any{"workers": 4},
})
```

## Parameters

| Parameter | Type             | Description                                               |
| --------- | ---------------- | --------------------------------------------------------- |
| `target`  | `string`         | Target host: `_any`, `_all`, hostname, or label selector. |
| `opts`    | `FileDeployOpts` | Deploy options (refer to the table below).                |

### FileDeployOpts

| Field         | Type             | Required | Description                                              |
| ------------- | ---------------- | -------- | -------------------------------------------------------- |
| `ObjectName`  | `string`         | Yes      | The name of the object in the NATS Object Store.         |
| `Path`        | `string`         | Yes      | The destination path on the target filesystem.           |
| `ContentType` | `string`         | Yes      | `"raw"` for literal content, or `"template"` to render a |
|               |                  |          | Go template with facts and vars.                         |
| `Mode`        | `string`         | No       | The file permission mode (for example, `"0644"`).        |
| `Owner`       | `string`         | No       | The user that owns the file.                             |
| `Group`       | `string`         | No       | The group that owns the file.                            |
| `Vars`        | `map[string]any` | No       | The template variables, if ContentType is `"template"`.  |

## Result type

```go
var result osapi.FileDeployResult
err := results.Decode("file.deploy.execute-1", &result)
```

| Field     | Type     | Description                                               |
| --------- | -------- | --------------------------------------------------------- |
| `Changed` | `bool`   | `true` if the operation wrote the file. `false` if        |
|           |          | the file was already current.                             |
| `SHA256`  | `string` | The SHA-256 checksum of the content of the deployed file. |
| `Path`    | `string` | The destination path where the operation wrote the file.  |

## Idempotency

**Idempotent.** The operation compares the SHA-256 checksum of the current file
with the checksum of the expected content. If the checksums are the same, the
operation does not write the file again, and the step returns `Changed: false`.
If the file is missing or different, the operation writes it, and the step
returns `Changed: true`.

## Permissions

The token must have the `file:write` permission.

## Example

Refer to
[`examples/operations/file-deploy.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/file-deploy.go)
for a complete example.
