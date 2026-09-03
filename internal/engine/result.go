// Package engine provides DAG-based task orchestration primitives.
package engine

import (
	"time"
)

// Status represents the outcome of a task execution.
type Status string

// Task execution statuses.
const (
	// StatusPending and StatusRunning are reserved for future
	// streaming status support. The runner does not currently
	// assign these — tasks go directly to a terminal status.
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusChanged   Status = "changed"
	StatusUnchanged Status = "unchanged"
	StatusSkipped   Status = "skipped"
	StatusFailed    Status = "failed"
)

// HostResult represents a single host's response within a broadcast
// operation.
type HostResult struct {
	Hostname    string
	Status      string
	Changed     bool
	Error       string
	Data        map[string]any
	JobDuration time.Duration
}

// Result is the outcome of a single task execution.
type Result struct {
	JobID       string
	Changed     bool
	Data        map[string]any
	Status      Status
	JobDuration time.Duration
	HostResults []HostResult
}

// TaskResult records the full execution details of a task.
type TaskResult struct {
	JobID       string
	Name        string
	Status      Status
	Changed     bool
	Duration    time.Duration
	JobDuration time.Duration
	Error       error
	Data        map[string]any
	HostResults []HostResult
}

// Results is a map of task name to Result, used for conditional logic.
type Results map[string]*Result

// Get returns the Result for the named task, or nil if not found.
func (r Results) Get(
	name string,
) *Result {
	return r[name]
}
