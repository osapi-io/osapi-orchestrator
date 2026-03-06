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

## Quick Reference

```bash
just fetch / just deps / just test / just go::unit / just go::vet / just go::fmt
```

## Package Structure

- **`pkg/orchestrator/`** — User-facing DSL
  - Typed operation constructors (NodeHostnameGet, CommandExec, etc.)
  - Typed result structs (HostnameResult, DiskResult, etc.)
  - Porcelain over osapi-sdk's orchestrator engine

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
- One suite method per function under test — all scenarios (success, errors, edge cases) as rows in one table

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

### Task Tracking

Implementation planning and execution uses the superpowers plugin workflows
(`writing-plans` and `executing-plans`). Plans live in `docs/plans/`.

### Commit Messages

See @docs/development.md#commit-messages for full conventions.

Follow [Conventional Commits](https://www.conventionalcommits.org/) with the
50/72 rule. Format: `type(scope): description`.

When committing via Claude Code, end with:
- `🤖 Generated with [Claude Code](https://claude.ai/code)`
- `Co-Authored-By: Claude <noreply@anthropic.com>`

[OSAPI]: https://github.com/osapi-io/osapi
[osapi-sdk]: https://github.com/osapi-io/osapi-sdk
