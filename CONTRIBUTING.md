# Contributing

Contributions to osapi-orchestrator are very welcome, but we ask that you read
this document before submitting a PR. It covers everything you need:
prerequisites, setup, the conventions code follows, how to add an operation, and
the pull request workflow.

## Before you start

- Read the [Code of Conduct](CODE_OF_CONDUCT.md). It applies to every
  interaction in this repo.

- **Design records** — Reasoning behind a change lives in
  [osapi-io/specs](https://github.com/osapi-io/specs) as a change, not in this
  repository. A design document kept here goes stale the moment the code moves
  past it, with nothing to catch the drift.

- **Check existing work** — Is there an existing PR? Are there issues discussing
  the feature/change you want to make? Please make sure you consider/address
  these discussions in your work.

- **Backwards compatibility** — Will your change break existing consumers of
  osapi-orchestrator? It is much more likely that your change will be merged if
  it is backwards compatible. Is there an approach you can take that maintains
  this compatibility? If not, consider opening an issue first so that API
  changes can be discussed before you invest your time into a PR.

## Prerequisites

Install tools using [mise]:

```bash
mise install
```

- **[Go]** — osapi-orchestrator is written in Go. We always support the latest
  two major Go versions, so make sure your version is recent enough.
- **[just]** — Task runner used for building, testing, formatting, and other
  development workflows. Install with `brew install just`.

### Claude Code

If you use [Claude Code] for development, install these plugins from the default
marketplace:

```
/plugin install commit-commands@claude-plugins-official
/plugin install superpowers@claude-plugins-official
```

- **commit-commands** — provides `/commit` and `/commit-push-pr` slash commands
  that follow the project's commit conventions automatically.
- **superpowers** — provides structured workflows for planning, TDD, debugging,
  code review, and git worktree isolation.

## Setup

Fetch shared justfiles and install all dependencies:

```bash
just fetch
just deps
```

## Project Structure

```
pkg/orchestrator/          # User-facing DSL
  ops.go                   # Typed operation constructors
  step.go                  # Step chaining (guards, retry, ordering)
  result.go                # Result types (HostResult, Results, Report)
  host_status.go           # HostStatusOk/Skipped/Failed constants
  renderer_lipgloss.go     # Terminal output renderer
  types.go                 # Orchestrator struct
  options.go               # Option types (verbose, upload, retry)
docs/
  operations/              # Operation reference (one doc per operation)
    README.md              # Master index linking to domains
    node/                  # Domain subdirectory with landing page
      README.md            # "Node Management" — ops table, permissions, example
      hostname-get.md      # Individual operation reference
      ...
    network/, command/, docker/, cron/, file/, agent/, health/
  features/                # Cross-cutting feature guides
    guards.md, broadcast.md, retry.md, ...
examples/
  operations/              # Runnable workflow examples (one per domain)
    node-info.go, docker.go, cron.go, ...
  features/                # Runnable feature examples
    host-status.go, guards.go, broadcast.go, ...
```

## Package Structure

- **`pkg/orchestrator/`** — User-facing DSL
  - Typed operation constructors (NodeHostnameGet, CommandExec, etc.)
  - Uses SDK types directly (`osapi.HostnameResult`, `osapi.Agent`, etc.)
  - Porcelain over osapi-sdk's orchestrator engine

## Code style

Go code should be formatted by [`gofumpt`][gofumpt] and linted using
[`golangci-lint`][golangci-lint]. This style is enforced by CI.

```bash
just go-fmt-check   # Check formatting
just go-fmt         # Auto-fix formatting
just go-vet         # Run linter
```

### Documentation

Markdown files are formatted with [mdformat] through `uvx`. This style is
enforced by CI.

```bash
just md-fmt-check   # Check formatting
just md-fmt         # Auto-fix formatting
```

## Code standards

The Go conventions this repository follows — multi-line function signatures,
naming a file for what it holds, table-driven `testify/suite` tests, suite
naming, generated mocks, and the error-wrapping and import-order baseline — are
specified in the `go-code-standards` capability in
[osapi-io/specs](https://github.com/osapi-io/specs). That capability is the
source, and this guide does not restate it.

golangci-lint runs errcheck, errname, goimports, govet, prealloc, predeclared,
revive, and staticcheck. Generated files (`*.gen.go`, `*.pb.go`) are excluded
from formatting.

Tests exercise a real HTTP server via `httptest.Server` rather than mocking the
SDK client.

## Testing

```bash
just test           # Run all tests (lint + unit + coverage)
just go-unit       # Run unit tests only
just go-unit-cov   # Generate coverage report
go test -run TestName -v ./pkg/orchestrator/...  # Run a single test
```

Coverage is gated at 100%. `just test` fails if total coverage drops below it,
so a change that adds untested code fails locally and in CI:

```bash
just go-unit-cov-check   # Report coverage and fail below the target
```

The target is declared in `.github/codecov.yml` and in the shared `go` justfile
module — change both together.

### Test file conventions

Test package layout, suite naming, and table-driven cases are specified in
`go-code-standards`. In this repository external tests live in
`package orchestrator_test`, and tables carry `validateFunc` callbacks.

## Adding a new operation

When adding a new typed constructor (e.g., `NodeRebootDo`), follow these steps
in order. Every operation must ship with tests, docs, and an example.

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
        ) (*engine.Result, error) {
            resp, err := c.Node.Reboot(ctx, target)
            if err != nil {
                return nil, fmt.Errorf("reboot node: %w", err)
            }

            return engine.CollectionResult(resp.Data, resp.RawJSON(),
                func(r osapi.RebootResult) engine.HostResult {
                    return engine.HostResult{
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
- Always include `Status: r.Status` in the `engine.HostResult` mapper
- Wrap errors with context: `fmt.Errorf("verb noun: %w", err)`
- Use `engine.CollectionResult` for node-targeted operations (returns per-host
  results), `engine.StructToMap` for non-collection responses
- The `engine` import is `internal/engine` — only used inside
  `pkg/orchestrator/`, never by external consumers

### Step 2: Tests

Two test files must be updated:

**`pkg/orchestrator/ops_test.go`** (internal, httptest pattern) — tests the full
HTTP round-trip with a mock server:

- Create an `httptest.Server` that returns canned JSON
- Exercise the constructor and verify the result via `report.Decode()`
- Cover success, error, and edge-case scenarios as table rows

**`pkg/orchestrator/ops_public_test.go`** (public, step-creation pattern) —
tests that the constructor creates valid steps:

- Verify the step is non-nil and has the expected task name
- One suite method per operation, all scenarios as rows

Target 100% coverage on both files.

### Step 3: Operation Doc

Create `docs/operations/{domain}/{operation}.md` following the existing template
in that domain directory. Every doc must include these sections:

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

Add the operation to an existing workflow example in `examples/operations/` that
covers the same domain. Domain groupings:

| Domain      | Example file         |
| ----------- | -------------------- |
| Node        | `node-info.go`       |
| Node        | `hostname-update.go` |
| Network     | `dns-update.go`      |
| Network     | `ping.go`            |
| Interface   | `interface.go`       |
| Route       | `route.go`           |
| Command     | `command.go`         |
| File        | `file-deploy.go`     |
| File        | `file-changed.go`    |
| Agent       | `agent-drain.go`     |
| Docker      | `docker.go`          |
| Cron        | `cron.go`            |
| Sysctl      | `sysctl.go`          |
| NTP         | `ntp.go`             |
| Timezone    | `timezone.go`        |
| Service     | `service.go`         |
| Package     | `package.go`         |
| User        | `user.go`            |
| Group       | `group.go`           |
| Certificate | `certificate.go`     |
| Process     | `process.go`         |
| Power       | `power.go`           |
| Log         | `log.go`             |
| Health      | (used as gate)       |

If no domain match exists, create a new `{domain}.go` file. Every operation must
appear in at least one runnable example.

#### Example conventions

- **Self-contained**: cleanup at the start (separate plan with
  `ContinueOnError()`), execute, verify. Must be repeatable.
- **One purpose per file**: demonstrate one domain's operations. Don't mix in
  other features (parallel, verbose, broadcast).
- **Cleanup plan pattern**: use a separate `orchestrator.New()` for cleanup with
  `ContinueOnError()`, then a main plan for the workflow.
- **Platform safety**: operations that may not work everywhere use
  `ContinueOnError()` so the example doesn't crash.
- **Decode and print**: decode at least one result so the example isn't silent.
  Use `report.Decode("step-name", &typedStruct)`.
- **Keep it short**: under ~100 lines of code (excluding license). If longer,
  you're demonstrating too much — split it.
- **Operation docs link to examples**: every operation doc in
  `docs/operations/{domain}/` must link to the example file where that operation
  is demonstrated.

### Step 6: Update README.md

Update the operation count and tables in the root `README.md` if the total
number of operations changes.

### Step 7: Verify

```bash
go build ./...                                       # compiles
go test ./... -count=1                               # tests pass
cd examples/operations && go build *.go              # examples compile
cd examples/features && go build *.go                # feature examples compile
```

## Before committing

Run `just ready` before committing to ensure generated code, package docs,
formatting, and lint are all up to date:

```bash
just ready   # generate, go-docs, go-fmt, go-vet
```

## Branching

All changes should be developed on feature branches. Create a branch from `main`
using the naming convention `type/short-description`, where `type` matches the
[Conventional Commits] type:

- `feat/add-retry-logic`
- `fix/null-pointer-crash`
- `docs/update-api-reference`
- `refactor/simplify-handler`
- `chore/update-dependencies`

When using Claude Code's `/commit` command, a branch will be created
automatically if you are on `main`.

## Commit messages

Follow [Conventional Commits] with the 50/72 rule:

- **Subject line**: max 50 characters, imperative mood, capitalized, no period
- **Body**: wrap at 72 characters, separated from subject by a blank line
- **Format**: `type(scope): description`
- **Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`
- Summarize the "what" and "why", not the "how"

Try to write meaningful commit messages and avoid having too many commits on a
PR. Most PRs should likely have a single commit (although for bigger PRs it may
be reasonable to split it in a few). Git squash and rebase is your friend!

## Submitting a PR

- **Describe your changes** — Ensure that you provide a comprehensive
  description of your changes.
- **Issue/PR links** — Link any previous work such as related issues or PRs.
  Please describe how your changes differ to/extend this work.
- **Examples** — Add any examples or screenshots that you think are useful to
  demonstrate the effect of your changes.
- **Draft PRs** — If your changes are incomplete, but you would like to discuss
  them, open the PR as a draft and add a comment to start a discussion. Using
  comments rather than the PR description allows the description to be updated
  later while preserving any discussions.

## AI usage

All contributions are subject to the [AI Usage Policy](AI_POLICY.md) — disclose
the tool you used, and make sure you can explain what your change does without
the aid of AI tools.

## FAQ

> I want to contribute, where do I start?

All kinds of contributions are welcome, whether it's a typo fix or a shiny new
feature. You can also contribute by upvoting/commenting on issues or helping to
answer questions.

> I'm stuck, where can I get help?

If you have questions, feel free to open a [Discussion] on GitHub.

[claude code]: https://claude.ai/code
[conventional commits]: https://www.conventionalcommits.org
[discussion]: https://github.com/osapi-io/osapi-orchestrator/discussions
[go]: https://go.dev
[gofumpt]: https://github.com/mvdan/gofumpt
[golangci-lint]: https://golangci-lint.run
[just]: https://just.systems
[mdformat]: https://pypi.org/project/mdformat/
[mise]: https://mise.jdx.dev
