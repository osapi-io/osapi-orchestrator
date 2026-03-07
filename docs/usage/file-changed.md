# FileChanged

Checks whether local content differs from the version stored in the Object
Store. Computes SHA-256 locally and compares against the stored hash. Pairs with
`OnlyIfChanged` to skip uploads or deploys when content is unchanged.

## Usage

```go
check := o.FileChanged("config.yaml", localContent)
o.FileUpload("config.yaml", "raw", localContent).
    After(check).
    OnlyIfChanged()
```

## Parameters

| Parameter | Type     | Description                      |
| --------- | -------- | -------------------------------- |
| `name`    | `string` | Object name in the Object Store. |
| `data`    | `[]byte` | Local file content to compare.   |

## Result Type

```go
var result orchestrator.FileChangedResult
err := results.Decode("check-file", &result)
```

| Field     | Type     | Description                                    |
| --------- | -------- | ---------------------------------------------- |
| `Name`    | `string` | Object name that was checked.                  |
| `Changed` | `bool`   | Whether the local content differs from stored. |
| `SHA256`  | `string` | SHA-256 hash of the local content.             |

## Idempotency

**Read-only.** Does not modify any state. Returns `Changed: true` if the file
does not exist in the Object Store or if the SHA-256 hash differs.

## Permissions

Requires `file:read` permission.
