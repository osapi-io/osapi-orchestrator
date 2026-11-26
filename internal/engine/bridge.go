// Copyright (c) 2026 John Dewey

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package engine

import (
	"encoding/json"
	"fmt"
	"strings"

	client "github.com/osapi-io/osapi/pkg/sdk/client"
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
//
// When one or more hosts failed (HostFailed), the step fails: the
// result comes back with an error that names each failed host, and
// the runner keeps the result for the guards that inspect hosts.
func CollectionResult[T any](
	col client.Collection[T],
	rawJSON []byte,
	toHostResult func(T) HostResult,
) (*Result, error) {
	hostResults := make([]HostResult, 0, len(col.Results))
	changed := false

	var failures []string

	for _, r := range col.Results {
		hr := toHostResult(r)

		if hr.Data == nil {
			hr.Data = StructToMap(r)
		}

		if hr.Changed {
			changed = true
		}

		if HostFailed(hr) {
			failures = append(failures, fmt.Sprintf("%s: %s", hr.Hostname, hostError(hr)))
		}

		hostResults = append(hostResults, hr)
	}

	var data map[string]any
	if len(rawJSON) > 0 {
		if err := jsonUnmarshalFn(rawJSON, &data); err != nil {
			return nil, fmt.Errorf("unmarshal response data: %w", err)
		}
	}

	result := &Result{
		JobID:       col.JobID,
		Changed:     changed,
		Data:        data,
		HostResults: hostResults,
	}

	if len(failures) > 0 {
		return result, fmt.Errorf(
			"%d of %d hosts failed: %s",
			len(failures),
			len(hostResults),
			strings.Join(failures, "; "),
		)
	}

	return result, nil
}

// HostFailed reports whether a host failed: its status is failed, or it
// has an error and was not skipped. A skipped host is not a failure.
func HostFailed(
	hr HostResult,
) bool {
	return hr.Status == "failed" || (hr.Error != "" && hr.Status != "skipped")
}

// hostError returns the error of a failed host, or its status when it
// gave none.
func hostError(
	hr HostResult,
) string {
	if hr.Error != "" {
		return hr.Error
	}

	return hr.Status
}
