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
- **No backward compatibility.** When a change makes something old or
  inconsistent, change it or delete it everywhere, in the same change. Do not
  add a deprecated flag, a field for old callers, or a shim for an old setup.
  For a change to the API, open an issue first. Then the maintainers can discuss
  it before you do the work.

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

- **`pkg/orchestrator/`.** The DSL that users call. It has a typed constructor
  for each operation, such as `HostnameGet` and `CommandExec`, and uses the SDK
  types directly, such as `osapi.HostnameEntry` and `osapi.AgentInfo`.
  - `orchestrator.go`. `New`, `Run`, `TaskFunc`.
  - `ops.go`. The step helpers and the controller operations.
  - `ops_hosts.go`. One constructor for each SDK call to the hosts.
  - `step.go`. Step chaining: guards, retry and ordering.
  - `result.go`. The result types: `Results`, `TaskStatus`, `Report`.
  - `host_status.go`. The `HostStatusOk`, `Skipped` and `Failed` constants.
  - `discover.go`. `Discover` and `GroupByFact`.
  - `predicate.go`. The agent predicates: OS, Arch, HasLabel and more.
  - `renderer.go`. The renderer interface.
  - `renderer_lipgloss.go`. The terminal output renderer.
  - `generate.go`. The mockgen directive, with no code.
  - `options.go`. The option types: verbose, upload and retry.
  - `types.go`. The `Orchestrator` struct, and aliases onto the engine.
- **`internal/engine/`.** The engine that runs a plan.
- **`docs/`.**
  - `operations/`. The operation reference, one doc for each operation.
    `README.md` is the index. A host operation is `{service}/{method}.md`, such
    as `service/start.md`. `agent/`, `health/` and `files/file/` are the
    controller operations.
  - `features/`. The feature guides, such as `guards.md`, `broadcast.md` and
    `retry.md`.
- **`examples/`.**
  - `operations/`. A runnable workflow example for each domain, such as
    `node-info.go`, `docker.go` and `schedule.go`.
  - `features/`. A runnable example for each feature, such as `host-status.go`,
    `guards.go` and `broadcast.go`.

A host operation has the name and the parameters of its osapi SDK call:
`ServiceStart` calls `c.Service.Start`, and its doc is at
`docs/operations/service/start.md`.

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

Do not write `//nolint:revive`. It hides a problem and does not fix it. Fix the
name or the code. When a name must be different in an API, give it an explicit
name, such as a struct tag.

Do not write an em dash, in a page or in a Go comment. The generated docs come
from the comments, so a comment is a page too.

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

A test entry point, such as `func TestServerPublicTestSuite(t *testing.T)`, is
the one exception. Write it on one line, as `go test` shows it. Each other
function with parameters has the signature on more than one line.

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

A test file has the name of the production file that it tests: `foo.go` has
`foo_public_test.go` and `foo_test.go`. `export_test.go` is the exception. If a
test file becomes too large to read, split the production file first. Then each
test file keeps its production file. `just go-test-name-check` fails a test file
that has no production file of the same name in its directory.

### Go patterns

- Do not call `panic` in code that is not a test. Return an error, also for a
  registration error or a programming error. `main` reports the error and exits
  with a code that is not zero. A `recover` that stops a panic from a dependency
  is permitted. The `forbidigo` linter fails a `panic` outside a test.
- Read and write files through [avfs], so that a test can use a file system in
  memory. This applies to tools that run at build time too: generators,
  scaffolders and the tools that write docs.
- Wrap errors with `fmt.Errorf("context: %w", err)`. Then the chain names each
  layer. `errors.Is` and `errors.As` can examine the chain.
- Return early. Do not put the successful path inside conditions.
- Give an unused parameter the name `_`.
- Put imports in three groups, with a blank line between them: the standard
  library, third-party packages, then local packages.
- A package holds one kind of thing, and its name tells what. Do not make a
  package named `util`, `common`, `helpers` or `shared`.
- The name of a package is the last part of its path. An exported name does not
  repeat it inside a longer name: `atomicfile.Write`, not
  `atomicfile.WriteFileAtomic`. The main type can have the name of the package,
  as `time.Time` does.
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
target. Change the two together. Put the Codecov settings only in
`.github/codecov.yml`. A `codecov.yml` at the root replaces it, and nothing
fails. The `threshold` of `.github/codecov.yml` absorbs the rounding of Codecov
only. `just go-unit-cov-check` is the exact gate at the target.

### Test file conventions

- Public tests: `*_public_test.go`, in package `<name>_test`. They test the
  exported API. Write this type of test first.
- Internal tests: `*_test.go`, in package `<name>`. Use them only for what the
  exported API cannot reach.
- A test in package `<name>_test` has a name that ends in `_public_test.go`. A
  `*_test.go` file is always an internal test.
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

A host operation wraps one call of the osapi SDK. When the SDK gains a call, add
its operation in these steps, in order. Each operation has tests, a doc, and an
example.

### Step 1: operation constructor

Add the method to `pkg/orchestrator/ops_hosts.go`, next to the other methods of
its SDK service. Give it the name of the service and the method, and the SDK's
own parameters after `target`:

```go
// ServiceStart creates a step that runs Service.Start on the target: start a
// service.
func (o *Orchestrator) ServiceStart(
    target string,
    name string,
) *Step {
    return collectionStep(
        o, "start-service", false,
        func(
            ctx context.Context,
            c *osapi.Client,
        ) (*osapi.Response[osapi.Collection[osapi.ServiceResult]], error) {
            return c.Service.Start(ctx, target, name)
        },
    )
}
```

Rules:

- The step name is the method and the service in kebab case, such as
  `start-service` or `create-interface-route`.
- `collectionStep` reads the hostname, status, changed and error fields of each
  row through its JSON. A row that needs more, such as the exit code of a
  command, uses `rowStep` with its own mapper.
- Pass `true` for an operation that is not idempotent, such as a restart or a
  reboot. `Retry` then does not run it again unless the step adds `RetryUnsafe`.

### Step 2: tests

Add a row for the operation to `TestHostOps` in
`pkg/orchestrator/ops_public_test.go`. The row gives the constructor, the HTTP
method that the SDK sends, the step name, and whether the row type has
`Changed`. The test runs the step against a fake controller and checks the
request and the report. Keep 100% coverage.

### Step 3: operation doc

Write `docs/operations/{service}/{method}.md` in the form of the docs next to
it:

- **The name** as the h1, then what the call does on each OS family. Take this
  from the doc comment of the SDK method.
- **Usage.** The constructor call.
- **Parameters.** `target`, then each parameter of the SDK call, with a field
  table for an input or a params type.
- **Result.** The row type, a `Decode()` snippet with the step name, and a field
  table.
- **Retries.** For an operation that is not idempotent.
- **Permissions.** The permission and the route, from the SDK's doc comment.
- **Example.** A link to the example that uses the operation.

Then add the operation to `docs/operations/README.md`, and change the count in
the sentence at the top.

### Step 4: example

Add the operation to the workflow example in `examples/operations/` for its
service, or make a new `{service}.go`. Every operation must be in one or more
runnable examples.

#### Example conventions

- **Self-contained.** Clean the host at the start, with a different plan that
  has `ContinueOnError()`. Then do the work. Then verify it. The example must
  give the same result each time it runs.
- **One purpose in each file.** Show the operations of one service. Do not add
  other features (parallel, verbose, broadcast).
- **Cleanup plan pattern.** Use a different `orchestrator.New()` for the
  cleanup, with `ContinueOnError()`. Then use a main plan for the workflow.
- **Platform safety.** If an operation does not work on all platforms, use
  `ContinueOnError()`, so the example does not crash.
- **Decode and print.** Decode one or more results, so the example shows output.
  Use `report.Decode("step-name", &row)` for one row, and
  `report.HostResults("step-name")` for each row of a list.
- **Short.** Less than 100 lines of code (without the license). If it is longer,
  the example shows too much. Split it.

### Step 5: update README.md

Update the operation count and the table of services in the root `README.md`.

### Step 6: verify

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

[avfs]: https://github.com/avfs/avfs
[conventional commits]: https://www.conventionalcommits.org
[discussion]: https://github.com/osapi-io/osapi-orchestrator/discussions
[go]: https://go.dev
[gofumpt]: https://github.com/mvdan/gofumpt
[golangci-lint]: https://golangci-lint.run
[just]: https://just.systems
[mdformat]: https://pypi.org/project/mdformat/
[mise]: https://mise.jdx.dev
[uv]: https://docs.astral.sh/uv/
