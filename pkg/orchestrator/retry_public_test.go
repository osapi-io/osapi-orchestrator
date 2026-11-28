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

package orchestrator_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	osapi "github.com/osapi-io/osapi/pkg/sdk/client"
	"github.com/stretchr/testify/suite"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
)

type RetryPublicTestSuite struct {
	suite.Suite
}

func TestRetryPublicTestSuite(t *testing.T) {
	suite.Run(t, new(RetryPublicTestSuite))
}

func (s *RetryPublicTestSuite) TestRetry() {
	tests := []struct {
		name         string
		status       int
		build        func(o *orchestrator.Orchestrator)
		expectedRuns int32
		// defaults, when set, are the backoff defaults for the row.
		defaults [2]time.Duration
		// minGaps are the least waits between the calls.
		minGaps []time.Duration
	}{
		{
			name:   "a command that fails with a 503 does not run again",
			status: http.StatusServiceUnavailable,
			build: func(o *orchestrator.Orchestrator) {
				o.CommandExec("_any", osapi.CommandExecInput{Command: "uptime"}).Retry(2)
			},
			expectedRuns: 1,
		},
		{
			name:   "a command that allows it runs again",
			status: http.StatusServiceUnavailable,
			build: func(o *orchestrator.Orchestrator) {
				o.CommandExec("_any", osapi.CommandExecInput{Command: "uptime"}).
					Retry(2).
					RetryUnsafe()
			},
			expectedRuns: 3,
		},
		{
			name:   "a reboot that fails does not run again",
			status: http.StatusServiceUnavailable,
			build: func(o *orchestrator.Orchestrator) {
				o.PowerReboot("_any", osapi.PowerRebootInput{}).Retry(2)
			},
			expectedRuns: 1,
		},
		{
			name:   "a read that fails with a 503 runs again",
			status: http.StatusServiceUnavailable,
			build: func(o *orchestrator.Orchestrator) {
				o.HostnameGet("_any").Retry(2)
			},
			expectedRuns: 3,
		},
		{
			name:   "a read that fails with a 400 does not run again",
			status: http.StatusBadRequest,
			build: func(o *orchestrator.Orchestrator) {
				o.HostnameGet("_any").Retry(2)
			},
			expectedRuns: 1,
		},
		{
			name:         "Retry waits between attempts with no backoff set",
			status:       http.StatusTooManyRequests,
			build:        func(o *orchestrator.Orchestrator) { o.HostnameGet("_any").Retry(2) },
			expectedRuns: 3,
			defaults:     [2]time.Duration{40 * time.Millisecond, time.Second},
			// The jitter keeps at least half of each delay: 20ms, then 40ms.
			minGaps: []time.Duration{20 * time.Millisecond, 40 * time.Millisecond},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.defaults[0] > 0 {
				restore := orchestrator.SetRetryDefaultsForTest(tt.defaults[0], tt.defaults[1])
				defer restore()
			}

			var (
				calls atomic.Int32
				mu    sync.Mutex
				times []time.Time
			)
			server := httptest.NewServer(http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				calls.Add(1)
				mu.Lock()
				times = append(times, time.Now())
				mu.Unlock()
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(`{"error":"no"}`))
			}))
			defer server.Close()

			o := orchestrator.New(server.URL, "token")
			tt.build(o)

			_, err := o.Run(context.Background())

			s.Error(err)
			s.Equal(tt.expectedRuns, calls.Load())

			mu.Lock()
			defer mu.Unlock()
			for i, gap := range tt.minGaps {
				s.GreaterOrEqual(times[i+1].Sub(times[i]), gap)
			}
		})
	}
}

func (s *RetryPublicTestSuite) TestRetryable() {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "an error with no type", err: errors.New("reset"), want: true},
		{name: "a permanent error", err: orchestrator.Permanent(errors.New("bad")), want: false},
		{name: "a cancelled context", err: context.Canceled, want: false},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.Equal(tt.want, orchestrator.Retryable(tt.err))
		})
	}
}

func (s *RetryPublicTestSuite) TestPlanOptions() {
	fails := func(
		_ context.Context,
		_ *osapi.Client,
		_ orchestrator.Results,
	) (*orchestrator.Result, error) {
		return nil, errors.New("boom")
	}

	tests := []struct {
		name         string
		opts         []orchestrator.Option
		build        func(o *orchestrator.Orchestrator, runs *atomic.Int32)
		validateFunc func(report *orchestrator.Report, err error, runs int32)
	}{
		{
			name: "WithOnError(Continue) lets the independent steps run",
			opts: []orchestrator.Option{orchestrator.WithOnError(orchestrator.Continue)},
			build: func(o *orchestrator.Orchestrator, runs *atomic.Int32) {
				o.TaskFunc("fail", fails)
				o.TaskFunc("other", func(
					_ context.Context,
					_ *osapi.Client,
					_ orchestrator.Results,
				) (*orchestrator.Result, error) {
					runs.Add(1)

					return &orchestrator.Result{}, nil
				})
			},
			validateFunc: func(report *orchestrator.Report, _ error, runs int32) {
				s.Require().NotNil(report)
				s.Equal(int32(1), runs)

				for _, tr := range report.Tasks {
					if tr.Name == "fail" {
						s.Equal(orchestrator.HostStatusFailed, string(tr.Status))
					}
				}
			},
		},
		{
			name: "WithRetry retries a step that sets no retry of its own",
			opts: []orchestrator.Option{orchestrator.WithRetry(2)},
			build: func(o *orchestrator.Orchestrator, runs *atomic.Int32) {
				o.TaskFunc("flaky", func(
					_ context.Context,
					_ *osapi.Client,
					_ orchestrator.Results,
				) (*orchestrator.Result, error) {
					if runs.Add(1) < 3 {
						return nil, errors.New("not yet")
					}

					return &orchestrator.Result{Changed: true}, nil
				})
			},
			validateFunc: func(_ *orchestrator.Report, err error, runs int32) {
				s.NoError(err)
				s.Equal(int32(3), runs)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			var runs atomic.Int32

			o := orchestrator.New("http://127.0.0.1:0", "token", tt.opts...)
			tt.build(o, &runs)

			report, err := o.Run(context.Background())
			tt.validateFunc(report, err, runs.Load())
		})
	}
}
