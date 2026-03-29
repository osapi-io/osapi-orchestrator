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
		name     string
		tasks    []engine.TaskResult
		contains []string
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
			contains: []string{"4 tasks", "2 changed", "1 unchanged", "1 skipped"},
		},
		{
			name: "all statuses including failed",
			tasks: []engine.TaskResult{
				{Name: "a", Status: engine.StatusChanged, Changed: true},
				{Name: "b", Status: engine.StatusUnchanged},
				{Name: "c", Status: engine.StatusSkipped},
				{Name: "d", Status: engine.StatusFailed},
			},
			contains: []string{"4 tasks", "1 changed", "1 unchanged", "1 skipped", "1 failed"},
		},
		{
			name:     "empty report",
			tasks:    nil,
			contains: []string{"0 tasks"},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			report := engine.Report{Tasks: tt.tasks}
			summary := report.Summary()
			for _, c := range tt.contains {
				s.Contains(summary, c)
			}
		})
	}
}

func (s *ResultPublicTestSuite) TestResultStatusField() {
	tests := []struct {
		name       string
		result     *engine.Result
		wantStatus engine.Status
		wantChange bool
	}{
		{
			name: "changed result carries status",
			result: &engine.Result{
				Changed: true,
				Data:    map[string]any{"hostname": "web-01"},
				Status:  engine.StatusChanged,
			},
			wantStatus: engine.StatusChanged,
			wantChange: true,
		},
		{
			name: "unchanged result carries status",
			result: &engine.Result{
				Changed: false,
				Status:  engine.StatusUnchanged,
			},
			wantStatus: engine.StatusUnchanged,
			wantChange: false,
		},
		{
			name: "failed result carries status",
			result: &engine.Result{
				Changed: false,
				Status:  engine.StatusFailed,
			},
			wantStatus: engine.StatusFailed,
			wantChange: false,
		},
		{
			name: "skipped result carries status",
			result: &engine.Result{
				Changed: false,
				Status:  engine.StatusSkipped,
			},
			wantStatus: engine.StatusSkipped,
			wantChange: false,
		},
		{
			name:       "zero value has empty status",
			result:     &engine.Result{},
			wantStatus: "",
			wantChange: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.Equal(tt.wantStatus, tt.result.Status)
			s.Equal(tt.wantChange, tt.result.Changed)
		})
	}
}

func (s *ResultPublicTestSuite) TestResultHostResults() {
	tests := []struct {
		name       string
		result     *engine.Result
		wantLen    int
		validateFn func(hrs []engine.HostResult)
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
			wantLen: 2,
			validateFn: func(hrs []engine.HostResult) {
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
			wantLen: 0,
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
			wantLen: 1,
			validateFn: func(hrs []engine.HostResult) {
				s.Equal("db-01", hrs[0].Hostname)
				s.Equal("migrated", hrs[0].Data["stdout"])
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.Len(tt.result.HostResults, tt.wantLen)

			if tt.validateFn != nil {
				tt.validateFn(tt.result.HostResults)
			}
		})
	}
}

func (s *ResultPublicTestSuite) TestResultsGet() {
	tests := []struct {
		name       string
		results    engine.Results
		lookupName string
		wantNil    bool
		wantChange bool
	}{
		{
			name: "found",
			results: engine.Results{
				"install": {Changed: true},
			},
			lookupName: "install",
			wantNil:    false,
			wantChange: true,
		},
		{
			name:       "not found",
			results:    engine.Results{},
			lookupName: "missing",
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got := tt.results.Get(tt.lookupName)
			if tt.wantNil {
				s.Nil(got)
			} else {
				s.Require().NotNil(got)
				s.Equal(tt.wantChange, got.Changed)
			}
		})
	}
}
