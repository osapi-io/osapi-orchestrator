package engine

import (
	"context"
	"fmt"
	"strings"

	osapiclient "github.com/osapi-io/osapi/pkg/sdk/client"
)

// Plan is a DAG of tasks with dependency edges.
type Plan struct {
	client *osapiclient.Client
	tasks  []*Task
	config PlanConfig
}

// NewPlan creates a new plan bound to an OSAPI client.
func NewPlan(
	client *osapiclient.Client,
	opts ...PlanOption,
) *Plan {
	cfg := PlanConfig{
		OnErrorStrategy: StopAll,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return &Plan{
		client: client,
		config: cfg,
	}
}

// Client returns the OSAPI client bound to this plan.
func (p *Plan) Client() *osapiclient.Client {
	return p.client
}

// Config returns the plan configuration.
func (p *Plan) Config() PlanConfig {
	return p.config
}

// TaskFunc creates a functional task, adds it to the plan, and
// returns it.
func (p *Plan) TaskFunc(
	name string,
	fn TaskFn,
) *Task {
	t := NewTaskFunc(name, fn)
	p.tasks = append(p.tasks, t)

	return t
}

// TaskFuncWithResults creates a functional task that receives
// completed results from prior tasks, adds it to the plan, and
// returns it.
func (p *Plan) TaskFuncWithResults(
	name string,
	fn TaskFnWithResults,
) *Task {
	t := NewTaskFuncWithResults(name, fn)
	p.tasks = append(p.tasks, t)

	return t
}

// Tasks returns all tasks in the plan.
func (p *Plan) Tasks() []*Task {
	return p.tasks
}

// Explain returns a human-readable representation of the execution
// plan showing levels, parallelism, dependencies, and guards.
func (p *Plan) Explain() string {
	levels, err := p.Levels()
	if err != nil {
		return fmt.Sprintf("invalid plan: %s", err)
	}

	var b strings.Builder

	writef(&b, "Plan: %d tasks, %d levels\n", len(p.tasks), len(levels))

	for i, level := range levels {
		writef(&b, "\nLevel %d%s:\n", i, parallelSuffix(len(level)))

		for _, t := range level {
			writef(&b, "  %s\n", explainTask(t))
		}
	}

	return b.String()
}

// writef writes to a strings.Builder, whose Write never returns an error. The
// ignore is here once rather than at each of the call sites.
func writef(b *strings.Builder, format string, a ...any) {
	_, _ = fmt.Fprintf(b, format, a...)
}

// parallelSuffix marks a level holding more than one task.
func parallelSuffix(taskCount int) string {
	if taskCount > 1 {
		return " (parallel)"
	}

	return ""
}

// explainTask renders one task as its name, what it waits on, and its flags.
func explainTask(t *Task) string {
	var b strings.Builder

	writef(&b, "%s [fn]", t.name)

	if len(t.deps) > 0 {
		names := make([]string, len(t.deps))
		for i, dep := range t.deps {
			names[i] = dep.name
		}

		writef(&b, " <- %s", strings.Join(names, ", "))
	}

	if flags := taskFlags(t); len(flags) > 0 {
		writef(&b, " (%s)", strings.Join(flags, ", "))
	}

	return b.String()
}

// taskFlags lists the conditions attached to a task.
func taskFlags(t *Task) []string {
	var flags []string

	if t.requiresChange {
		flags = append(flags, "only-if-changed")
	}

	if t.guard != nil {
		flags = append(flags, "when")
	}

	return flags
}

// Levels returns the levelized DAG -- tasks grouped into execution
// levels where all tasks in a level can run concurrently.
// Returns an error if the plan fails validation.
func (p *Plan) Levels() ([][]*Task, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}

	return levelize(p.tasks), nil
}

// Validate checks the plan for errors: duplicate names and cycles.
func (p *Plan) Validate() error {
	names := make(map[string]bool, len(p.tasks))

	for _, t := range p.tasks {
		if names[t.name] {
			return fmt.Errorf("duplicate task name: %q", t.name)
		}

		names[t.name] = true
	}

	return p.detectCycle()
}

// Run validates the plan, resolves the DAG, and executes tasks.
func (p *Plan) Run(
	ctx context.Context,
) (*Report, error) {
	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("plan validation: %w", err)
	}

	runner := newRunner(p)

	return runner.run(ctx)
}

// detectCycle uses DFS to find cycles in the dependency graph.
// nodeColor tracks depth-first search state. A gray node is on the path
// being walked right now, so an edge back to one closes a cycle.
type nodeColor int

const (
	nodeWhite nodeColor = iota // unvisited, and the zero value the map returns
	nodeGray                   // on the current path
	nodeBlack                  // fully explored
)

func (p *Plan) detectCycle() error {
	color := make(map[string]nodeColor, len(p.tasks))

	for _, t := range p.tasks {
		if color[t.name] != nodeWhite {
			continue
		}

		if err := visitTask(t, color); err != nil {
			return err
		}
	}

	return nil
}

// visitTask walks a task's dependencies depth-first, reporting the first
// edge that closes a cycle.
func visitTask(
	t *Task,
	color map[string]nodeColor,
) error {
	color[t.name] = nodeGray

	for _, dep := range t.deps {
		switch color[dep.name] {
		case nodeGray:
			return fmt.Errorf(
				"cycle detected: %q depends on %q",
				t.name,
				dep.name,
			)
		case nodeWhite:
			if err := visitTask(dep, color); err != nil {
				return err
			}
		default:
			// black: already fully explored, and no cycle came of it.
		}
	}

	color[t.name] = nodeBlack

	return nil
}
