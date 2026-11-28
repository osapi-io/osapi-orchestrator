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

	osapi "github.com/osapi-io/osapi/pkg/sdk/client"
)

// jsonUnmarshalFn is the JSON unmarshal function (injectable for testing).
var jsonUnmarshalFn = json.Unmarshal

// StructToMap converts a struct to map[string]any using its JSON tags.
// It returns nil for a nil v, and an error when v cannot be encoded as a
// JSON object.
func StructToMap(
	v any,
) (map[string]any, error) {
	if v == nil {
		return nil, nil
	}

	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("encode %T: %w", v, err)
	}

	var m map[string]any
	if err := jsonUnmarshalFn(b, &m); err != nil {
		return nil, fmt.Errorf("decode %T: %w", v, err)
	}

	return m, nil
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
//
// An error from toHostResult, or from StructToMap, is the error of the
// step, so that a result never comes back with its data missing.
func CollectionResult[T any](
	col osapi.Collection[T],
	rawJSON []byte,
	toHostResult func(T) (HostResult, error),
) (*Result, error) {
	hostResults := make([]HostResult, 0, len(col.Results))
	changed := false

	var failures []string

	for _, r := range col.Results {
		hr, err := toHostResult(r)
		if err != nil {
			return nil, fmt.Errorf("host result: %w", err)
		}

		if hr.Data == nil {
			hr.Data, err = StructToMap(r)
			if err != nil {
				return nil, fmt.Errorf("host result data: %w", err)
			}
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

// HostFailed reports whether a host failed: its status is failed or
// timeout, or it has an error and was not skipped. A host that timed out
// did not answer, so nobody knows whether the operation ran there; it is
// a failure. A skipped host is not a failure.
func HostFailed(
	hr HostResult,
) bool {
	switch hr.Status {
	case "failed", "timeout":
		return true
	case "skipped":
		return false
	}

	return hr.Error != ""
}

// hostError returns the error of a failed host, or its status when it
// gave none.
func hostError(
	hr HostResult,
) string {
	if hr.Error != "" {
		return hr.Error
	}

	if hr.Status == "timeout" {
		return "the host did not answer"
	}

	return hr.Status
}
