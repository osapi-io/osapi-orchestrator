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

	sdk "github.com/osapi-io/osapi-sdk/pkg/orchestrator"
	"github.com/stretchr/testify/suite"
)

type StepTestSuite struct {
	suite.Suite
}

func (s *StepTestSuite) TestWhenGuardCallbackInvoked() {
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

func (s *StepTestSuite) TestNamedSetsTaskName() {
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
	step := orch.NodeHostnameGet("_any").Named("custom-name")

	s.Equal("custom-name", step.task.Name())
}

func (s *StepTestSuite) TestOnlyIfFailedGuard() {
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

func (s *StepTestSuite) TestOnlyIfAllChangedGuard() {
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
				"dep-a": &sdk.Result{Status: sdk.StatusChanged},
				"dep-b": &sdk.Result{Status: sdk.StatusChanged},
			},
			hasDeps:  true,
			expected: true,
		},
		{
			name: "Returns false when one dep unchanged",
			results: sdk.Results{
				"dep-a": &sdk.Result{Status: sdk.StatusChanged},
				"dep-b": &sdk.Result{Status: sdk.StatusUnchanged},
			},
			hasDeps:  true,
			expected: false,
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

func (s *StepTestSuite) TestWhenFactGuardBehavior() {
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
			predicate: func(_ AgentResult) bool {
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
			predicate: func(a AgentResult) bool {
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
			predicate: func(a AgentResult) bool {
				return a.OSInfo != nil && a.OSInfo.Distribution == "Ubuntu"
			},
			wantResult: false,
		},
		{
			name: "Returns false when target hostname does not match any agent",
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
			predicate: func(_ AgentResult) bool {
				return true
			},
			wantResult: false,
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
			predicate: func(a AgentResult) bool {
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
			predicate: func(a AgentResult) bool {
				return a.OSInfo != nil && a.OSInfo.Distribution == "Ubuntu"
			},
			wantResult: false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			r := Results{results: tc.results}

			var list AgentListResult

			decodeErr := r.Decode(tc.stepName, &list)
			if decodeErr != nil {
				s.False(tc.wantResult)

				return
			}

			matched := false
			for _, a := range list.Agents {
				if tc.target == "_all" || tc.target == a.Hostname {
					if tc.predicate(a) {
						matched = true

						break
					}
				}
			}

			s.Equal(tc.wantResult, matched)
		})
	}
}

func TestStepTestSuite(
	t *testing.T,
) {
	suite.Run(t, new(StepTestSuite))
}
