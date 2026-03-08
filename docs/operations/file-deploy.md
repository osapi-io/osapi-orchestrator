# FileDeploy

Deploys a file from the NATS Object Store to the target agent's filesystem. The
object must be uploaded first (see [FileUpload](file-upload.md)). ContentType
controls whether the content is written as-is (`raw`) or rendered through Go
templates (`template`) with agent facts and user-supplied variables.

## Usage

```go
step := o.FileDeploy("web-01", orchestrator.FileDeployOpts{
    ObjectName:  "app-config.yaml",
    Path:        "/etc/myapp/config.yaml",
    ContentType: "raw",
    Mode:        "0644",
    Owner:       "root",
    Group:       "root",
})
```

With template rendering:

```go
step := o.FileDeploy("web-01", orchestrator.FileDeployOpts{
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
| `opts`    | `FileDeployOpts` | Deploy options (see below).                               |

### FileDeployOpts

| Field         | Type             | Required | Description                                          |
| ------------- | ---------------- | -------- | ---------------------------------------------------- |
| `ObjectName`  | `string`         | Yes      | Name of the object in the NATS Object Store.         |
| `Path`        | `string`         | Yes      | Destination path on the target filesystem.           |
| `ContentType` | `string`         | Yes      | `"raw"` for literal content or `"template"` for Go   |
|               |                  |          | template rendering with facts and vars.              |
| `Mode`        | `string`         | No       | File permission mode (e.g., `"0644"`).               |
| `Owner`       | `string`         | No       | File owner user.                                     |
| `Group`       | `string`         | No       | File owner group.                                    |
| `Vars`        | `map[string]any` | No       | Template variables when ContentType is `"template"`. |

## Result Type

```go
var result orchestrator.FileDeployResult
err := results.Decode("file.deploy.execute-1", &result)
```

| Field     | Type     | Description                                    |
| --------- | -------- | ---------------------------------------------- |
| `Changed` | `bool`   | Whether the file was written or already up to  |
|           |          | date.                                          |
| `SHA256`  | `string` | SHA-256 checksum of the deployed file content. |
| `Path`    | `string` | Destination path the file was written to.      |

## Idempotency

**Idempotent.** Compares the SHA-256 checksum of the existing file against the
expected content. If the checksums match, the file is not rewritten and the step
returns `Changed: false`. When the file is missing or differs, it is written and
the step returns `Changed: true`.

## Permissions

Requires `file:write` permission.

## Example

See
[`examples/operations/file-deploy.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/file-deploy.go)
for a complete working example.
