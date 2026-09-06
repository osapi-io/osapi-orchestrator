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
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/osapi-io/osapi-orchestrator/internal/engine"
)

type ResultPublicTestSuite struct {
	suite.Suite
}

func TestResultPublicTestSuite(t *testing.T) {
	suite.Run(t, new(ResultPublicTestSuite))
}

func (s *ResultPublicTestSuite) TestReportSummary() {
	tests := []struct {
		name         string
		tasks        []engine.TaskResult
		validateFunc func(string)
	}{
		{
			name: "mixed results",
			tasks: []engine.TaskResult{
				{
					Name:     "a",
					Status:   engine.StatusChanged,
					Changed:  true,
					Duration: time.Second,
				},
				{
					Name:     "b",
					Status:   engine.StatusUnchanged,
					Changed:  false,
					Duration: 2 * time.Second,
				},
				{Name: "c", Status: engine.StatusSkipped, Changed: false, Duration: 0},
				{
					Name:     "d",
					Status:   engine.StatusChanged,
					Changed:  true,
					Duration: 500 * time.Millisecond,
				},
			},
			validateFunc: func(summary string) {
				s.Contains(summary, "4 tasks")
				s.Contains(summary, "2 changed")
				s.Contains(summary, "1 unchanged")
				s.Contains(summary, "1 skipped")
			},
		},
		{
			name: "all statuses including failed",
			tasks: []engine.TaskResult{
				{Name: "a", Status: engine.StatusChanged, Changed: true},
				{Name: "b", Status: engine.StatusUnchanged},
				{Name: "c", Status: engine.StatusSkipped},
				{Name: "d", Status: engine.StatusFailed},
			},
			validateFunc: func(summary string) {
				s.Contains(summary, "4 tasks")
				s.Contains(summary, "1 changed")
				s.Contains(summary, "1 unchanged")
				s.Contains(summary, "1 skipped")
				s.Contains(summary, "1 failed")
			},
		},
		{
			name:  "empty report",
			tasks: nil,
			validateFunc: func(summary string) {
				s.Contains(summary, "0 tasks")
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			report := engine.Report{Tasks: tt.tasks}

			tt.validateFunc(report.Summary())
		})
	}
}

func (s *ResultPublicTestSuite) TestResultStatusField() {
	tests := []struct {
		name         string
		result       *engine.Result
		validateFunc func(*engine.Result)
	}{
		{
			name: "changed result carries status",
			result: &engine.Result{
				Changed: true,
				Data:    map[string]any{"hostname": "web-01"},
				Status:  engine.StatusChanged,
			},
			validateFunc: func(result *engine.Result) {
				s.Equal(engine.StatusChanged, result.Status)
				s.True(result.Changed)
			},
		},
		{
			name: "unchanged result carries status",
			result: &engine.Result{
				Changed: false,
				Status:  engine.StatusUnchanged,
			},
			validateFunc: func(result *engine.Result) {
				s.Equal(engine.StatusUnchanged, result.Status)
				s.False(result.Changed)
			},
		},
		{
			name: "failed result carries status",
			result: &engine.Result{
				Changed: false,
				Status:  engine.StatusFailed,
			},
			validateFunc: func(result *engine.Result) {
				s.Equal(engine.StatusFailed, result.Status)
				s.False(result.Changed)
			},
		},
		{
			name: "skipped result carries status",
			result: &engine.Result{
				Changed: false,
				Status:  engine.StatusSkipped,
			},
			validateFunc: func(result *engine.Result) {
				s.Equal(engine.StatusSkipped, result.Status)
				s.False(result.Changed)
			},
		},
		{
			name:   "zero value has empty status",
			result: &engine.Result{},
			validateFunc: func(result *engine.Result) {
				s.Empty(result.Status)
				s.False(result.Changed)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.validateFunc(tt.result)
		})
	}
}

func (s *ResultPublicTestSuite) TestResultHostResults() {
	tests := []struct {
		name         string
		result       *engine.Result
		validateFunc func([]engine.HostResult)
	}{
		{
			name: "result with multiple host results",
			result: &engine.Result{
				Changed: true,
				Status:  engine.StatusChanged,
				HostResults: []engine.HostResult{
					{
						Hostname: "web-01",
						Changed:  true,
						Data:     map[string]any{"stdout": "ok"},
					},
					{
						Hostname: "web-02",
						Changed:  false,
						Error:    "connection timeout",
					},
				},
			},
			validateFunc: func(hrs []engine.HostResult) {
				s.Len(hrs, 2)
				s.Equal("web-01", hrs[0].Hostname)
				s.True(hrs[0].Changed)
				s.Equal("web-02", hrs[1].Hostname)
				s.Equal("connection timeout", hrs[1].Error)
			},
		},
		{
			name: "result with no host results",
			result: &engine.Result{
				Changed: false,
				Status:  engine.StatusUnchanged,
			},
			validateFunc: func(hrs []engine.HostResult) {
				s.Empty(hrs)
			},
		},
		{
			name: "host result with data map",
			result: &engine.Result{
				Changed: true,
				Status:  engine.StatusChanged,
				HostResults: []engine.HostResult{
					{
						Hostname: "db-01",
						Changed:  true,
						Data: map[string]any{
							"stdout":    "migrated",
							"exit_code": float64(0),
						},
					},
				},
			},
			validateFunc: func(hrs []engine.HostResult) {
				s.Len(hrs, 1)
				s.Equal("db-01", hrs[0].Hostname)
				s.Equal("migrated", hrs[0].Data["stdout"])
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.validateFunc(tt.result.HostResults)
		})
	}
}

func (s *ResultPublicTestSuite) TestResultsGet() {
	tests := []struct {
		name         string
		results      engine.Results
		lookupName   string
		validateFunc func(*engine.Result)
	}{
		{
			name: "found",
			results: engine.Results{
				"install": {Changed: true},
			},
			lookupName: "install",
			validateFunc: func(got *engine.Result) {
				s.Require().NotNil(got)
				s.True(got.Changed)
			},
		},
		{
			name:       "not found",
			results:    engine.Results{},
			lookupName: "missing",
			validateFunc: func(got *engine.Result) {
				s.Nil(got)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.validateFunc(tt.results.Get(tt.lookupName))
		})
	}
}
