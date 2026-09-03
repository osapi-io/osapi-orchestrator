package engine

import (
	"encoding/json"
	"fmt"

	"github.com/osapi-io/osapi/pkg/sdk/client"
)

// jsonUnmarshalFn is the JSON unmarshal function (injectable for testing).
var jsonUnmarshalFn = json.Unmarshal

// StructToMap converts a struct to map[string]any using its JSON tags.
// Returns nil if v is nil or cannot be marshaled.
func StructToMap(
	v any,
) map[string]any {
	if v == nil {
		return nil
	}

	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}

	var m map[string]any
	if err := jsonUnmarshalFn(b, &m); err != nil {
		return nil
	}

	return m
}

// CollectionResult builds a Result from a Collection response.
// It iterates all results, applies the toHostResult mapper to build
// per-host details, and auto-populates HostResult.Data via StructToMap
// when the mapper leaves it nil. Changed is true if any host reported
// a change.
//
// When rawJSON is non-nil, it is unmarshaled into Result.Data to
// provide the full response for downstream consumers (e.g., guards
// or Results.Decode). Pass resp.RawJSON() for this, or nil to skip.
// collectHosts converts each item to a host result, filling in its data
// from the item itself where the converter left it empty, and reports
// whether any host changed.
func collectHosts[T any](
	items []T,
	toHostResult func(T) HostResult,
) ([]HostResult, bool) {
	hostResults := make([]HostResult, 0, len(items))
	changed := false

	for _, r := range items {
		hr := toHostResult(r)

		if hr.Data == nil {
			hr.Data = StructToMap(r)
		}

		if hr.Changed {
			changed = true
		}

		hostResults = append(hostResults, hr)
	}

	return hostResults, changed
}

// decodeRaw unmarshals the aggregate response body, if there was one.
func decodeRaw(rawJSON []byte) (map[string]any, error) {
	if len(rawJSON) == 0 {
		return nil, nil
	}

	var data map[string]any
	if err := jsonUnmarshalFn(rawJSON, &data); err != nil {
		return nil, fmt.Errorf("unmarshal response data: %w", err)
	}

	return data, nil
}

func CollectionResult[T any](
	col client.Collection[T],
	rawJSON []byte,
	toHostResult func(T) HostResult,
) (*Result, error) {
	hostResults, changed := collectHosts(col.Results, toHostResult)

	data, err := decodeRaw(rawJSON)
	if err != nil {
		return nil, err
	}

	return &Result{
		JobID:       col.JobID,
		Changed:     changed,
		Data:        data,
		HostResults: hostResults,
	}, nil
}
