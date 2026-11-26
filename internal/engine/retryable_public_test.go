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

package engine_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/osapi-io/osapi-orchestrator/internal/engine"
	osapiclient "github.com/osapi-io/osapi/pkg/sdk/client"
)

type RetryablePublicTestSuite struct {
	suite.Suite
}

func TestRetryablePublicTestSuite(t *testing.T) {
	suite.Run(t, new(RetryablePublicTestSuite))
}

func (s *RetryablePublicTestSuite) TestRetryable() {
	api := func(code int) error {
		return &osapiclient.APIError{StatusCode: code}
	}

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "a 500", err: api(500), want: true},
		{name: "a 503", err: api(503), want: true},
		{name: "a 408", err: api(408), want: true},
		{name: "a 429", err: api(429), want: true},
		{name: "a 400", err: api(400), want: false},
		{
			name: "a 401 auth error",
			err:  &osapiclient.AuthError{APIError: osapiclient.APIError{StatusCode: 401}},
			want: false,
		},
		{
			name: "a 409 conflict",
			err:  &osapiclient.ConflictError{APIError: osapiclient.APIError{StatusCode: 409}},
			want: false,
		},
		{name: "a cancelled context", err: fmt.Errorf("call: %w", context.Canceled), want: false},
		{name: "a deadline of one attempt", err: context.DeadlineExceeded, want: true},
		{name: "a permanent error", err: engine.Permanent(api(503)), want: false},
		{name: "an error with no type", err: errors.New("connection reset"), want: true},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.Equal(tt.want, engine.Retryable(tt.err))
		})
	}
}

func (s *RetryablePublicTestSuite) TestPermanent() {
	tests := []struct {
		name         string
		err          error
		validateFunc func(err error)
	}{
		{
			name: "a nil error stays nil",
			err:  nil,
			validateFunc: func(err error) {
				s.NoError(err)
			},
		},
		{
			name: "it keeps the message and the wrapped error",
			err:  api404(),
			validateFunc: func(err error) {
				s.EqualError(err, api404().Error())

				var nf *osapiclient.NotFoundError
				s.True(errors.As(err, &nf))
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.validateFunc(engine.Permanent(tt.err))
		})
	}
}

func api404() error {
	return &osapiclient.NotFoundError{
		APIError: osapiclient.APIError{StatusCode: 404, Message: "gone"},
	}
}
