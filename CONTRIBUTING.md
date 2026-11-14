# Contributing

We welcome contributions to osapi-orchestrator. Read this document before you
send a PR. It tells the prerequisites, the setup, the conventions that the code
follows, how to add an operation, and the pull request workflow.

## Before you start

- Read the [Code of Conduct](CODE_OF_CONDUCT.md). It applies to all the work and
  all the discussion in this repo.

- **Design docs.** The design of this repository is in
  [osapi-io/specs](https://github.com/osapi-io/specs) under
  `components/osapi-orchestrator/`. Read it before you change the behaviour.
  When you change the behaviour, change the design doc too. A design doc in this
  repository becomes old when the code changes, and no check finds the
  difference.

- **Check existing work.** Look for a PR that exists for the change. Look for
  issues that discuss the feature or the change that you want to make. Make sure
  that your work agrees with these discussions, or answers them.

- **Backwards compatibility.** Make sure that your change does not break the
  code that uses osapi-orchestrator now. A backwards compatible change is much
  more likely to merge. Look for a method that keeps this compatibility. If you
  cannot find one, open an issue first. Then we can discuss the API change
  before you spend time on a PR.

## Prerequisites

Install the tools with [mise]:

```bash
mise install
```

- **[Go].** osapi-orchestrator is written in Go. We support the two latest major
  Go versions. Make sure that your version is one of them.
- **[just].** The task runner. It builds, tests, and formats the code, and runs
  other development workflows. Install it with `brew install just`.

### Claude Code

If you use [Claude Code] for development, install this plugin from the default
marketplace:

```
/plugin install commit-commands@claude-plugins-official
```

- **commit-commands.** It gives the `/commit` and `/commit-push-pr` slash
  commands. These commands follow the commit conventions of the project.

**The design is in the design docs.** Write or change the page in
[osapi-io/specs](https://github.com/osapi-io/specs) before you build. If the
build shows that the page is wrong, correct the page. A design document in this
repository becomes old when the code changes, and no check finds the difference.

## Setup

Fetch the shared justfiles and install all the dependencies:

```bash
just fetch
just deps
```

## Project structure

```
pkg/orchestrator/          # User-facing DSL
  orchestrator.go          # New, Run, TaskFunc
  ops.go                   # Typed operation constructors
  step.go                  # Step chaining (guards, retry, ordering)
  result.go                # Result types (Results, TaskStatus, Report)
  host_status.go           # HostStatusOk/Skipped/Failed constants
  discover.go              # Discover and GroupByFact
  predicate.go             # Agent predicates (OS, Arch, HasLabel, ...)
  renderer.go              # The renderer interface
  renderer_lipgloss.go     # Terminal output renderer
  generate.go              # mockgen directive, no code
  options.go               # Option types (verbose, upload, retry)
  types.go                 # Orchestrator struct, aliases onto the engine
internal/engine/           # Bridge onto osapi-sdk's orchestrator engine
docs/
  operations/              # Operation reference (one doc per operation)
    README.md              # Master index linking to each group
    health/                # Group whose domain is the group itself
      README.md            # Landing page — ops table, permissions, example
      check.md             # Individual operation reference
    node/                  # Group holding both its own ops and domains
      README.md
      status-get.md
      power/               # Domain subdirectory with its own landing page
        README.md
        reboot.md
    agent/, cmd/, config/, containers/, files/, hardware/,
    networking/, security/, services/, software/
  features/                # Cross-cutting feature guides
    guards.md, broadcast.md, retry.md, ...
examples/
  operations/              # Runnable workflow examples (one per domain)
    node-info.go, docker.go, cron.go, ...
  features/                # Runnable feature examples
    host-status.go, guards.go, broadcast.go, ...
```

If the group has no inner domains, an operation doc is at
`docs/operations/{group}/{operation}.md`. If the group has inner domains, the
doc is at `docs/operations/{group}/{domain}/{operation}.md`. The landing page of
the group links to the docs that it holds.

- **`pkg/orchestrator/`.** The DSL that users call.
  - Typed operation constructors (NodeHostnameGet, CommandExec, and more)
  - It uses the SDK types directly (`osapi.HostnameResult`, `osapi.Agent`, and
    more)
  - It is a high-level layer on the orchestrator engine of osapi-sdk

## Code style

[gofumpt] formats the Go code and [golangci-lint] lints it. CI makes these
checks mandatory.

```bash
just go-fmt-check   # Check formatting
just go-fmt         # Auto-fix formatting
just go-vet         # Run linter
```

`.golangci.yml` declares the linters that run. Read them there. This page does
not copy the list, because a copy becomes wrong when the configuration changes.
The formatter does not format generated files (`*.gen.go`, `*.pb.go`).

### Documentation

[mdformat] formats the Markdown files, through `uvx`. CI makes this style
mandatory.

```bash
just md-fmt-check   # Check formatting
just md-fmt         # Auto-fix formatting
```

## Code standards

### Function signatures

A function with parameters uses more than one line. Put one parameter on each
line. Put the closing parenthesis and the return types on a line of their own:

```go
func FunctionName(
    param1 type1,
    param2 type2,
) (returnType, error) {
}
```

A function without parameters stays on one line:

```go
func Name() string {
}
```

Then a new parameter shows as one added line in the diff, not as a new
signature.

### File naming

Give a file the name of what it holds. Do not use `helpers.go`, `utils.go`, or
names of that kind. These names tell where the code is, not what the code is.
Such files also collect all the code that has no other location.

`types.go` holds only type declarations: structs, interfaces, constants, and
aliases. Put a function in a file with the name of what the function does.

A test file has the name of the production file that it tests. If the tests
become too large to read, split the production file first. Then each test file
keeps one production file. Do not move tests away from the file that they test.

### Go patterns

- Wrap errors with `fmt.Errorf("context: %w", err)`. Then the chain names each
  layer that the error went through. You can also examine the chain with
  `errors.Is` and `errors.As`.
- Return early. Do not put the success path inside conditionals.
- Rename an unused parameter to `_`.
- Put the imports in this order: standard library, third party, then local. Put
  a blank line between the groups.

### Test doubles

If this organization defines the interface, generate its double with `mockgen`
and commit it. Do not write a struct by hand for that interface.

Generated mocks are in a `mocks` package next to the code that they mock. A
`generate.go` file holds the directive that makes them:

```go
package mocks

//go:generate go tool go.uber.org/mock/mockgen -source=../types.go -destination=types.gen.go -package=mocks
```

Go finds the generator through the tool dependencies of the module. Thus every
checkout runs the version that `go.mod` records. The destination files end in
`.gen.go`, and you commit them. Do not use `gen/` for mocks. The API code
generator uses that name.

If the interface is **unexported**, a sibling package cannot work. The mock must
import the package to name the types in the interface. The tests of the package
must import the mock. Thus generate the mock into the package. Use a destination
that only tests compile. Then the mock library is not in the dependency graph of
code that imports the package:

```go
// generate.go, in the package that declares the interface
package thispackage

//go:generate go tool go.uber.org/mock/mockgen -source=thing.go -destination=thing.gen_test.go -package=thispackage
```

In the two cases, the directives are in a `generate.go` that holds no code. The
generated file name has `.gen`, so a reader knows not to edit it.

If many call sites repeat the same expectations, write a constructor that
returns a configured mock. Do not add a hand-written type. The generated mock
still satisfies the interface.

Write these three doubles by hand, because a generated double gives no benefit:

- A double for a standard library interface: `net.Conn`, `fs.File`, `io.Writer`,
  `slog.Handler`. These interfaces do not change when our code changes.
- A double that holds a real implementation of the behavior under test. An
  example is a signature with a key pair that a real generator made.
- A recorder for a dependency that a goroutine calls, if the test cannot join
  that goroutine. A generated mock would assert a call count at a time that the
  test cannot know. Write that reason where you define the recorder.

Tests use a real HTTP server through `httptest.Server`. They do not mock the SDK
client.

### File headers

Every `.go` file MUST start with the MIT license header. Refer to a Go file in
the repo for the exact format. A file with a build tag puts `//go:build` on line
1, then a blank line, then the header.

## Testing

```bash
just test           # Run all tests (lint + unit + coverage)
just go-unit       # Run unit tests only
just go-unit-cov   # Generate coverage report
go test -run TestName -v ./pkg/orchestrator/...  # Run a single test
```

The coverage gate is 100%. `just test` fails if the total coverage is less than
100%. Thus a change that adds code without tests fails locally and in CI:

```bash
just go-unit-cov-check   # Report coverage and fail below the target
```

`.github/codecov.yml` and the shared `go` justfile module declare the target. If
you change one, change the other.

### Test file conventions

- Public tests: `*_public_test.go` in the `_test` package of the package. They
  test the exported API. This is the default.
- Internal tests: `*_test.go` in the same package. They test what the exported
  API cannot reach.
- Suite names: `*_public_test.go` → `{Name}PublicTestSuite`, `*_test.go` →
  `{Name}TestSuite`.
- Use `testify/suite` with table-driven cases.
- Write one suite method for each function under test. Success, errors, and edge
  cases are rows in one table, not different methods.
- `export_test.go` gives external tests access to unexported symbols, through an
  alias or a setter. Do not use an alias to test again the behavior that the
  test of the caller already reaches. The pattern is for a helper with its own
  contract.

External tests in this repository are in `package orchestrator_test`. Their
tables have `validateFunc` callbacks.

## Adding a new operation

To add a new typed constructor (for example, `NodeRebootDo`), do these steps in
order. Every operation must have tests, docs, and an example.

### Step 1: operation constructor

Add the method to `pkg/orchestrator/ops.go`. Use the pattern of the other
methods:

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

Rules:

- Use `o.nextOpName("verb-noun")` for the step name.
- Always include `Status: r.Status` in the `engine.HostResult` mapper.
- Wrap errors with context: `fmt.Errorf("verb noun: %w", err)`.
- Use `engine.CollectionResult` for operations that target a host. It returns a
  result for each host. Use `engine.StructToMap` for responses that are not a
  collection.
- The `engine` import is `internal/engine`. Use it only in `pkg/orchestrator/`.
  External code must not use it.

### Step 2: tests

Update these two test files:

**`pkg/orchestrator/ops_test.go`** (internal, httptest pattern) tests the full
HTTP round-trip with a mock server:

- Create an `httptest.Server` that returns canned JSON.
- Call the constructor. Verify the result with `report.Decode()`.
- Write the success, error, and edge cases as table rows.

**`pkg/orchestrator/ops_public_test.go`** (public, step-creation pattern) tests
that the constructor creates valid steps:

- Verify that the step is not nil and has the expected task name.
- Write one suite method for each operation, with all the cases as rows.

Get 100% coverage on the two files.

### Step 3: operation doc

Create the doc next to the other docs for its domain. The path is
`docs/operations/{group}/{operation}.md` or
`docs/operations/{group}/{domain}/{operation}.md`, as the group layout tells.
Use the template of the docs next to it. Every doc must have these sections:

- **Description** (h1 heading with the method name)
- **Usage.** A small Go snippet that shows the constructor call.
- **Parameters.** A table of all the parameters, with types and descriptions.
- **Result type.** A `Decode()` snippet and a field table.
- **Idempotency.** One of: Read-only, Idempotent (Yes), Non-idempotent (No).
- **Permissions.** The OSAPI permission that the token must have (for example,
  `node:write`).
- **Example.** A link to the example file that uses this operation:
  ```
  Refer to
  [`examples/operations/{domain}.go`](https://github.com/osapi-io/osapi-orchestrator/blob/main/examples/operations/{domain}.go)
  for a complete example.
  ```

### Step 4: update the domain landing page and operation index

Add the operation to the table in the landing page of the directory that holds
the doc. Then update `docs/operations/README.md`. Change the count in the row of
the domain, and the total in the sentence above the group tables. The sum of the
rows must be equal to that total.

### Step 5: example

Add the operation to the workflow example in `examples/operations/` for the same
domain. The domain groups are:

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

If no example has the domain, create a new `{domain}.go` file. Every operation
must be in one or more runnable examples.

#### Example conventions

- **Self-contained.** Clean the host at the start (a different plan with
  `ContinueOnError()`), then do the work, then verify it. The example must give
  the same result each time it runs.
- **One purpose in each file.** Show the operations of one domain. Do not add
  other features (parallel, verbose, broadcast).
- **Cleanup plan pattern.** Use a different `orchestrator.New()` for the
  cleanup, with `ContinueOnError()`. Then use a main plan for the workflow.
- **Platform safety.** If an operation does not work on all platforms, use
  `ContinueOnError()`, so the example does not crash.
- **Decode and print.** Decode one or more results, so the example shows output.
  Use `report.Decode("step-name", &typedStruct)`.
- **Short.** Less than 100 lines of code (without the license). If it is longer,
  the example shows too much. Split it.
- **Operation docs link to examples.** Every operation doc in
  `docs/operations/{domain}/` must link to the example file that shows that
  operation.

### Step 6: update README.md

If the total number of operations changes, update the operation count and the
tables in the root `README.md`.

### Step 7: verify

```bash
go build ./...                                       # compiles
go test ./... -count=1                               # tests pass

# Each example is its own `package main`, so they are built one at a time
# rather than as a package.
cd examples/operations && for f in *.go; do go build -o /dev/null "$f"; done
cd examples/features   && for f in *.go; do go build -o /dev/null "$f"; done
```

## Before committing

Run `just ready` before you commit. It makes sure that the generated code, the
package docs, the format, and the lint are all current:

```bash
just ready
```

## Branching

Make all changes on feature branches. Create a branch from `main` with the name
`type/short-description`. The `type` is the [Conventional Commits] type:

- `feat/add-retry-logic`
- `fix/null-pointer-crash`
- `docs/update-api-reference`
- `refactor/simplify-handler`
- `chore/update-dependencies`

If you are on `main`, the `/commit` command of Claude Code creates a branch for
you.

## Commit messages

Follow [Conventional Commits] with the 50/72 rule:

- **Subject line.** No more than 50 characters, imperative mood, capital first
  letter, no period.
- **Body.** Wrap at 72 characters. Put a blank line between the subject and the
  body.
- **Format.** `type(scope): description`
- **Types.** `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`
- Tell "what" changed and "why", not "how"

Write commit messages that tell the reader something. Do not put too many
commits on a PR. Most PRs have one commit. A large PR can have a few commits.
Use git squash and rebase to combine commits.

## Submitting a PR

- **Describe your changes.** Tell what changed and why. A reviewer must not have
  to read the diff to learn the reason for the change.
- **Issue/PR links.** Link to the related issues or PRs. Tell how your changes
  are different from that work, or how they add to it.
- **Examples.** Add examples or screenshots that show the effect of your
  changes.
- **Draft PRs.** If your changes are not complete and you want to discuss them,
  open the PR as a draft. Then add a comment to start a discussion. Use
  comments, not the PR description, for the discussion. Then you can change the
  description later, and the discussion stays.

## AI usage

We write this repo with help from AI tools. The [AI Usage Policy](AI_POLICY.md)
applies to all contributions. Tell which tool you used. Make sure that you can
explain what your change does without help from AI tools.

## FAQ

> I want to contribute, where do I start?

We welcome all contributions, from a typo fix to a new feature. You can also
vote on issues, comment on issues, or help to answer questions.

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
