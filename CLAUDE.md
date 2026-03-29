# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with
code in this repository.

## Project Overview

User-facing orchestration DSL for [OSAPI][] — declarative infrastructure
operations with typed constructors, chaining, and result decoding. Built on
top of the [osapi-sdk][] orchestrator engine.

## Development Reference

For setup, prerequisites, and contributing guidelines:

- @docs/development.md - Prerequisites, setup, code style, testing, commits
- @docs/contributing.md - PR workflow and contribution guidelines
- @docs/operations/README.md - Per-operation reference (37 typed constructors)
- @docs/features/README.md - Feature guides (guards, retry, discovery, etc.)

## Code Standards (MANDATORY)

### Function Signatures

ALL function signatures MUST use multi-line format:
```go
func FunctionName(
    param1 type1,
    param2 type2,
) (returnType, error) {
}
```

### Testing

- Public tests: `*_public_test.go` in test package
  (`package orchestrator_test`) for exported functions
- Internal tests: `*_test.go` in same package (`package orchestrator`)
  for private functions
- Suite naming: `*_public_test.go` → `{Name}PublicTestSuite`,
  `*_test.go` → `{Name}TestSuite`
- Use `testify/suite` with table-driven patterns
- Table-driven structure with `validateFunc` callbacks
- One suite method per function under test — all scenarios (success,
  errors, edge cases) as rows in one table
- Avoid generic file names like `helpers.go` or `utils.go` — name
  files after what they contain

### Go Patterns

- Error wrapping: `fmt.Errorf("context: %w", err)`
- Early returns over nested if-else
- Unused parameters: rename to `_`
- Import order: stdlib, third-party, local (blank-line separated)

### Linting

golangci-lint with: errcheck, errname, goimports, govet, prealloc,
predeclared, revive, staticcheck. Generated files (`*.gen.go`, `*.pb.go`)
are excluded from formatting.

### Branching

See @docs/development.md#branching for full conventions.

When committing changes via `/commit`, create a feature branch first if
currently on `main`. Branch names use the pattern `type/short-description`
(e.g., `feat/add-dns-retry`, `fix/memory-leak`, `docs/update-readme`).

### Commit Messages

See @docs/development.md#commit-messages for full conventions.

Follow [Conventional Commits](https://www.conventionalcommits.org/) with the
50/72 rule. Format: `type(scope): description`.

When committing via Claude Code, end with:
- `🤖 Generated with [Claude Code](https://claude.ai/code)`
- `Co-Authored-By: Claude <noreply@anthropic.com>`

## Task Tracking

Implementation planning and execution uses the superpowers plugin workflows
(`writing-plans` and `executing-plans`). Plans live in `docs/plans/`.

## Adding a New Operation

When adding a new typed constructor (e.g., `NodeRebootDo`), follow these
steps in order. Every operation must ship with tests, docs, and an example.

### Step 1: Operation Constructor

Add the method to `pkg/orchestrator/ops.go`, following the existing pattern:

```go
// NodeRebootDo creates a step that reboots the target node.
func (o *Orchestrator) NodeRebootDo(
    target string,
) *Step {
    name := o.nextOpName("reboot-node")

    task := o.plan.TaskFunc(
        name,
        func(
            ctx context.Context,
            c *osapi.Client,
        ) (*sdk.Result, error) {
            resp, err := c.Node.Reboot(ctx, target)
            if err != nil {
                return nil, fmt.Errorf("reboot node: %w", err)
            }

            return sdk.CollectionResult(resp.Data, resp.RawJSON(),
                func(r osapi.RebootResult) sdk.HostResult {
                    return sdk.HostResult{
                        Hostname: r.Hostname,
                        Status:   r.Status,
                        Changed:  r.Changed,
                        Error:    r.Error,
                    }
                },
            )
        },
    )

    return &Step{task: task}
}
```

Key rules:
- Use `o.nextOpName("verb-noun")` for the step name
- Always include `Status: r.Status` in the `HostResult` mapper
- Wrap errors with context: `fmt.Errorf("verb noun: %w", err)`
- Use `sdk.CollectionResult` for node-targeted operations (returns
  per-host results), `sdk.StructToMap` for non-collection responses

### Step 2: Tests

Two test files must be updated:

**`pkg/orchestrator/ops_test.go`** (internal, httptest pattern) — tests
the full HTTP round-trip with a mock server:
- Create an `httptest.Server` that returns canned JSON
- Exercise the constructor and verify the result via `report.Decode()`
- Cover success, error, and edge-case scenarios as table rows

**`pkg/orchestrator/ops_public_test.go`** (public, step-creation
pattern) — tests that the constructor creates valid steps:
- Verify the step is non-nil and has the expected task name
- One suite method per operation, all scenarios as rows

Target 100% coverage on both files.

### Step 3: Operation Doc

Create `docs/operations/{domain}/{operation}.md` following the existing
template in that domain directory. Every doc must include these sections:

- **Description** (h1 heading with the method name)
- **Usage** — minimal Go snippet showing the constructor call
- **Parameters** — table of all parameters with types and descriptions
- **Result Type** — `Decode()` snippet and field table
- **Idempotency** — one of: Read-only, Idempotent (Yes), Non-idempotent (No)
- **Permissions** — required OSAPI permission (e.g., `node:write`)
- **Example** — link to the example file where this operation is used:
  ```
  See
  [`examples/operations/{domain}.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/{domain}.go)
  for a complete working example.
  ```

### Step 4: Update Domain Landing Page and Operation Index

Add the operation to the table in the domain landing page
`docs/operations/{domain}/README.md`. Update the operation count in
`docs/operations/README.md` if the total changes.

### Step 5: Example

Add the operation to an existing workflow example in
`examples/operations/` that covers the same domain. Domain groupings:

| Domain  | Example file         |
| ------- | -------------------- |
| Node    | `node-info.go`       |
| Node    | `hostname-update.go` |
| Network | `dns-update.go`      |
| Command | `command.go`         |
| File    | `file-deploy.go`     |
| File    | `file-changed.go`    |
| Agent   | `agent-drain.go`     |
| Docker  | `docker.go`          |
| Cron    | `cron.go`            |
| Health  | (used as gate)       |

If no domain match exists, create a new `{domain}.go` file. Every
operation must appear in at least one runnable example.

### Step 6: Update README.md

Update the operation count and tables in the root `README.md` if the
total number of operations changes.

### Step 7: Verify

```bash
go build ./...                                       # compiles
go test ./... -count=1                               # tests pass
cd examples/operations && go build *.go              # examples compile
cd examples/features && go build *.go                # feature examples compile
```

[OSAPI]: https://github.com/osapi-io/osapi
[osapi-sdk]: https://github.com/osapi-io/osapi/tree/main/pkg/sdk
