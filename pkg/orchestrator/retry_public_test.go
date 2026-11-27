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
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
	osapi "github.com/osapi-io/osapi/pkg/sdk/client"
)

type RetryPublicTestSuite struct {
	suite.Suite
}

func TestRetryPublicTestSuite(t *testing.T) {
	suite.Run(t, new(RetryPublicTestSuite))
}

func (s *RetryPublicTestSuite) TestRetryOfAnUnsafeOp() {
	tests := []struct {
		name         string
		status       int
		build        func(o *orchestrator.Orchestrator)
		expectedRuns int32
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
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(
				w http.ResponseWriter,
				_ *http.Request,
			) {
				calls.Add(1)
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
