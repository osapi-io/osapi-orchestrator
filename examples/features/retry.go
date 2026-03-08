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

// Package main demonstrates automatic retry on failure. Each TaskFunc
// simulates transient errors that succeed after a few attempts, showing
// three retry strategies:
//   - immediate-retry: fails twice, retries immediately (no backoff)
//   - default-backoff: fails twice, retries with default exponential backoff
//   - custom-backoff: fails twice, retries with custom backoff intervals
//
// DAG:
//
//	immediate-retry  [retry:3, fails 2x]
//	default-backoff  [retry:3, backoff:1s-30s, fails 2x]
//	custom-backoff   [retry:5, backoff:500ms-5s, fails 2x]
//
// Run with: go run retry.go
package main

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	sdk "github.com/retr0h/osapi/pkg/sdk/orchestrator"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
)

// failNTimes returns a TaskFunc that fails the first n calls with a
// transient error, then succeeds.
func failNTimes(n int32) func(context.Context, orchestrator.Results) (*sdk.Result, error) {
	var calls atomic.Int32

	return func(_ context.Context, _ orchestrator.Results) (*sdk.Result, error) {
		attempt := calls.Add(1)
		if attempt <= n {
			return nil, fmt.Errorf("transient failure (attempt %d/%d)", attempt, n)
		}

		return &sdk.Result{Changed: false}, nil
	}
}

func main() {
	// No server needed — TaskFuncs simulate failures locally.
	o := orchestrator.New("http://localhost:8080", "unused")

	// Retry up to 3 times immediately (no backoff).
	// Fails twice, succeeds on the 3rd attempt.
	o.TaskFunc("immediate-retry", failNTimes(2)).
		Retry(3)

	// Retry with default exponential backoff (1s initial, 30s max).
	// Fails twice, succeeds on the 3rd attempt with ~1s, ~2s delays.
	o.TaskFunc("default-backoff", failNTimes(2)).
		Retry(3, orchestrator.WithExponentialBackoff())

	// Retry with custom backoff (500ms initial, 5s max).
	// Fails twice, succeeds on the 3rd attempt with ~500ms, ~1s delays.
	o.TaskFunc("custom-backoff", failNTimes(2)).
		Retry(5, orchestrator.WithBackoff(500*time.Millisecond, 5*time.Second))

	if _, err := o.Run(); err != nil {
		log.Fatal(err)
	}
}
