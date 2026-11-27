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
	"time"

	"github.com/osapi-io/osapi-orchestrator/internal/engine"
)

// Option configures the Orchestrator.
type Option func(*config)

type config struct {
	verbose  bool
	strategy *ErrorStrategy
	retry    *planRetry
}

// planRetry is the retry that WithRetry sets for each step of the plan.
type planRetry struct {
	count int
	opts  []RetryOption
}

// WithOnError sets the error strategy of each step that does not set
// its own with OnError or ContinueOnError. The default is StopAll.
func WithOnError(
	strategy ErrorStrategy,
) Option {
	return func(c *config) {
		c.strategy = &strategy
	}
}

// WithRetry sets the retry of each step that does not set its own with
// Retry. It backs off as Retry does.
func WithRetry(
	n int,
	opts ...RetryOption,
) Option {
	return func(c *config) {
		c.retry = &planRetry{count: n, opts: opts}
	}
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
) engine.ErrorStrategy {
	switch s {
	case Continue:
		return engine.Continue
	default:
		return engine.StopAll
	}
}

// RetryOption configures the Retry behavior.
type RetryOption func(*retryConfig)

type retryConfig struct {
	initialInterval time.Duration
	maxInterval     time.Duration
}

// defaultRetryInitial and defaultRetryMax are the backoff of a Retry
// with no WithBackoff: 1s, doubling, and never above 30s. They are
// variables so that a test can make them short.
var (
	defaultRetryInitial = time.Second
	defaultRetryMax     = 30 * time.Second
)

// WithBackoff sets the initial and max intervals of the exponential
// backoff between retry attempts. The delay doubles on each attempt,
// clamped to maxInterval. A maxInterval of 0 or less caps the delay at
// 30s, or at the initial interval when that is longer. Without
// WithBackoff, Retry waits 1s, then doubles up to 30s. Each delay is
// spread over its upper half, so that clients that failed together do
// not retry together.
func WithBackoff(
	initial time.Duration,
	maxInterval time.Duration,
) RetryOption {
	return func(c *retryConfig) {
		c.initialInterval = initial
		c.maxInterval = maxInterval
	}
}

// enginePolicy makes the engine's retry policy for n attempts, with the
// default backoff when no option set one.
func enginePolicy(
	n int,
	opts []RetryOption,
) engine.RetryPolicy {
	cfg := &retryConfig{
		initialInterval: defaultRetryInitial,
		maxInterval:     defaultRetryMax,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.initialInterval <= 0 {
		cfg.initialInterval = defaultRetryInitial
	}

	return engine.Retry(
		n,
		engine.WithRetryBackoff(cfg.initialInterval, cfg.maxInterval),
	)
}
