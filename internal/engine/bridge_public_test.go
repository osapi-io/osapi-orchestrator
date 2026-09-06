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
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/osapi-io/osapi-orchestrator/internal/engine"
	client "github.com/osapi-io/osapi/pkg/sdk/client"
)

type BridgePublicTestSuite struct {
	suite.Suite
}

func TestBridgePublicTestSuite(t *testing.T) {
	suite.Run(t, new(BridgePublicTestSuite))
}

type testStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type testNested struct {
	Label string     `json:"label"`
	Inner testStruct `json:"inner"`
}

func (s *BridgePublicTestSuite) TestStructToMap() {
	tests := []struct {
		name         string
		setupFn      func()
		teardownFn   func()
		input        any
		validateFunc func(m map[string]any)
	}{
		{
			name:  "converts struct with json tags to map",
			input: testStruct{Name: "web-01", Value: 42},
			validateFunc: func(m map[string]any) {
				s.Require().NotNil(m)
				s.Equal("web-01", m["name"])
				s.Equal(float64(42), m["value"])
			},
		},
		{
			name:  "returns nil for nil input",
			input: nil,
			validateFunc: func(m map[string]any) {
				s.Nil(m)
			},
		},
		{
			name: "handles nested structs",
			input: testNested{
				Label: "parent",
				Inner: testStruct{Name: "child", Value: 7},
			},
			validateFunc: func(m map[string]any) {
				s.Require().NotNil(m)
				s.Equal("parent", m["label"])

				inner, ok := m["inner"].(map[string]any)
				s.Require().True(ok)
				s.Equal("child", inner["name"])
				s.Equal(float64(7), inner["value"])
			},
		},
		{
			name: "converts SDK type with json tags",
			input: client.HostnameResult{
				Hostname: "web-01",
				Changed:  true,
			},
			validateFunc: func(m map[string]any) {
				s.Require().NotNil(m)
				s.Equal("web-01", m["hostname"])
				s.Equal(true, m["changed"])
			},
		},
		{
			name:  "returns nil for unmarshalable input",
			input: make(chan int),
			validateFunc: func(m map[string]any) {
				s.Nil(m)
			},
		},
		{
			name: "returns nil when unmarshal fails",
			setupFn: func() {
				engine.SetJSONUnmarshalFn(func(
					_ []byte,
					_ any,
				) error {
					return fmt.Errorf("forced unmarshal error")
				})
			},
			teardownFn: engine.ResetJSONUnmarshalFn,
			input:      testStruct{Name: "test"},
			validateFunc: func(m map[string]any) {
				s.Nil(m)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.setupFn != nil {
				tt.setupFn()
			}
			if tt.teardownFn != nil {
				defer tt.teardownFn()
			}

			got := engine.StructToMap(tt.input)
			tt.validateFunc(got)
		})
	}
}

func (s *BridgePublicTestSuite) TestCollectionResult() {
	mapper := func(r client.HostnameResult) engine.HostResult {
		return engine.HostResult{
			Hostname: r.Hostname,
			Changed:  r.Changed,
			Error:    r.Error,
		}
	}

	tests := []struct {
		name         string
		col          client.Collection[client.HostnameResult]
		rawJSON      []byte
		toHost       func(client.HostnameResult) engine.HostResult
		expectErr    bool
		validateFunc func(result *engine.Result)
	}{
		{
			name: "single result with auto-populated data",
			col: client.Collection[client.HostnameResult]{
				Results: []client.HostnameResult{
					{Hostname: "web-01", Changed: false},
				},
				JobID: "job-123",
			},
			toHost: mapper,
			validateFunc: func(result *engine.Result) {
				s.Equal("job-123", result.JobID)
				s.False(result.Changed)
				s.Require().Len(result.HostResults, 1)

				hr := result.HostResults[0]
				s.Equal("web-01", hr.Hostname)
				s.False(hr.Changed)
				s.Require().NotNil(hr.Data, "Data should be auto-populated via StructToMap")
				s.Equal("web-01", hr.Data["hostname"])
			},
		},
		{
			name: "multiple results with changed true when any host changed",
			col: client.Collection[client.HostnameResult]{
				Results: []client.HostnameResult{
					{Hostname: "web-01", Changed: false},
					{Hostname: "web-02", Changed: true},
				},
				JobID: "job-456",
			},
			toHost: mapper,
			validateFunc: func(result *engine.Result) {
				s.Equal("job-456", result.JobID)
				s.True(result.Changed)
				s.Len(result.HostResults, 2)
				s.False(result.HostResults[0].Changed)
				s.True(result.HostResults[1].Changed)
			},
		},
		{
			name: "empty results returns result with empty host results",
			col: client.Collection[client.HostnameResult]{
				Results: []client.HostnameResult{},
				JobID:   "job-789",
			},
			toHost: mapper,
			validateFunc: func(result *engine.Result) {
				s.Equal("job-789", result.JobID)
				s.False(result.Changed)
				s.Empty(result.HostResults)
			},
		},
		{
			name: "data auto-populated via StructToMap when mapper leaves it nil",
			col: client.Collection[client.HostnameResult]{
				Results: []client.HostnameResult{
					{Hostname: "db-01", Changed: false, Error: "timeout"},
				},
				JobID: "job-auto",
			},
			toHost: func(r client.HostnameResult) engine.HostResult {
				return engine.HostResult{
					Hostname: r.Hostname,
					Changed:  r.Changed,
					Error:    r.Error,
				}
			},
			validateFunc: func(result *engine.Result) {
				hr := result.HostResults[0]
				s.Require().NotNil(hr.Data)
				s.Equal("db-01", hr.Data["hostname"])
				s.Equal("timeout", hr.Data["error"])
			},
		},
		{
			name: "data preserved when mapper sets it explicitly",
			col: client.Collection[client.HostnameResult]{
				Results: []client.HostnameResult{
					{Hostname: "app-01", Changed: true},
				},
				JobID: "job-explicit",
			},
			toHost: func(r client.HostnameResult) engine.HostResult {
				return engine.HostResult{
					Hostname: r.Hostname,
					Changed:  r.Changed,
					Data:     map[string]any{"custom": "value"},
				}
			},
			validateFunc: func(result *engine.Result) {
				hr := result.HostResults[0]
				s.Require().NotNil(hr.Data)
				s.Equal("value", hr.Data["custom"])
				_, hasHostname := hr.Data["hostname"]
				s.False(hasHostname, "mapper-set Data should not be overwritten")
			},
		},
		{
			name: "rawJSON populates Result.Data",
			col: client.Collection[client.HostnameResult]{
				Results: []client.HostnameResult{
					{Hostname: "web-01"},
				},
				JobID: "job-raw",
			},
			rawJSON: []byte(`{"job_id":"job-raw","results":[{"hostname":"web-01"}]}`),
			toHost:  mapper,
			validateFunc: func(result *engine.Result) {
				s.Require().NotNil(result.Data)
				s.Equal("job-raw", result.Data["job_id"])
			},
		},
		{
			name: "nil rawJSON leaves Result.Data nil",
			col: client.Collection[client.HostnameResult]{
				Results: []client.HostnameResult{
					{Hostname: "web-01"},
				},
				JobID: "job-nil",
			},
			rawJSON: nil,
			toHost:  mapper,
			validateFunc: func(result *engine.Result) {
				s.Nil(result.Data)
			},
		},
		{
			name: "invalid rawJSON returns error",
			col: client.Collection[client.HostnameResult]{
				Results: []client.HostnameResult{
					{Hostname: "web-01"},
				},
				JobID: "job-bad",
			},
			rawJSON:   []byte(`not valid json`),
			toHost:    mapper,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result, err := engine.CollectionResult(
				tt.col,
				tt.rawJSON,
				tt.toHost,
			)

			if tt.expectErr {
				s.Error(err)
				s.Nil(result)

				return
			}

			s.NoError(err)
			s.Require().NotNil(result)
			tt.validateFunc(result)
		})
	}
}
