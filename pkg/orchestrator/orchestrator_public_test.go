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
	"net/http"
	"net/http/httptest"
	"testing"

	sdk "github.com/retr0h/osapi/pkg/sdk/orchestrator"
	"github.com/stretchr/testify/suite"

	"github.com/osapi-io/osapi-orchestrator/pkg/orchestrator"
)

type OrchestratorPublicTestSuite struct {
	suite.Suite

	server *httptest.Server
}

func (s *OrchestratorPublicTestSuite) SetupTest() {
	s.server = httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
}

func (s *OrchestratorPublicTestSuite) TearDownTest() {
	s.server.Close()
}

func (s *OrchestratorPublicTestSuite) TestNew() {
	tests := []struct {
		name         string
		validateFunc func(o *orchestrator.Orchestrator)
	}{
		{
			name: "Returns non-nil orchestrator",
			validateFunc: func(o *orchestrator.Orchestrator) {
				s.NotNil(o)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			o := orchestrator.New(s.server.URL, "test-token")
			tc.validateFunc(o)
		})
	}
}

func (s *OrchestratorPublicTestSuite) TestRun() {
	tests := []struct {
		name         string
		setupFunc    func(o *orchestrator.Orchestrator)
		validateFunc func(report *orchestrator.Report, err error)
	}{
		{
			name:      "Empty plan returns empty report",
			setupFunc: func(_ *orchestrator.Orchestrator) {},
			validateFunc: func(
				report *orchestrator.Report,
				err error,
			) {
				s.Require().NoError(err)
				s.NotNil(report)
				s.Empty(report.Tasks)
			},
		},
		{
			name: "TaskFunc executes and returns result",
			setupFunc: func(o *orchestrator.Orchestrator) {
				o.TaskFunc(
					"custom",
					func(
						_ context.Context,
						_ orchestrator.Results,
					) (*sdk.Result, error) {
						return &sdk.Result{
							Changed: true,
							Data:    map[string]any{"key": "val"},
						}, nil
					},
				)
			},
			validateFunc: func(
				report *orchestrator.Report,
				err error,
			) {
				s.Require().NoError(err)
				s.Len(report.Tasks, 1)
				s.Equal("custom", report.Tasks[0].Name)
				s.True(report.Tasks[0].Changed)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			o := orchestrator.New(s.server.URL, "test-token")
			tc.setupFunc(o)

			report, err := o.Run()
			tc.validateFunc(report, err)
		})
	}
}

func (s *OrchestratorPublicTestSuite) TestTaskFunc() {
	tests := []struct {
		name         string
		validateFunc func(step *orchestrator.Step, called *bool)
	}{
		{
			name: "Returns non-nil step",
			validateFunc: func(
				step *orchestrator.Step,
				_ *bool,
			) {
				s.NotNil(step)
			},
		},
		{
			name: "Executes callback when plan runs",
			validateFunc: func(
				_ *orchestrator.Step,
				called *bool,
			) {
				s.True(*called)
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			o := orchestrator.New(s.server.URL, "test-token")

			var called bool
			step := o.TaskFunc(
				"summarize",
				func(
					_ context.Context,
					_ orchestrator.Results,
				) (*sdk.Result, error) {
					called = true
					return &sdk.Result{Changed: true}, nil
				},
			)

			_, err := o.Run()
			s.Require().NoError(err)

			tc.validateFunc(step, &called)
		})
	}
}

func TestOrchestratorPublicTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(OrchestratorPublicTestSuite))
}
