# ImagePull

Pull an image.

Runs docker pull, and answers the ID and the size of the image that the host has
after it. A pull that gets nothing new, because the host had the image of that
reference already, reports changed: false.

## Usage

```go
step := o.ImagePull("web-01", input)
```

## Parameters

| Parameter | Type                   | Description                                                   |
| --------- | ---------------------- | ------------------------------------------------------------- |
| `target`  | `string`               | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.ImagePullInput` | Refer to `ImagePullInput` below.                              |

### ImagePullInput

| Field   | Type     | Description                     |
| ------- | -------- | ------------------------------- |
| `Image` | `string` | The image, such as alpine:3.20. |

## Result

Each host answers with one `osapi.ImagePullResult` row. A list answers with one
row for each item.

```go
var row osapi.ImagePullResult
err := report.Decode("pull-image", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`  | `bool`            | Whether the pull got a new image: false when the host had it already.                                                                                                                                                    |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Image`    | `string`          | The image.                                                                                                                                                                                                               |
| `ImageID`  | `string`          | The ID of the image that the host has now.                                                                                                                                                                               |
| `Size`     | `int64`           | Its size in bytes.                                                                                                                                                                                                       |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `image:execute` permission. The step sends
`POST /api/node/{hostname}/image/pull`.

## Example

Refer to
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go).
