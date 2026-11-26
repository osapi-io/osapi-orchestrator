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

import engine "github.com/osapi-io/osapi-orchestrator/internal/engine"

// Permanent marks err as an error that a retry cannot fix, so a step
// with Retry does not run again. A TaskFunc returns it for a failure
// that is the same on each attempt, such as a bad input.
func Permanent(
	err error,
) error {
	return engine.Permanent(err)
}

// Retryable reports whether a retry can fix err. An API answer of 5xx,
// 408 or 429 can, and any other API answer cannot. A cancelled context
// and an error that Permanent marked cannot. Any other error, such as a
// lost connection or a timeout, can.
func Retryable(
	err error,
) bool {
	return engine.Retryable(err)
}
