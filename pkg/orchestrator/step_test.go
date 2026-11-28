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
	"time"

	osapi "github.com/osapi-io/osapi/pkg/sdk/client"
	"github.com/stretchr/testify/suite"

	"github.com/osapi-io/osapi-orchestrator/internal/engine"
)

type StepTestSuite struct {
	suite.Suite
}

func (s *StepTestSuite) TestWhen() {
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
			step := orch.HostnameGet("_any")
			step.When(tc.guardFunc)

			// Access the SDK task's guard and invoke it to
			// exercise the wrapper lambda.
			guard := evaluate(step)
			s.Require().NotNil(guard)

			got := guard(engine.Results{})
			s.Equal(tc.expected, got)
		})
	}
}

func (s *StepTestSuite) TestNamed() {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	tests := []struct {
		name         string
		customName   string
		expectedName string
	}{
		{
			name:         "Named sets custom name",
			customName:   "custom-name",
			expectedName: "custom-name",
		},
		{
			name:         "Named overrides auto-generated name",
			customName:   "my-hostname-check",
			expectedName: "my-hostname-check",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			orch := New(server.URL, "test-token")
			step := orch.HostnameGet("_any").Named(tc.customName)

			s.Equal(tc.expectedName, step.task.Name())
		})
	}
}

func (s *StepTestSuite) TestOnlyIfFailed() {
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
		name     string
		results  engine.Results
		expected bool
	}{
		{
			name: "Returns true when dependency failed",
			results: engine.Results{
				"dep": &engine.Result{Status: engine.StatusFailed},
			},
			expected: true,
		},
		{
			name: "Returns false when dependency succeeded",
			results: engine.Results{
				"dep": &engine.Result{Status: engine.StatusChanged},
			},
			expected: false,
		},
		{
			name:     "Returns false when no results",
			results:  engine.Results{},
			expected: false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			dep := orch.HostnameGet("_any").Named("dep")
			step := orch.HostnameGet("_any").After(dep).OnlyIfFailed()

			guard := evaluate(step)
			s.Require().NotNil(guard)
			s.Equal(tc.expected, guard(tc.results))
		})
	}
}

func (s *StepTestSuite) TestOnlyIfAllChanged() {
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
		name     string
		results  engine.Results
		hasDeps  bool
		expected bool
	}{
		{
			name: "Returns true when all deps changed",
			results: engine.Results{
				"dep-a": &engine.Result{Status: engine.StatusChanged, Changed: true},
				"dep-b": &engine.Result{Status: engine.StatusChanged, Changed: true},
			},
			hasDeps:  true,
			expected: true,
		},
		{
			name: "Returns false when one dep unchanged",
			results: engine.Results{
				"dep-a": &engine.Result{Status: engine.StatusChanged, Changed: true},
				"dep-b": &engine.Result{Status: engine.StatusUnchanged},
			},
			hasDeps:  true,
			expected: false,
		},
		{
			name: "Returns true when dep failed but Changed=true (broadcast partial failure)",
			results: engine.Results{
				"dep-a": &engine.Result{Status: engine.StatusFailed, Changed: true},
				"dep-b": &engine.Result{Status: engine.StatusChanged, Changed: true},
			},
			hasDeps:  true,
			expected: true,
		},
		{
			name: "Returns false when dep missing from results",
			results: engine.Results{
				"dep-a": &engine.Result{Status: engine.StatusChanged},
			},
			hasDeps:  true,
			expected: false,
		},
		{
			name:     "Returns false with no dependencies",
			results:  engine.Results{},
			hasDeps:  false,
			expected: false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			var step *Step
			if tc.hasDeps {
				depA := orch.HostnameGet("_any").Named("dep-a")
				depB := orch.HostnameGet("_any").Named("dep-b")
				step = orch.HostnameGet("_any").
					After(depA, depB).
					OnlyIfAllChanged()
			} else {
				step = orch.HostnameGet("_any").OnlyIfAllChanged()
			}

			guard := evaluate(step)
			s.Require().NotNil(guard)
			s.Equal(tc.expected, guard(tc.results))
		})
	}
}

func (s *StepTestSuite) TestOnlyIfAnyHostFailed() {
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
		name     string
		results  engine.Results
		hasDeps  bool
		expected bool
	}{
		{
			name: "Returns true when one host has Status failed",
			results: engine.Results{
				"dep": &engine.Result{
					Status: engine.StatusFailed,
					HostResults: []engine.HostResult{
						{Hostname: "web-01", Status: "failed", Error: "timeout"},
						{Hostname: "web-02", Status: "ok"},
					},
				},
			},
			hasDeps:  true,
			expected: true,
		},
		{
			name: "Returns false when no hosts have Error",
			results: engine.Results{
				"dep": &engine.Result{
					Status: engine.StatusChanged,
					HostResults: []engine.HostResult{
						{Hostname: "web-01", Changed: true},
						{Hostname: "web-02", Changed: true},
					},
				},
			},
			hasDeps:  true,
			expected: false,
		},
		{
			name:     "Returns false with no dependencies",
			results:  engine.Results{},
			hasDeps:  false,
			expected: false,
		},
		{
			name: "Returns false when dep has no HostResults (unicast)",
			results: engine.Results{
				"dep": &engine.Result{
					Status:  engine.StatusChanged,
					Changed: true,
				},
			},
			hasDeps:  true,
			expected: false,
		},
		{
			name: "Returns false when host is skipped (not failed)",
			results: engine.Results{
				"dep": &engine.Result{
					Status: engine.StatusFailed,
					HostResults: []engine.HostResult{
						{Hostname: "web-01", Status: "skipped", Error: "unsupported"},
						{Hostname: "web-02", Status: "ok"},
					},
				},
			},
			hasDeps:  true,
			expected: false,
		},
		{
			name: "Returns true when host has Status failed",
			results: engine.Results{
				"dep": &engine.Result{
					Status: engine.StatusFailed,
					HostResults: []engine.HostResult{
						{Hostname: "web-01", Status: "failed", Error: "permission denied"},
						{Hostname: "web-02", Status: "ok"},
					},
				},
			},
			hasDeps:  true,
			expected: true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			var step *Step
			if tc.hasDeps {
				dep := orch.HostnameGet("_all").Named("dep")
				step = orch.HostnameGet("_any").
					After(dep).
					OnlyIfAnyHostFailed()
			} else {
				step = orch.HostnameGet("_any").OnlyIfAnyHostFailed()
			}

			guard := evaluate(step)
			s.Require().NotNil(guard)
			s.Equal(tc.expected, guard(tc.results))
		})
	}
}

func (s *StepTestSuite) TestOnlyIfAnyHostSkipped() {
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
		name     string
		results  engine.Results
		hasDeps  bool
		expected bool
	}{
		{
			name:     "Returns false with no dependencies",
			results:  engine.Results{},
			hasDeps:  false,
			expected: false,
		},
		{
			name: "Returns false when dep has no HostResults (unicast)",
			results: engine.Results{
				"dep": &engine.Result{
					Status:  engine.StatusChanged,
					Changed: true,
				},
			},
			hasDeps:  true,
			expected: false,
		},
		{
			name: "Returns false when all hosts are ok",
			results: engine.Results{
				"dep": &engine.Result{
					Status: engine.StatusChanged,
					HostResults: []engine.HostResult{
						{Hostname: "web-01", Status: "ok", Changed: true},
						{Hostname: "web-02", Status: "ok", Changed: true},
					},
				},
			},
			hasDeps:  true,
			expected: false,
		},
		{
			name: "Returns true when any host has Status skipped",
			results: engine.Results{
				"dep": &engine.Result{
					Status: engine.StatusFailed,
					HostResults: []engine.HostResult{
						{Hostname: "web-01", Status: "skipped", Error: "unsupported"},
						{Hostname: "web-02", Status: "ok"},
					},
				},
			},
			hasDeps:  true,
			expected: true,
		},
		{
			name: "Returns false when hosts are failed but not skipped",
			results: engine.Results{
				"dep": &engine.Result{
					Status: engine.StatusFailed,
					HostResults: []engine.HostResult{
						{Hostname: "web-01", Status: "failed", Error: "timeout"},
						{Hostname: "web-02", Status: "ok"},
					},
				},
			},
			hasDeps:  true,
			expected: false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			var step *Step
			if tc.hasDeps {
				dep := orch.HostnameGet("_all").Named("dep")
				step = orch.HostnameGet("_any").
					After(dep).
					OnlyIfAnyHostSkipped()
			} else {
				step = orch.HostnameGet("_any").OnlyIfAnyHostSkipped()
			}

			guard := evaluate(step)
			s.Require().NotNil(guard)
			s.Equal(tc.expected, guard(tc.results))
		})
	}
}

func (s *StepTestSuite) TestOnlyIfAllHostsFailed() {
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
		name     string
		results  engine.Results
		hasDeps  bool
		expected bool
	}{
		{
			name: "Returns true when all hosts have Status failed",
			results: engine.Results{
				"dep": &engine.Result{
					Status: engine.StatusFailed,
					HostResults: []engine.HostResult{
						{Hostname: "web-01", Status: "failed", Error: "timeout"},
						{Hostname: "web-02", Status: "failed", Error: "connection refused"},
					},
				},
			},
			hasDeps:  true,
			expected: true,
		},
		{
			name: "Returns false when one host succeeded",
			results: engine.Results{
				"dep": &engine.Result{
					Status: engine.StatusFailed,
					HostResults: []engine.HostResult{
						{Hostname: "web-01", Status: "failed", Error: "timeout"},
						{Hostname: "web-02", Status: "ok", Changed: true},
					},
				},
			},
			hasDeps:  true,
			expected: false,
		},
		{
			name:     "Returns false with no dependencies",
			results:  engine.Results{},
			hasDeps:  false,
			expected: false,
		},
		{
			name: "Returns false when dep has no HostResults (unicast)",
			results: engine.Results{
				"dep": &engine.Result{
					Status:  engine.StatusFailed,
					Changed: true,
				},
			},
			hasDeps:  true,
			expected: false,
		},
		{
			name: "Returns false when some hosts are skipped not failed",
			results: engine.Results{
				"dep": &engine.Result{
					Status: engine.StatusFailed,
					HostResults: []engine.HostResult{
						{Hostname: "web-01", Status: "failed", Error: "timeout"},
						{Hostname: "web-02", Status: "skipped", Error: "unsupported"},
					},
				},
			},
			hasDeps:  true,
			expected: false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			var step *Step
			if tc.hasDeps {
				dep := orch.HostnameGet("_all").Named("dep")
				step = orch.HostnameGet("_any").
					After(dep).
					OnlyIfAllHostsFailed()
			} else {
				step = orch.HostnameGet("_any").OnlyIfAllHostsFailed()
			}

			guard := evaluate(step)
			s.Require().NotNil(guard)
			s.Equal(tc.expected, guard(tc.results))
		})
	}
}

func (s *StepTestSuite) TestOnlyIfAnyHostChanged() {
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
		name     string
		results  engine.Results
		hasDeps  bool
		expected bool
	}{
		{
			name: "Returns true when one host Changed=true",
			results: engine.Results{
				"dep": &engine.Result{
					Status: engine.StatusChanged,
					HostResults: []engine.HostResult{
						{Hostname: "web-01", Changed: true},
						{Hostname: "web-02", Changed: false},
					},
				},
			},
			hasDeps:  true,
			expected: true,
		},
		{
			name: "Returns false when no hosts Changed",
			results: engine.Results{
				"dep": &engine.Result{
					Status: engine.StatusUnchanged,
					HostResults: []engine.HostResult{
						{Hostname: "web-01", Changed: false},
						{Hostname: "web-02", Changed: false},
					},
				},
			},
			hasDeps:  true,
			expected: false,
		},
		{
			name:     "Returns false with no dependencies",
			results:  engine.Results{},
			hasDeps:  false,
			expected: false,
		},
		{
			name: "Returns false when dep has no HostResults (unicast)",
			results: engine.Results{
				"dep": &engine.Result{
					Status:  engine.StatusChanged,
					Changed: true,
				},
			},
			hasDeps:  true,
			expected: false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			var step *Step
			if tc.hasDeps {
				dep := orch.HostnameGet("_all").Named("dep")
				step = orch.HostnameGet("_any").
					After(dep).
					OnlyIfAnyHostChanged()
			} else {
				step = orch.HostnameGet("_any").OnlyIfAnyHostChanged()
			}

			guard := evaluate(step)
			s.Require().NotNil(guard)
			s.Equal(tc.expected, guard(tc.results))
		})
	}
}

func (s *StepTestSuite) TestOnlyIfAllHostsChanged() {
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
		name     string
		results  engine.Results
		hasDeps  bool
		expected bool
	}{
		{
			name: "Returns true when all hosts Changed=true",
			results: engine.Results{
				"dep": &engine.Result{
					Status: engine.StatusChanged,
					HostResults: []engine.HostResult{
						{Hostname: "web-01", Changed: true},
						{Hostname: "web-02", Changed: true},
					},
				},
			},
			hasDeps:  true,
			expected: true,
		},
		{
			name: "Returns false when one host Changed=false",
			results: engine.Results{
				"dep": &engine.Result{
					Status: engine.StatusChanged,
					HostResults: []engine.HostResult{
						{Hostname: "web-01", Changed: true},
						{Hostname: "web-02", Changed: false},
					},
				},
			},
			hasDeps:  true,
			expected: false,
		},
		{
			name:     "Returns false with no dependencies",
			results:  engine.Results{},
			hasDeps:  false,
			expected: false,
		},
		{
			name: "Returns false when dep has no HostResults (unicast)",
			results: engine.Results{
				"dep": &engine.Result{
					Status:  engine.StatusChanged,
					Changed: true,
				},
			},
			hasDeps:  true,
			expected: false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			var step *Step
			if tc.hasDeps {
				dep := orch.HostnameGet("_all").Named("dep")
				step = orch.HostnameGet("_any").
					After(dep).
					OnlyIfAllHostsChanged()
			} else {
				step = orch.HostnameGet("_any").OnlyIfAllHostsChanged()
			}

			guard := evaluate(step)
			s.Require().NotNil(guard)
			s.Equal(tc.expected, guard(tc.results))
		})
	}
}

func (s *StepTestSuite) TestWhenFact() {
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
		name       string
		results    engine.Results
		stepName   string
		target     string
		predicate  Predicate
		wantResult bool
	}{
		{
			name:     "Returns false when agent list step not found",
			results:  engine.Results{},
			stepName: "list-agents",
			target:   "web-01",
			predicate: func(_ osapi.AgentInfo) bool {
				return true
			},
			wantResult: false,
		},
		{
			name: "Returns true when target hostname matches and predicate passes",
			results: engine.Results{
				"list-agents": &engine.Result{
					Data: map[string]any{
						"items": []any{
							map[string]any{
								"hostname": "web-01",
								"os_info": map[string]any{
									"distribution": "Ubuntu",
								},
							},
						},
						"total": float64(1),
					},
				},
			},
			stepName: "list-agents",
			target:   "web-01",
			predicate: func(a osapi.AgentInfo) bool {
				return a.OSInfo.Distribution == "Ubuntu"
			},
			wantResult: true,
		},
		{
			name: "Returns false when target hostname matches but predicate fails",
			results: engine.Results{
				"list-agents": &engine.Result{
					Data: map[string]any{
						"items": []any{
							map[string]any{
								"hostname": "web-01",
								"os_info": map[string]any{
									"distribution": "Debian",
								},
							},
						},
						"total": float64(1),
					},
				},
			},
			stepName: "list-agents",
			target:   "web-01",
			predicate: func(a osapi.AgentInfo) bool {
				return a.OSInfo.Distribution == "Ubuntu"
			},
			wantResult: false,
		},
		{
			name: "Returns true when predicate matches any agent regardless of target",
			results: engine.Results{
				"list-agents": &engine.Result{
					Data: map[string]any{
						"items": []any{
							map[string]any{
								"hostname": "web-02",
							},
						},
						"total": float64(1),
					},
				},
			},
			stepName: "list-agents",
			target:   "web-01",
			predicate: func(_ osapi.AgentInfo) bool {
				return true
			},
			wantResult: true,
		},
		{
			name: "Returns true for _all target when any agent matches predicate",
			results: engine.Results{
				"list-agents": &engine.Result{
					Data: map[string]any{
						"items": []any{
							map[string]any{
								"hostname": "web-01",
								"os_info": map[string]any{
									"distribution": "Debian",
								},
							},
							map[string]any{
								"hostname": "web-02",
								"os_info": map[string]any{
									"distribution": "Ubuntu",
								},
							},
						},
						"total": float64(2),
					},
				},
			},
			stepName: "list-agents",
			target:   "_all",
			predicate: func(a osapi.AgentInfo) bool {
				return a.OSInfo.Distribution == "Ubuntu"
			},
			wantResult: true,
		},
		{
			name: "Returns false for _all target when no agent matches predicate",
			results: engine.Results{
				"list-agents": &engine.Result{
					Data: map[string]any{
						"items": []any{
							map[string]any{
								"hostname": "web-01",
								"os_info": map[string]any{
									"distribution": "Debian",
								},
							},
						},
						"total": float64(1),
					},
				},
			},
			stepName: "list-agents",
			target:   "_all",
			predicate: func(a osapi.AgentInfo) bool {
				return a.OSInfo.Distribution == "Ubuntu"
			},
			wantResult: false,
		},
		{
			name: "Returns false when agent list is empty",
			results: engine.Results{
				"list-agents": &engine.Result{
					Data: map[string]any{
						"items": []any{},
						"total": float64(0),
					},
				},
			},
			stepName: "list-agents",
			target:   "_all",
			predicate: func(_ osapi.AgentInfo) bool {
				return true
			},
			wantResult: false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			list := orch.AgentList().Named(tc.stepName)
			step := orch.HostnameGet(tc.target).
				WhenFact(list, tc.predicate)

			s.Contains(step.task.Dependencies(), list.task)

			guard := evaluate(step)
			s.Require().NotNil(guard)
			s.Equal(tc.wantResult, guard(tc.results))
		})
	}
}

func (s *StepTestSuite) TestRetry() {
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
		name               string
		retryCount         int
		opts               []RetryOption
		expectedString     string
		expectedRetryCount int
	}{
		{
			name:               "Retry without options",
			retryCount:         3,
			opts:               nil,
			expectedString:     "retry(3)",
			expectedRetryCount: 3,
		},
		{
			name:               "Retry with a zero backoff",
			retryCount:         3,
			opts:               []RetryOption{WithBackoff(0, 0)},
			expectedString:     "retry(3)",
			expectedRetryCount: 3,
		},
		{
			name:               "Retry with custom backoff",
			retryCount:         5,
			opts:               []RetryOption{WithBackoff(2*time.Second, 30*time.Second)},
			expectedString:     "retry(5)",
			expectedRetryCount: 5,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			step := orch.HostnameGet("_any").Retry(tc.retryCount, tc.opts...)

			policy := step.task.RetryPolicy()
			s.Require().NotNil(policy)
			s.Equal(tc.expectedString, policy.String())
			s.Equal(tc.expectedRetryCount, policy.Count())
			s.Nil(step.task.ErrorStrategy())
		})
	}
}

// evaluate returns the guards of a step as one predicate.
func evaluate(
	step *Step,
) func(engine.Results) bool {
	return func(results engine.Results) bool {
		run, _ := step.task.Evaluate(results)

		return run
	}
}

func TestStepTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(StepTestSuite))
}

func (s *StepTestSuite) TestUnsafeOps() {
	orch := New("http://127.0.0.1:1", "token")

	tests := []struct {
		name string
		step *Step
		want bool
	}{
		{
			name: "CommandExec",
			step: orch.CommandExec("_any", osapi.CommandExecInput{Command: "uptime"}),
			want: true,
		},
		{
			name: "CommandShell",
			step: orch.CommandShell("_any", osapi.CommandShellInput{Command: "uptime"}),
			want: true,
		},
		{
			name: "DockerCreate",
			step: orch.DockerCreate("_any", osapi.DockerCreateInput{}),
			want: true,
		},
		{
			name: "DockerExec",
			step: orch.DockerExec("_any", osapi.DockerExecInput{Name: "web"}),
			want: true,
		},
		{name: "ServiceRestart", step: orch.ServiceRestart("_any", "nginx"), want: true},
		{
			name: "ProcessSignal",
			step: orch.ProcessSignal("_any", osapi.ProcessSignalInput{PID: 1}),
			want: true,
		},
		{name: "PowerReboot", step: orch.PowerReboot("_any", osapi.PowerRebootInput{}), want: true},
		{
			name: "PowerShutdown",
			step: orch.PowerShutdown("_any", osapi.PowerShutdownInput{}),
			want: true,
		},
		{name: "HostnameGet", step: orch.HostnameGet("_any"), want: false},
		{name: "ServiceStart", step: orch.ServiceStart("_any", "nginx"), want: false},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.Equal(tt.want, tt.step.task.Unsafe())
		})
	}
}
