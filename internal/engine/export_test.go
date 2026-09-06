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
	"context"
	"encoding/json"
	"time"
)

// SetJSONUnmarshalFn overrides the jsonUnmarshalFn injectable for testing.
func SetJSONUnmarshalFn(
	fn func([]byte, any) error,
) {
	jsonUnmarshalFn = fn
}

// ResetJSONUnmarshalFn restores the default jsonUnmarshalFn.
func ResetJSONUnmarshalFn() {
	jsonUnmarshalFn = json.Unmarshal
}

// ExportLevelize exposes the private levelize function for testing.
func ExportLevelize(
	tasks []*Task,
) [][]*Task {
	return levelize(tasks)
}

// ExportRunner wraps a private runner for testing.
type ExportRunner struct {
	r *runner
}

// ExportNewRunner creates a new ExportRunner backed by a private runner.
func ExportNewRunner(
	plan *Plan,
) *ExportRunner {
	return &ExportRunner{r: newRunner(plan)}
}

// Run executes the runner.
func (er *ExportRunner) Run(
	ctx context.Context,
) (*Report, error) {
	return er.r.run(ctx)
}

// GetResult returns the result for the named task.
func (er *ExportRunner) GetResult(
	name string,
) *Result {
	return er.r.results.Get(name)
}

// ExportBackoffDelay exposes the private backoffDelay method on ErrorStrategy
// for testing.
func ExportBackoffDelay(
	initial time.Duration,
	maxInterval time.Duration,
	attempt int,
) time.Duration {
	s := ErrorStrategy{
		kind:            "retry",
		retryCount:      3,
		initialInterval: initial,
		maxInterval:     maxInterval,
	}

	return s.backoffDelay(attempt)
}
