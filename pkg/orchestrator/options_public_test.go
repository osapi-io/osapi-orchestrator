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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
	"github.com/stretchr/testify/suite"
)

type OptionsPublicTestSuite struct {
	suite.Suite
}

func (s *OptionsPublicTestSuite) TestWithVerbose() {
	tests := []struct {
		name string
		opts []orchestrator.Option
	}{
		{
			name: "Creates orchestrator without verbose",
			opts: nil,
		},
		{
			name: "Creates orchestrator with verbose",
			opts: []orchestrator.Option{orchestrator.WithVerbose()},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			server := httptest.NewServer(
				http.HandlerFunc(func(
					w http.ResponseWriter,
					_ *http.Request,
				) {
					w.WriteHeader(http.StatusOK)
				}),
			)
			defer server.Close()

			o := orchestrator.New(
				server.URL,
				"test-token",
				tc.opts...,
			)

			s.NotNil(o)
		})
	}
}

func TestOptionsPublicTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(OptionsPublicTestSuite))
}
