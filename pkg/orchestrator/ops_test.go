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

	"github.com/stretchr/testify/suite"

	"github.com/osapi-io/osapi-orchestrator/internal/engine"
)

// OpsTestSuite tests what the exported steps cannot reach: a row or a
// response that JSON cannot encode, which no SDK type is.
type OpsTestSuite struct {
	suite.Suite
}

func (s *OpsTestSuite) TestHostResultOf() {
	tests := []struct {
		name         string
		row          any
		validateFunc func(engine.HostResult, error)
	}{
		{
			name: "reads the host fields of a row",
			row: map[string]any{
				"hostname": "web-01",
				"status":   "ok",
				"changed":  true,
				"error":    "",
			},
			validateFunc: func(hr engine.HostResult, err error) {
				s.Require().NoError(err)
				s.Equal("web-01", hr.Hostname)
				s.Equal("ok", hr.Status)
				s.True(hr.Changed)
			},
		},
		{
			name: "returns the error of a row that JSON cannot encode",
			row:  make(chan int),
			validateFunc: func(hr engine.HostResult, err error) {
				s.ErrorContains(err, "encode chan int:")
				s.Empty(hr.Hostname)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			tc.validateFunc(hostResultOf(tc.row))
		})
	}
}

func (s *OpsTestSuite) TestDataResult() {
	tests := []struct {
		name         string
		changed      bool
		v            any
		validateFunc func(*engine.Result, error)
	}{
		{
			name:    "keeps changed and the fields of the value",
			changed: true,
			v:       map[string]any{"name": "app.conf"},
			validateFunc: func(r *engine.Result, err error) {
				s.Require().NoError(err)
				s.True(r.Changed)
				s.Equal("app.conf", r.Data["name"])
			},
		},
		{
			name: "returns the error of a value that JSON cannot encode",
			v:    make(chan int),
			validateFunc: func(r *engine.Result, err error) {
				s.ErrorContains(err, "result data: encode chan int:")
				s.Nil(r)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			tc.validateFunc(dataResult(tc.changed, tc.v))
		})
	}
}

func TestOpsTestSuite(t *testing.T) {
	suite.Run(t, new(OpsTestSuite))
}
