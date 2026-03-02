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
	"net/http"
	"net/http/httptest"
	"testing"

	sdk "github.com/osapi-io/osapi-sdk/pkg/orchestrator"
	"github.com/stretchr/testify/suite"
)

type StepTestSuite struct {
	suite.Suite
}

func (s *StepTestSuite) TestWhenGuardCallbackInvoked() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	orch := New(server.URL, "test-token")

	tests := []struct {
		name      string
		guardFunc func(Results) bool
		expected  bool
	}{
		{
			name: "Guard returns true",
			guardFunc: func(
				_ Results,
			) bool {
				return true
			},
			expected: true,
		},
		{
			name: "Guard returns false",
			guardFunc: func(
				_ Results,
			) bool {
				return false
			},
			expected: false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			step := orch.NodeHostnameGet("_any")
			step.When(tc.guardFunc)

			// Access the SDK task's guard and invoke it to
			// exercise the wrapper lambda.
			guard := step.task.Guard()
			s.Require().NotNil(guard)

			got := guard(sdk.Results{})
			s.Equal(tc.expected, got)
		})
	}
}

func TestStepTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(StepTestSuite))
}
