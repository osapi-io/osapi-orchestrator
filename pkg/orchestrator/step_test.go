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

	osapi "github.com/retr0h/osapi/pkg/sdk/client"
	sdk "github.com/retr0h/osapi/pkg/sdk/orchestrator"
	"github.com/stretchr/testify/suite"
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
			step := orch.NodeHostnameGet("_any").Named(tc.customName)

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
		results  sdk.Results
		expected bool
	}{
		{
			name: "Returns true when dependency failed",
			results: sdk.Results{
				"dep": &sdk.Result{Status: sdk.StatusFailed},
			},
			expected: true,
		},
		{
			name: "Returns false when dependency succeeded",
			results: sdk.Results{
				"dep": &sdk.Result{Status: sdk.StatusChanged},
			},
			expected: false,
		},
		{
			name:     "Returns false when no results",
			results:  sdk.Results{},
			expected: false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			dep := orch.NodeHostnameGet("_any").Named("dep")
			step := orch.NodeHostnameGet("_any").After(dep).OnlyIfFailed()

			guard := step.task.Guard()
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
		results  sdk.Results
		hasDeps  bool
		expected bool
	}{
		{
			name: "Returns true when all deps changed",
			results: sdk.Results{
				"dep-a": &sdk.Result{Status: sdk.StatusChanged, Changed: true},
				"dep-b": &sdk.Result{Status: sdk.StatusChanged, Changed: true},
			},
			hasDeps:  true,
			expected: true,
		},
		{
			name: "Returns false when one dep unchanged",
			results: sdk.Results{
				"dep-a": &sdk.Result{Status: sdk.StatusChanged, Changed: true},
				"dep-b": &sdk.Result{Status: sdk.StatusUnchanged},
			},
			hasDeps:  true,
			expected: false,
		},
		{
			name: "Returns true when dep failed but Changed=true (broadcast partial failure)",
			results: sdk.Results{
				"dep-a": &sdk.Result{Status: sdk.StatusFailed, Changed: true},
				"dep-b": &sdk.Result{Status: sdk.StatusChanged, Changed: true},
			},
			hasDeps:  true,
			expected: true,
		},
		{
			name: "Returns false when dep missing from results",
			results: sdk.Results{
				"dep-a": &sdk.Result{Status: sdk.StatusChanged},
			},
			hasDeps:  true,
			expected: false,
		},
		{
			name:     "Returns false with no dependencies",
			results:  sdk.Results{},
			hasDeps:  false,
			expected: false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			var step *Step
			if tc.hasDeps {
				depA := orch.NodeHostnameGet("_any").Named("dep-a")
				depB := orch.NodeHostnameGet("_any").Named("dep-b")
				step = orch.NodeHostnameGet("_any").
					After(depA, depB).
					OnlyIfAllChanged()
			} else {
				step = orch.NodeHostnameGet("_any").OnlyIfAllChanged()
			}

			guard := step.task.Guard()
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
		results  sdk.Results
		hasDeps  bool
		expected bool
	}{
		{
			name: "Returns true when one host has Status failed",
			results: sdk.Results{
				"dep": &sdk.Result{
					Status: sdk.StatusFailed,
					HostResults: []sdk.HostResult{
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
			results: sdk.Results{
				"dep": &sdk.Result{
					Status: sdk.StatusChanged,
					HostResults: []sdk.HostResult{
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
			results:  sdk.Results{},
			hasDeps:  false,
			expected: false,
		},
		{
			name: "Returns false when dep has no HostResults (unicast)",
			results: sdk.Results{
				"dep": &sdk.Result{
					Status:  sdk.StatusChanged,
					Changed: true,
				},
			},
			hasDeps:  true,
			expected: false,
		},
		{
			name: "Returns false when host is skipped (not failed)",
			results: sdk.Results{
				"dep": &sdk.Result{
					Status: sdk.StatusFailed,
					HostResults: []sdk.HostResult{
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
			results: sdk.Results{
				"dep": &sdk.Result{
					Status: sdk.StatusFailed,
					HostResults: []sdk.HostResult{
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
				dep := orch.NodeHostnameGet("_all").Named("dep")
				step = orch.NodeHostnameGet("_any").
					After(dep).
					OnlyIfAnyHostFailed()
			} else {
				step = orch.NodeHostnameGet("_any").OnlyIfAnyHostFailed()
			}

			guard := step.task.Guard()
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
		results  sdk.Results
		hasDeps  bool
		expected bool
	}{
		{
			name:     "Returns false with no dependencies",
			results:  sdk.Results{},
			hasDeps:  false,
			expected: false,
		},
		{
			name: "Returns false when dep has no HostResults (unicast)",
			results: sdk.Results{
				"dep": &sdk.Result{
					Status:  sdk.StatusChanged,
					Changed: true,
				},
			},
			hasDeps:  true,
			expected: false,
		},
		{
			name: "Returns false when all hosts are ok",
			results: sdk.Results{
				"dep": &sdk.Result{
					Status: sdk.StatusChanged,
					HostResults: []sdk.HostResult{
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
			results: sdk.Results{
				"dep": &sdk.Result{
					Status: sdk.StatusFailed,
					HostResults: []sdk.HostResult{
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
			results: sdk.Results{
				"dep": &sdk.Result{
					Status: sdk.StatusFailed,
					HostResults: []sdk.HostResult{
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
				dep := orch.NodeHostnameGet("_all").Named("dep")
				step = orch.NodeHostnameGet("_any").
					After(dep).
					OnlyIfAnyHostSkipped()
			} else {
				step = orch.NodeHostnameGet("_any").OnlyIfAnyHostSkipped()
			}

			guard := step.task.Guard()
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
		results  sdk.Results
		hasDeps  bool
		expected bool
	}{
		{
			name: "Returns true when all hosts have Status failed",
			results: sdk.Results{
				"dep": &sdk.Result{
					Status: sdk.StatusFailed,
					HostResults: []sdk.HostResult{
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
			results: sdk.Results{
				"dep": &sdk.Result{
					Status: sdk.StatusFailed,
					HostResults: []sdk.HostResult{
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
			results:  sdk.Results{},
			hasDeps:  false,
			expected: false,
		},
		{
			name: "Returns false when dep has no HostResults (unicast)",
			results: sdk.Results{
				"dep": &sdk.Result{
					Status:  sdk.StatusFailed,
					Changed: true,
				},
			},
			hasDeps:  true,
			expected: false,
		},
		{
			name: "Returns false when some hosts are skipped not failed",
			results: sdk.Results{
				"dep": &sdk.Result{
					Status: sdk.StatusFailed,
					HostResults: []sdk.HostResult{
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
				dep := orch.NodeHostnameGet("_all").Named("dep")
				step = orch.NodeHostnameGet("_any").
					After(dep).
					OnlyIfAllHostsFailed()
			} else {
				step = orch.NodeHostnameGet("_any").OnlyIfAllHostsFailed()
			}

			guard := step.task.Guard()
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
		results  sdk.Results
		hasDeps  bool
		expected bool
	}{
		{
			name: "Returns true when one host Changed=true",
			results: sdk.Results{
				"dep": &sdk.Result{
					Status: sdk.StatusChanged,
					HostResults: []sdk.HostResult{
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
			results: sdk.Results{
				"dep": &sdk.Result{
					Status: sdk.StatusUnchanged,
					HostResults: []sdk.HostResult{
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
			results:  sdk.Results{},
			hasDeps:  false,
			expected: false,
		},
		{
			name: "Returns false when dep has no HostResults (unicast)",
			results: sdk.Results{
				"dep": &sdk.Result{
					Status:  sdk.StatusChanged,
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
				dep := orch.NodeHostnameGet("_all").Named("dep")
				step = orch.NodeHostnameGet("_any").
					After(dep).
					OnlyIfAnyHostChanged()
			} else {
				step = orch.NodeHostnameGet("_any").OnlyIfAnyHostChanged()
			}

			guard := step.task.Guard()
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
		results  sdk.Results
		hasDeps  bool
		expected bool
	}{
		{
			name: "Returns true when all hosts Changed=true",
			results: sdk.Results{
				"dep": &sdk.Result{
					Status: sdk.StatusChanged,
					HostResults: []sdk.HostResult{
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
			results: sdk.Results{
				"dep": &sdk.Result{
					Status: sdk.StatusChanged,
					HostResults: []sdk.HostResult{
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
			results:  sdk.Results{},
			hasDeps:  false,
			expected: false,
		},
		{
			name: "Returns false when dep has no HostResults (unicast)",
			results: sdk.Results{
				"dep": &sdk.Result{
					Status:  sdk.StatusChanged,
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
				dep := orch.NodeHostnameGet("_all").Named("dep")
				step = orch.NodeHostnameGet("_any").
					After(dep).
					OnlyIfAllHostsChanged()
			} else {
				step = orch.NodeHostnameGet("_any").OnlyIfAllHostsChanged()
			}

			guard := step.task.Guard()
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
		results    sdk.Results
		stepName   string
		target     string
		predicate  Predicate
		wantResult bool
	}{
		{
			name:     "Returns false when agent list step not found",
			results:  sdk.Results{},
			stepName: "list-agents",
			target:   "web-01",
			predicate: func(_ osapi.Agent) bool {
				return true
			},
			wantResult: false,
		},
		{
			name: "Returns true when target hostname matches and predicate passes",
			results: sdk.Results{
				"list-agents": &sdk.Result{
					Data: map[string]any{
						"agents": []any{
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
			predicate: func(a osapi.Agent) bool {
				return a.OSInfo != nil && a.OSInfo.Distribution == "Ubuntu"
			},
			wantResult: true,
		},
		{
			name: "Returns false when target hostname matches but predicate fails",
			results: sdk.Results{
				"list-agents": &sdk.Result{
					Data: map[string]any{
						"agents": []any{
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
			predicate: func(a osapi.Agent) bool {
				return a.OSInfo != nil && a.OSInfo.Distribution == "Ubuntu"
			},
			wantResult: false,
		},
		{
			name: "Returns true when predicate matches any agent regardless of target",
			results: sdk.Results{
				"list-agents": &sdk.Result{
					Data: map[string]any{
						"agents": []any{
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
			predicate: func(_ osapi.Agent) bool {
				return true
			},
			wantResult: true,
		},
		{
			name: "Returns true for _all target when any agent matches predicate",
			results: sdk.Results{
				"list-agents": &sdk.Result{
					Data: map[string]any{
						"agents": []any{
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
			predicate: func(a osapi.Agent) bool {
				return a.OSInfo != nil && a.OSInfo.Distribution == "Ubuntu"
			},
			wantResult: true,
		},
		{
			name: "Returns false for _all target when no agent matches predicate",
			results: sdk.Results{
				"list-agents": &sdk.Result{
					Data: map[string]any{
						"agents": []any{
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
			predicate: func(a osapi.Agent) bool {
				return a.OSInfo != nil && a.OSInfo.Distribution == "Ubuntu"
			},
			wantResult: false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			step := orch.NodeHostnameGet(tc.target).
				WhenFact(tc.stepName, tc.predicate)

			guard := step.task.Guard()
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
			name:               "Retry without options (backwards compatible)",
			retryCount:         3,
			opts:               nil,
			expectedString:     "retry(3)",
			expectedRetryCount: 3,
		},
		{
			name:               "Retry with default exponential backoff",
			retryCount:         3,
			opts:               []RetryOption{WithExponentialBackoff()},
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
			step := orch.NodeHostnameGet("_any").Retry(tc.retryCount, tc.opts...)

			strategy := step.task.ErrorStrategy()
			s.Require().NotNil(strategy)
			s.Equal(tc.expectedString, strategy.String())
			s.Equal(tc.expectedRetryCount, strategy.RetryCount())
		})
	}
}

func TestStepTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(StepTestSuite))
}
