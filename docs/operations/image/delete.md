# ImageDelete

Remove an image.

Runs docker rmi, with -f when force is set, which removes an image that a
container uses. An image that the host does not have reports changed: false.

## Usage

```go
step := o.ImageDelete("web-01", input)
```

## Parameters

| Parameter | Type                     | Description                                                   |
| --------- | ------------------------ | ------------------------------------------------------------- |
| `target`  | `string`                 | Target host: `_any`, `_all`, a hostname, or a label selector. |
| `input`   | `osapi.ImageDeleteInput` | Refer to `ImageDeleteInput` below.                            |

### ImageDeleteInput

| Field   | Type     | Description                                                                                  |
| ------- | -------- | -------------------------------------------------------------------------------------------- |
| `Image` | `string` | The image of the path: the image, by its reference or its ID. It names the item in the path. |
| `Force` | `bool`   | Remove it also when a container uses it.                                                     |

## Result

Each host answers with one `osapi.ImageResult` row. A list answers with one row
for each item.

```go
var row osapi.ImageResult
err := report.Decode("delete-image", &row)
```

| Field      | Type              | Description                                                                                                                                                                                                              |
| ---------- | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Changed`  | `bool`            | Whether the host changed.                                                                                                                                                                                                |
| `Error`    | `string`          | Error message if the agent failed.                                                                                                                                                                                       |
| `Hostname` | `string`          | Hostname of the agent that produced this result.                                                                                                                                                                         |
| `Image`    | `string`          | The image.                                                                                                                                                                                                               |
| `Status`   | `types.RowStatus` | What happened on the host: ok, the operation ran and succeeded; failed, it ran and did not succeed; skipped, the host's OS family does not support it; timeout, the host did not answer, so nobody knows whether it ran. |

## Permissions

The token must have the `image:write` permission. The step sends
`DELETE /api/node/{hostname}/image/{image}`.

## Example

Refer to
[`examples/operations/docker.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/docker.go).
