---
name: go-tester
description: >
  The embed-code-go authority on writing and reviewing Go tests. Use when
  adding, restructuring, or reviewing tests for Go production code, parser
  behavior, embedding flows, fragmentation, CLI/config validation, fixtures,
  malformed inputs, regression coverage, focused go test commands, or full
  suite verification. go-engineer remains the baseline for Go inside test
  helper code.
---

# Go tester

This skill is the single source of truth for how tests are written in
`embed-code-go`. It does not decide what behavior to change; it decides how to
cover the known behavior once the cases are identified.

Two companions own neighboring concerns:

- `.agents/skills/go-engineer/SKILL.md` - Go implementation policy and production code correctness.
- `.agents/skills/review-docs/SKILL.md` - documentation review for test fixtures and
  Markdown docs when prose changes are part of the task.

## Core Policy

1. **Follow existing Ginkgo/Gomega style.** Packages use `Describe` and `It`
   suites with Gomega expectations. Match the surrounding package before adding
   a new style.
2. **Test behavior at the owning package.** Put the focused regression test near
   the package that owns the bug or feature.
3. **Use fixtures for document behavior.** Parser and embedding changes should
   add or update files under `test/resources/` when source or document shape is
   part of the contract.
4. **Cover both success and failure when parser behavior changes.** Malformed
   instruction tests are as important as successful embedding tests.
5. **Keep tests deterministic and portable.** Avoid machine-specific absolute
   paths, map iteration order, hidden local files, and OS-specific separators
   unless the code intentionally emits them.
6. **Do not commit or rewrite history.**

## Workflow

1. **Read first.** Read the code under test, its package tests, and relevant
   fixtures in full.
2. **Identify the owning suite.** Prefer extending an existing `Describe` block
   over creating a parallel test structure.
3. **Choose the level.** Use package-level unit tests for pure helpers, parser
   tests for state-machine behavior, and embedding tests for end-to-end
   document rewrites or stale checks.
4. **Add the minimal fixture.** Create or update only the source and document
   fixture needed to express the behavior.
5. **Assert behavior, not implementation trivia.** Prefer outputs, changed-file
   lists, typed errors, and line numbers over private state.
6. **Run the narrowest test first.** Broaden only after the focused test passes.

## Naming And Structure

- Keep test files beside the package they test, using the existing `_test.go`
  naming pattern.
- Match the package's current package name. Do not switch between internal and
  external test packages without a reason.
- Use `Describe` for the unit or component under test and `It` for the behavior sentence.
- Keep setup close to the test unless the package already has a clear helper pattern.
- Use `BeforeEach` only when it reduces duplication without hiding essential state.
- Avoid broad helper functions that make the fixture harder to read.

## Assertions

- Use Gomega matchers already used by the package.
- Prefer exact content assertions for rendered snippets when formatting is the contract.
- Prefer typed error checks with `errors.As` when the code exposes a typed error.
- Verify line numbers for malformed instructions, especially when the bug is an
  EOF or later-fence misreport.
- Avoid matching full error strings unless the user-facing message itself is
  the behavior under test.

## Fixture Guidance

### Documentation Fixtures

- Put documentation files under `test/resources/docs/`.
- Keep one behavior per fixture unless a multi-case document is already the
  package convention.
- Name malformed fixtures after the malformed condition, not the expected implementation fix.
- Preserve indentation and code-fence marker details when those are part of the behavior.

### Source Fixtures

- Put source files under the existing language folder in `test/resources/code/`.
- Add the smallest source file or fragment marker needed for the behavior.
- Reuse existing Java, Kotlin, plain-text, or literal-pattern fixtures when
  doing so keeps the test clearer.

### Prepared Fragments

- Use `test/resources/prepared-fragments/` only for behavior that already relies
  on prepared expected content.
- Keep expected files stable and easy to diff.

## Verification

Use the smallest useful command while iterating:

- Parser changes: `go test ./embedding/parsing`
- Embedding behavior: `go test ./embedding`
- Fragment extraction: `go test ./fragmentation`
- CLI/config behavior: `go test ./cli`
- Shared or uncertain behavior: `go test ./...`

After Go test code changes, run:

1. `gofmt -w <changed-go-files>`
2. the focused package test
3. `go test ./...`

Run `go vet ./...` when the test change also changes helper code, exported
test support, or production code.

## Repo Notes

- Tests use the repository fixtures heavily; prefer a focused fixture over a
  large inline string when document shape matters.
- Do not rely on local absolute paths except where the code intentionally emits
  a `file://` reference.
- If the implementation change is still unclear, return to `go-engineer`
  before writing tests.

## Output Format

When writing tests:

1. State the behavior being covered.
2. List fixtures added or changed.
3. Report focused and full test commands run.

When reviewing tests: group findings as `Must fix`, `Should fix`, and `Nits`.
End with `APPROVE`, `APPROVE WITH CHANGES`, or `REQUEST CHANGES`.
