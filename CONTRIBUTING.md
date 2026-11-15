# Contributing

Thank you for your contribution to osapi-orchestrator.

## Before you start

- Read the [Code of Conduct](CODE_OF_CONDUCT.md). It applies to all
  communication in this repository.
- **Design pages.** The design pages of this repository are in
  [osapi-io/specs](https://github.com/osapi-io/specs), under
  `components/osapi-orchestrator/`. Read them before you change the behavior.
  When you change the behavior, change the design page in the same work. A
  design page in this repository becomes incorrect when the code changes. No
  check finds the difference.
- **Examine the work that exists.** Look for a pull request or an issue about
  the change. If you find one, use what it says in your work.
- **Backward compatibility.** Make sure that your change does not stop a program
  that uses osapi-orchestrator. A change that keeps compatibility is more likely
  to merge. If your change cannot keep it, open an issue first. Then the
  maintainers can discuss the API change before you do the work.

## Prerequisites

Install the tools with [mise]:

```bash
mise install
```

- **[Go].** osapi-orchestrator is written in Go. The project supports the two
  latest major versions of Go. Use one of them.
- **[uv].** A Python package runner. `just md-fmt` uses it to run [mdformat]. It
  installs nothing into the repository.
- **[just].** The task runner. Use it to build, test and format. You can also
  install it with `brew install just`.

### Coding agents

This repository works with any coding agent that reads `AGENTS.md`. `CLAUDE.md`
only points to `AGENTS.md`. Prepare your agent as its documentation tells you.
This repository does not need a plugin or an extension.

## Setup

Download the shared justfiles and install the dependencies:

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

[gofumpt] formats the Go code. [golangci-lint] examines the Go code for
problems. CI makes sure that both pass.

```bash
just go-fmt-check   # Check formatting
just go-fmt         # Auto-fix formatting
just go-vet         # Run linter
```

`.golangci.yml` lists the linters that run. Read the list there. A copy of the
list in a page becomes incorrect when the configuration changes. Formatting does
not apply to generated files (`*.gen.go`, `*.pb.go`).

### Documentation

[mdformat] formats the markdown files, through `uvx`. CI makes sure that the
format is correct.

```bash
just md-fmt-check   # Check formatting
just md-fmt         # Auto-fix formatting
```

Write each page in Simplified Technical English, as
[the osapi-io constitution](https://github.com/osapi-io/specs/blob/main/CONSTITUTION.md)
requires.

## Code standards

### Function signatures

If a function has parameters, write its signature on more than one line. Put
each parameter on its own line. Put the last `)` and the return types on their
own line:

```go
func FunctionName(
    param1 type1,
    param2 type2,
) (returnType, error) {
}
```

Write a function that has no parameters on one line:

```go
func Name() string {
}
```

Then a new parameter shows as one added line in a diff, not as a changed
signature.

### File names

Give a file the name of what it holds. Do not use `helpers.go`, `utils.go` or
names of that type. They tell where code is, not what it is. They also collect
code that has no other place.

`types.go` holds only type declarations: structs, interfaces, constants and
aliases. Put a function in a file that has the name of what the function does.

A test file has the name of the production file that it tests. If a test file
becomes too large to read, split the production file first. Then each test file
keeps its production file.

### Go patterns

- Wrap errors with `fmt.Errorf("context: %w", err)`. Then the chain names each
  layer. `errors.Is` and `errors.As` can examine the chain.
- Return early. Do not put the successful path inside conditions.
- Give an unused parameter the name `_`.
- Put imports in three groups, with a blank line between them: the standard
  library, third-party packages, then local packages.
- A package holds one kind of thing, and its name tells what. Do not make a
  package named `util`, `common`, `helpers` or `shared`.
- The name of a package is the last part of its path. An exported name does not
  repeat it: `atomicfile.Write`, not `atomicfile.WriteFileAtomic`.
- A caller does not need an alias. If two packages that a file imports have the
  same name, rename one of them. Mock packages are the exception. Each is named
  `mocks`, so a test that uses two of them gives one an alias. A clash with a
  package outside the module, or between two names that are both correct, also
  takes an alias.
- A directory can group packages that belong together. Do not put many packages
  under one package.
  [A Go package is one thing](https://github.com/osapi-io/specs/blob/main/CONSTITUTION.md#a-go-package-is-one-thing)
  tells why.

### Test doubles

For an interface that this organization defines, generate the double with
`mockgen` and commit it. Do not write a struct by hand for it.

A generated mock is in a `mocks` package next to the code that it mocks. A
`generate.go` file in that package holds the directive:

```go
package mocks

//go:generate go tool go.uber.org/mock/mockgen -source=../types.go -destination=types.gen.go -package=mocks
```

The tool dependencies of the module supply the generator. Thus each checkout
runs the version that `go.mod` records. The names of the generated files end in
`.gen.go`. Commit these files. Do not put mocks in `gen/`. That name is for API
code generation.

If the interface is **unexported**, a sibling package cannot hold the mock. The
mock must import the package to name the types of the interface. The tests of
the package must import the mock. Thus generate the mock into the package. Give
the mock a destination that only tests use. Then the mock library is not a
dependency of code that imports the package:

```go
// generate.go, in the package that declares the interface
package thispackage

//go:generate go tool go.uber.org/mock/mockgen -source=thing.go -destination=thing.gen_test.go -package=thispackage
```

In both cases, the directives are in a `generate.go` that holds no code. The
name of the generated file has `.gen`, so a reader knows not to edit it.

If many tests set the same expectations, write a constructor that returns a
configured mock. Do not write a type by hand for it. The generated mock is still
the type that satisfies the interface.

Write these three doubles by hand, because a generated one gives no advantage:

- A double for a standard library interface: `net.Conn`, `fs.File`, `io.Writer`,
  `slog.Handler`. These interfaces do not change when our code changes.
- A double that has a real implementation of the behavior under test, for
  example signatures from a key pair that the test generates.
- A recorder for a dependency that a goroutine calls, if the test cannot wait
  for that goroutine. A generated mock counts the calls at a time that the test
  cannot control. Write that reason where you define the recorder.

### File headers

Start each `.go` file with the MIT license header. Copy the format from a Go
file in the repository. In a file with a build tag, put `//go:build` on line 1,
then a blank line, then the header.

## Tests

```bash
just test           # Run all tests (lint + unit + coverage)
just go-unit       # Run unit tests only
just go-unit-cov   # Generate coverage report
go test -run TestName -v ./pkg/orchestrator/...  # Run a single test
```

The coverage target is 100%. `just test` fails when the total coverage is less
than the target. Thus a change that adds code without tests fails on your
machine and in CI:

```bash
just go-unit-cov-check   # Report coverage and fail below the target
```

`.github/codecov.yml` and the shared `go` justfile module both declare the
target. Change the two together.

### Test file conventions

- Public tests: `*_public_test.go`, in the `_test` package of the package. They
  test the exported API. Write this type of test first.
- Internal tests: `*_test.go`, in the same package. Use them only for what the
  exported API cannot reach.
- Suite names: `*_public_test.go` has `{Name}PublicTestSuite`, and `*_test.go`
  has `{Name}TestSuite`.
- Use `testify/suite` with table-driven cases.
- Write one suite method for each function under test. Success, errors and edge
  cases are rows in one table, not different methods.
- `export_test.go` gives external tests access to unexported symbols, through an
  alias or a setter. Do not use an alias to test again what the test of the
  caller already reaches. Use it for a helper that has its own contract.

The external tests of this repository are in `package orchestrator_test`. Their
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
`docs/operations/{group}/{domain}/{operation}.md`, as the layout of the group
shows. Use the template of the docs next to it. Every doc must have these
sections:

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

- **Self-contained.** Clean the host at the start, with a different plan that
  has `ContinueOnError()`. Then do the work. Then verify it. The example must
  give the same result each time it runs.
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

## Before you commit

Run `just ready` before you commit. It makes sure that generated code, package
documentation, format and lint are current:

```bash
just ready
```

## Branches

Make all changes on a branch. Make the branch from `main`. Give it the name
`type/short-description`, where `type` is the [Conventional Commits] type:

- `feat/add-retry-logic`
- `fix/null-pointer-crash`
- `docs/update-api-reference`
- `refactor/simplify-handler`
- `chore/update-dependencies`

## Commit messages

Use [Conventional Commits] with the 50/72 rule:

- **Subject line.** 50 characters or fewer, in the imperative mood, with no
  period at the end.
- **Body.** Wrap it at 72 characters. Put a blank line between the subject and
  the body.
- **Format.** `type(scope): description`.
- **Types.** `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`,
  `chore`.
- Tell what changed and why, not how.

Write commit messages that help the reader. Keep the number of commits in a pull
request small. Most pull requests have one commit. A large pull request can have
a few. Use squash and rebase to combine commits.

## Submit a pull request

- **Describe your changes.** Tell what changed and why. Then a reviewer does not
  have to read the diff to find the reason.
- **Links to issues and pull requests.** Link to related issues and pull
  requests. Tell how your change is different from that work, or how it adds to
  it.
- **Examples.** Add examples or screenshots that show the effect of your change.
- **Draft pull requests.** If your change is not complete but you want to
  discuss it, open the pull request as a draft. Start the discussion in a
  comment, not in the description. Then you can change the description later and
  keep the discussion.

## AI usage

We write this repository with the help of AI tools. The
[AI Usage Policy](AI_POLICY.md) applies to all contributions. Tell which tool
you used. Make sure that you can explain what your change does without the help
of an AI tool.

## FAQ

> I want to contribute. Where do I start?

All contributions are welcome, from a typo fix to a new feature. You can also
vote or comment on issues, or help answer questions.

> I need help. Where do I ask?

Open a [Discussion] on GitHub.

[conventional commits]: https://www.conventionalcommits.org
[discussion]: https://github.com/osapi-io/osapi-orchestrator/discussions
[go]: https://go.dev
[gofumpt]: https://github.com/mvdan/gofumpt
[golangci-lint]: https://golangci-lint.run
[just]: https://just.systems
[mdformat]: https://pypi.org/project/mdformat/
[mise]: https://mise.jdx.dev
[uv]: https://docs.astral.sh/uv/
