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

package orchestrator

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	sdk "github.com/osapi-io/osapi-sdk/pkg/orchestrator"
)

// Results provides access to completed step results inside When guards.
type Results struct {
	results sdk.Results
}

// NewResults creates a Results from SDK results. Intended for testing.
func NewResults(
	sdkResults sdk.Results,
) Results {
	return Results{results: sdkResults}
}

// TaskStatus represents the outcome of a step for guard inspection.
type TaskStatus int

const (
	// TaskStatusUnknown indicates the step was not found or has not run.
	TaskStatusUnknown TaskStatus = iota
	// TaskStatusChanged indicates the step ran and reported changes.
	TaskStatusChanged
	// TaskStatusUnchanged indicates the step ran with no changes.
	TaskStatusUnchanged
	// TaskStatusSkipped indicates the step was skipped.
	TaskStatusSkipped
	// TaskStatusFailed indicates the step failed.
	TaskStatusFailed
)

// Status returns the terminal status of a completed dependency step.
func (r Results) Status(
	name string,
) TaskStatus {
	result := r.results.Get(name)
	if result == nil {
		return TaskStatusUnknown
	}

	switch result.Status {
	case sdk.StatusChanged:
		return TaskStatusChanged
	case sdk.StatusUnchanged:
		return TaskStatusUnchanged
	case sdk.StatusSkipped:
		return TaskStatusSkipped
	case sdk.StatusFailed:
		return TaskStatusFailed
	default:
		return TaskStatusUnknown
	}
}

// Decode retrieves the result of a named step and decodes it into
// the given typed struct.
func (r Results) Decode(
	name string,
	v any,
) error {
	result := r.results.Get(name)
	if result == nil {
		return fmt.Errorf("no result for %q", name)
	}

	b, err := json.Marshal(result.Data)
	if err != nil {
		return fmt.Errorf("marshal result data: %w", err)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("decode result data: %w", err)
	}

	return nil
}

// Report summarizes plan execution.
type Report struct {
	Tasks    []sdk.TaskResult
	Duration time.Duration
}

// Summary returns a human-readable summary of the plan execution.
func (r *Report) Summary() string {
	var changed, unchanged, skipped, failed int

	for _, t := range r.Tasks {
		switch t.Status {
		case sdk.StatusChanged:
			changed++
		case sdk.StatusUnchanged:
			unchanged++
		case sdk.StatusSkipped:
			skipped++
		case sdk.StatusFailed:
			failed++
		}
	}

	parts := []string{
		fmt.Sprintf("%d tasks", len(r.Tasks)),
	}

	if changed > 0 {
		parts = append(parts, fmt.Sprintf("%d changed", changed))
	}

	if unchanged > 0 {
		parts = append(
			parts,
			fmt.Sprintf("%d unchanged", unchanged),
		)
	}

	if skipped > 0 {
		parts = append(parts, fmt.Sprintf("%d skipped", skipped))
	}

	if failed > 0 {
		parts = append(parts, fmt.Sprintf("%d failed", failed))
	}

	return strings.Join(parts, ", ")
}
