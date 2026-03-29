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
	"testing"
	"time"

	engine "github.com/osapi-io/osapi-orchestrator/internal/engine"
	"github.com/stretchr/testify/suite"
)

type OptionsTestSuite struct {
	suite.Suite
}

func (s *OptionsTestSuite) TestToSDKStrategy() {
	tests := []struct {
		name     string
		input    ErrorStrategy
		expected engine.ErrorStrategy
	}{
		{
			name:     "StopAll maps to SDK StopAll",
			input:    StopAll,
			expected: engine.StopAll,
		},
		{
			name:     "Continue maps to SDK Continue",
			input:    Continue,
			expected: engine.Continue,
		},
		{
			name:     "Unknown defaults to SDK StopAll",
			input:    ErrorStrategy(99),
			expected: engine.StopAll,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			got := toSDKStrategy(tc.input)
			s.Equal(tc.expected, got)
		})
	}
}

func (s *OptionsTestSuite) TestWithExponentialBackoff() {
	tests := []struct {
		name                    string
		expectedInitialInterval time.Duration
		expectedMaxInterval     time.Duration
	}{
		{
			name:                    "Sets sensible defaults",
			expectedInitialInterval: 1 * time.Second,
			expectedMaxInterval:     30 * time.Second,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			cfg := &retryConfig{}
			WithExponentialBackoff()(cfg)

			s.Equal(tc.expectedInitialInterval, cfg.initialInterval)
			s.Equal(tc.expectedMaxInterval, cfg.maxInterval)
		})
	}
}

func (s *OptionsTestSuite) TestWithBackoff() {
	tests := []struct {
		name                    string
		initial                 time.Duration
		maxInterval             time.Duration
		expectedInitialInterval time.Duration
		expectedMaxInterval     time.Duration
	}{
		{
			name:                    "Sets custom intervals",
			initial:                 2 * time.Second,
			maxInterval:             30 * time.Second,
			expectedInitialInterval: 2 * time.Second,
			expectedMaxInterval:     30 * time.Second,
		},
		{
			name:                    "Sets short intervals",
			initial:                 100 * time.Millisecond,
			maxInterval:             5 * time.Second,
			expectedInitialInterval: 100 * time.Millisecond,
			expectedMaxInterval:     5 * time.Second,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			cfg := &retryConfig{}
			WithBackoff(tc.initial, tc.maxInterval)(cfg)

			s.Equal(tc.expectedInitialInterval, cfg.initialInterval)
			s.Equal(tc.expectedMaxInterval, cfg.maxInterval)
		})
	}
}

func TestOptionsTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(OptionsTestSuite))
}
