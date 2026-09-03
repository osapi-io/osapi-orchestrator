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

import "github.com/osapi-io/osapi-orchestrator/internal/engine"

// levelMode says whether the tasks in a level ran at the same time.
// Named rather than a bare bool so a call site says which it passed.
type levelMode int

const (
	modeSequential levelMode = iota
	modeParallel
)

// String names the mode as it appears in rendered output.
func (m levelMode) String() string {
	if m == modeParallel {
		return "parallel"
	}

	return "sequential"
}

// renderer defines the internal output contract for plan execution.
type renderer interface {
	PlanStart(summary engine.PlanSummary)
	PlanDone(report *Report)
	LevelStart(level int, tasks []string, mode levelMode)
	LevelDone(level int, changed int, total int, mode levelMode)
	TaskStart(name string, detail string)
	TaskDone(result engine.TaskResult)
	TaskRetry(name string, attempt int, err error)
	TaskSkip(name string, reason string)
}
