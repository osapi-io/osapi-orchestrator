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

import sdk "github.com/retr0h/osapi/pkg/sdk/orchestrator"

// Option configures the Orchestrator.
type Option func(*config)

type config struct {
	verbose bool
}

// WithVerbose enables verbose output showing stdout, stderr, and
// full response data for all tasks.
func WithVerbose() Option {
	return func(c *config) {
		c.verbose = true
	}
}

// UploadOption configures the FileUpload operation.
type UploadOption func(*uploadConfig)

type uploadConfig struct {
	force bool
}

// WithForce makes FileUpload bypass the SHA-256 pre-check and always
// upload regardless of whether the content has changed. Without this
// option FileUpload is idempotent — the SDK compares digests and
// skips the upload when content is unchanged.
func WithForce() UploadOption {
	return func(c *uploadConfig) { c.force = true }
}

// ErrorStrategy controls what happens when a step fails.
type ErrorStrategy int

const (
	// StopAll halts the entire plan on failure.
	StopAll ErrorStrategy = iota
	// Continue skips dependent steps and continues with the rest.
	Continue
)

// toSDKStrategy maps a porcelain ErrorStrategy to the SDK type.
func toSDKStrategy(
	s ErrorStrategy,
) sdk.ErrorStrategy {
	switch s {
	case Continue:
		return sdk.Continue
	default:
		return sdk.StopAll
	}
}
