// Package engine provides DAG-based task orchestration primitives.
package engine

import (
	"fmt"
	"strings"
	"time"
)

// StepSummary describes a single execution step (DAG level).
type StepSummary struct {
	Tasks    []string
	Parallel bool
}

// PlanSummary describes the execution plan before it runs.
type PlanSummary struct {
	TotalTasks int
	Steps      []StepSummary
}

// Report is the aggregate output of a plan execution.
type Report struct {
	Tasks    []TaskResult
	Duration time.Duration
}

// Summary returns a human-readable summary of the report.
// summaryOrder is the order counts appear in a summary line.
var summaryOrder = []Status{
	StatusChanged,
	StatusUnchanged,
	StatusSkipped,
	StatusFailed,
}

func (r *Report) Summary() string {
	counts := make(map[Status]int, len(summaryOrder))
	for _, t := range r.Tasks {
		counts[t.Status]++
	}

	parts := []string{
		fmt.Sprintf("%d tasks", len(r.Tasks)),
	}

	// Each status names itself — StatusChanged is "changed" — and a
	// status nothing landed in is left out rather than reported as zero.
	for _, status := range summaryOrder {
		if n := counts[status]; n > 0 {
			parts = append(parts, fmt.Sprintf("%d %s", n, status))
		}
	}

	return strings.Join(parts, ", ")
}
