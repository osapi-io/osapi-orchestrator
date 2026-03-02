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

	sdk "github.com/osapi-io/osapi-sdk/pkg/orchestrator"
	"github.com/stretchr/testify/suite"
)

type OptionsTestSuite struct {
	suite.Suite
}

func (s *OptionsTestSuite) TestToSDKStrategy() {
	tests := []struct {
		name     string
		input    ErrorStrategy
		expected sdk.ErrorStrategy
	}{
		{
			name:     "StopAll maps to SDK StopAll",
			input:    StopAll,
			expected: sdk.StopAll,
		},
		{
			name:     "Continue maps to SDK Continue",
			input:    Continue,
			expected: sdk.Continue,
		},
		{
			name:     "Unknown defaults to SDK StopAll",
			input:    ErrorStrategy(99),
			expected: sdk.StopAll,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			got := toSDKStrategy(tc.input)
			s.Equal(tc.expected, got)
		})
	}
}

func TestOptionsTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(OptionsTestSuite))
}
