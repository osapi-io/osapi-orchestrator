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
	"errors"
	"net/http"

	osapiclient "github.com/osapi-io/osapi/pkg/sdk/client"
)

// permanentError is an error that a retry cannot fix.
type permanentError struct {
	err error
}

// Error returns the message of the wrapped error.
func (e *permanentError) Error() string {
	return e.err.Error()
}

// Unwrap returns the wrapped error.
func (e *permanentError) Unwrap() error {
	return e.err
}

// Permanent marks err as an error that a retry cannot fix, so the
// runner does not run the task again. A TaskFunc returns it for a
// failure that is the same on each attempt, such as a bad input.
func Permanent(
	err error,
) error {
	if err == nil {
		return nil
	}

	return &permanentError{err: err}
}

// Retryable reports whether a retry can fix err. An API answer of 5xx,
// 408 or 429 can; any other API answer (a refusal, a missing item, a
// conflict, a bad input) cannot. A cancelled context and an error that
// Permanent marked cannot. Any other error, such as a lost connection
// or a timeout, can.
func Retryable(
	err error,
) bool {
	var permanent *permanentError
	if errors.As(err, &permanent) || errors.Is(err, context.Canceled) {
		return false
	}

	var api *osapiclient.APIError
	if errors.As(err, &api) {
		return api.StatusCode >= http.StatusInternalServerError ||
			api.StatusCode == http.StatusRequestTimeout ||
			api.StatusCode == http.StatusTooManyRequests
	}

	return true
}
