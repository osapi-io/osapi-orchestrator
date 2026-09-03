# Contributing

Contributions to osapi-orchestrator are very welcome, but we ask that you read
this document before submitting a PR. It covers everything you need:
prerequisites, setup, the conventions code follows, how to add an operation, and
the pull request workflow.

## Before you start

- Read the [Code of Conduct](CODE_OF_CONDUCT.md). It applies to every
  interaction in this repo.

- **Design records.** The conventions binding this repository are specified in
  [osapi-io/specs](https://github.com/osapi-io/specs) under
  `components/osapi-orchestrator/`, whose `.specify/memory/` is the standing
  record. Design reasoning for a change lives there, not here. A design document
  kept in this repository goes stale the moment the code moves past it, and
  nothing catches the drift.

- **Check existing work.** Is there an existing PR? Are there issues discussing
  the feature/change you want to make? Please make sure you consider/address
  these discussions in your work.

- **Backwards compatibility.** Will your change break existing consumers of
  osapi-orchestrator? It is much more likely that your change will be merged if
  it is backwards compatible. Is there an approach you can take that maintains
  this compatibility? If not, consider opening an issue first so that API
  changes can be discussed before you invest your time into a PR.

## Prerequisites

Install tools using [mise]:

```bash
mise install
```

- **[Go].** osapi-orchestrator is written in Go. We always support the latest
  two major Go versions, so make sure your version is recent enough.
- **[just].** Task runner used for building, testing, formatting, and other
  development workflows. Install with `brew install just`.

### Claude Code

If you use [Claude Code] for development, install this plugin from the default
marketplace:

```
/plugin install commit-commands@claude-plugins-official
```

- **commit-commands.** provides `/commit` and `/commit-push-pr` slash commands
  that follow the project's commit conventions automatically.

**Do not use superpowers.** Spec Kit governs specification, planning, and
implementation, and the design record for a change lives in
[osapi-io/specs](https://github.com/osapi-io/specs). A second workflow over that
ground gives two answers to which artifact is authoritative, and the answer that
loses is the one nobody reads. Nothing superpowers produces is committed.

## Setup

Fetch shared justfiles and install all dependencies:

```bash
just fetch
just deps
```

## Project structure

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

- **`pkg/orchestrator/`.** User-facing DSL
  - Typed operation constructors (NodeHostnameGet, CommandExec, etc.)
  - Uses SDK types directly (`osapi.HostnameResult`, `osapi.Agent`, etc.)
  - Porcelain over osapi-sdk's orchestrator engine

## Code style

Go code is formatted by [gofumpt] and linted using [golangci-lint], enforced by
CI.

```bash
just go-fmt-check   # Check formatting
just go-fmt         # Auto-fix formatting
just go-vet         # Run linter
```

The linters that run are declared in `.golangci.yml`. Read them there rather
than looking for a list here. A copied list goes stale the first time the
configuration changes. Generated files (`*.gen.go`, `*.pb.go`) are excluded from
formatting.

### Documentation

Markdown files are formatted with [mdformat] through `uvx`. This style is
enforced by CI.

```bash
just md-fmt-check   # Check formatting
just md-fmt         # Auto-fix formatting
```

## Code standards

### Function signatures

Functions with parameters use multi-line format, one parameter per line, with
the closing parenthesis and the return types on a line of their own:

```go
func FunctionName(
    param1 type1,
    param2 type2,
) (returnType, error) {
}
```

Functions taking no parameters stay on one line:

```go
func Name() string {
}
```

Adding a parameter then shows as one added line rather than a rewritten
signature.

### File naming

Name a file for what it holds. Avoid `helpers.go`, `utils.go`, and names of that
kind: they describe where code was put rather than what it is, and they
accumulate whatever has no other home.

`types.go` holds only type declarations: structs, interfaces, constants, and
aliases. A function belongs in a file named for what it does.

A test file is named for the production file it tests. Where tests grow too
large to read, split the production file first so each test file keeps a
counterpart, rather than splitting tests away from the file they cover.

### Go patterns

- Error wrapping: `fmt.Errorf("context: %w", err)`, so the chain names each
  layer it passed through and stays inspectable with `errors.Is` and
  `errors.As`.
- Early returns rather than nesting the successful path inside conditionals.
- Unused parameters: rename to `_`.
- Import order: standard library, third party, then local, separated by blank
  lines.

### Test doubles

A double for an interface this organization defines is generated with `mockgen`
and committed. Do not write a struct by hand to satisfy one.

Generated mocks live in a `mocks` package beside the code they mock, produced by
a `generate.go` holding the directive:

```go
package mocks

//go:generate go tool go.uber.org/mock/mockgen -source=../types.go -destination=types.gen.go -package=mocks
```

The generator is resolved through the module's tool dependencies, so every
checkout runs the version `go.mod` records. Destination files end in `.gen.go`
and are committed. Do not use `gen/` for mocks. That name is taken by API code
generation.

When the interface is **unexported**, a sibling package cannot work: the mock
has to import the package to name the types in the interface, and the package's
own tests have to import the mock. Generate it into the package instead, with a
destination scoped to tests so the mocking library stays out of the dependency
graph of anything that imports the package:

```go
// generate.go, in the package that declares the interface
package thispackage

//go:generate go tool go.uber.org/mock/mockgen -source=thing.go -destination=thing.gen_test.go -package=thispackage
```

Either way the directives live in a `generate.go` that holds no code, and the
generated file carries `.gen` so a reader knows not to edit it.

Where call sites would otherwise repeat the same expectations, write a
constructor returning a configured mock rather than introducing a hand-written
type. The generated mock is still what satisfies the interface.

Three doubles are written by hand, because generating them buys nothing:

- One standing in for a standard library interface: `net.Conn`, `fs.File`,
  `io.Writer`, `slog.Handler`. Those do not move when our code does.
- One carrying a real implementation of the behavior under test, such as signing
  with a genuinely generated key pair.
- A recorder for a dependency called from a goroutine the test cannot join,
  where a generated mock would assert a call count at a moment the test cannot
  establish. State that reason where the recorder is defined.

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
module. Change both together.

### Test file conventions

- Public tests: `*_public_test.go` in the package's `_test` package, exercising
  the exported surface. This is the default.
- Internal tests: `*_test.go` in the same package, for what the exported surface
  cannot reach.
- Suite naming: `*_public_test.go` → `{Name}PublicTestSuite`, `*_test.go` →
  `{Name}TestSuite`.
- `testify/suite` with table-driven cases.
- One suite method per function under test. Success, errors, and edge cases are
  rows in one table, not separate methods.
- `export_test.go` exposes unexported symbols to external tests, by alias or by
  setter. Do not use an alias to re-cover behavior the caller's own test already
  reaches; a helper with its own contract is what the pattern is for.

External tests in this repository live in `package orchestrator_test`, and
tables carry `validateFunc` callbacks.

## Adding a new operation

When adding a new typed constructor (e.g., `NodeRebootDo`), follow these steps
in order. Every operation must ship with tests, docs, and an example.

### Step 1: operation constructor

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
- The `engine` import is `internal/engine`, used only inside
  `pkg/orchestrator/`, never by external consumers

### Step 2: tests

Two test files must be updated:

**`pkg/orchestrator/ops_test.go`** (internal, httptest pattern) tests the full
HTTP round-trip with a mock server:

- Create an `httptest.Server` that returns canned JSON
- Exercise the constructor and verify the result via `report.Decode()`
- Cover success, error, and edge-case scenarios as table rows

**`pkg/orchestrator/ops_public_test.go`** (public, step-creation pattern) tests
that the constructor creates valid steps:

- Verify the step is non-nil and has the expected task name
- One suite method per operation, all scenarios as rows

Target 100% coverage on both files.

### Step 3: operation doc

Create `docs/operations/{domain}/{operation}.md` following the existing template
in that domain directory. Every doc must include these sections:

- **Description** (h1 heading with the method name)
- **Usage.** minimal Go snippet showing the constructor call
- **Parameters.** table of all parameters with types and descriptions
- **Result Type.** `Decode()` snippet and field table
- **Idempotency.** one of: Read-only, Idempotent (Yes), Non-idempotent (No)
- **Permissions.** required OSAPI permission (e.g., `node:write`)
- **Example.** link to the example file where this operation is used:
  ```
  See
  [`examples/operations/{domain}.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/{domain}.go)
  for a complete working example.
  ```

### Step 4: update the domain landing page and operation index

Add the operation to the table in the domain landing page
`docs/operations/{domain}/README.md`. Update the operation count in
`docs/operations/README.md` if the total changes.

### Step 5: example

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
  you're demonstrating too much. Split it.
- **Operation docs link to examples**: every operation doc in
  `docs/operations/{domain}/` must link to the example file where that operation
  is demonstrated.

### Step 6: update README.md

Update the operation count and tables in the root `README.md` if the total
number of operations changes.

### Step 7: verify

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

- **Describe your changes.** Say what changed and why. A reviewer should not
  have to read the diff to learn the reason for it.
- **Issue/PR links.** Link any previous work such as related issues or PRs.
  Please describe how your changes differ to/extend this work.
- **Examples.** Add any examples or screenshots that you think are useful to
  demonstrate the effect of your changes.
- **Draft PRs.** If your changes are incomplete, but you would like to discuss
  them, open the PR as a draft and add a comment to start a discussion. Using
  comments rather than the PR description allows the description to be updated
  later while preserving any discussions.

## AI usage

This repo is written with AI assistance. All contributions are subject to the
[AI Usage Policy](AI_POLICY.md). Disclose the tool you used, and make sure you
can explain what your change does without the aid of AI tools.

## FAQ

> I want to contribute, where do I start?

All kinds of contributions are welcome, whether it's a typo fix or a shiny new
feature. You can also contribute by upvoting/commenting on issues or helping to
answer questions.

> I'm stuck, where can I get help?

If you have questions, open a [Discussion] on GitHub.

[claude code]: https://claude.ai/code
[conventional commits]: https://www.conventionalcommits.org
[discussion]: https://github.com/osapi-io/osapi-orchestrator/discussions
[go]: https://go.dev
[gofumpt]: https://github.com/mvdan/gofumpt
[golangci-lint]: https://golangci-lint.run
[just]: https://just.systems
[mdformat]: https://pypi.org/project/mdformat/
[mise]: https://mise.jdx.dev
