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
	"os"
	"testing"
	"time"
)

// TestMain makes the default backoff of Retry short for each test of the
// package, so that a test of a retry does not wait seconds. A test of the
// defaults sets them back with SetRetryDefaultsForTest.
func TestMain(m *testing.M) {
	defaultRetryInitial, defaultRetryMax = time.Millisecond, 5*time.Millisecond

	os.Exit(m.Run())
}

// SetRetryDefaultsForTest makes the default backoff of Retry short, so
// that a test of a retry does not wait seconds. It returns a func that
// restores the defaults.
func SetRetryDefaultsForTest(
	initial time.Duration,
	maxInterval time.Duration,
) func() {
	oldInitial, oldMax := defaultRetryInitial, defaultRetryMax
	defaultRetryInitial, defaultRetryMax = initial, maxInterval

	return func() {
		defaultRetryInitial, defaultRetryMax = oldInitial, oldMax
	}
}
