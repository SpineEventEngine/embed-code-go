---
name: go-engineer
description: >
  Encodes the embed-code-go Go implementation policy and recurring pitfalls.
  Use whenever writing, modifying, refactoring, debugging, or reviewing Go code
  in this repository, especially CLI/config parsing, parser state transitions,
  embedding behavior, fragmentation, filesystem handling, error reporting,
  package APIs, formatting, vetting, and build verification.
---

# Go - policy & pitfalls

Baseline Go knowledge is assumed. This skill does not teach the language; it
encodes the project policy, package boundaries, and traps that recur in `embed-code-go` work.

## When to Use

Use `go-engineer` for implementation work in Go:

- Writing or changing `.go` files.
- Debugging `embed` or `check` behavior.
- Refactoring package APIs or shared helpers.
- Reviewing Go code for correctness and maintainability.

This skill is the baseline for production Go and helper code. Test-writing
conventions live in `.agents/skills/go-tester/SKILL.md`; documentation review
lives in `.agents/skills/review-docs/SKILL.md`.

## Fast Path for Agents

1. Read `AGENTS.md`, `PROJECT.md`, and relevant implementation files.
2. Ask clarifying questions before editing if the requested outcome, scope,
   compatibility constraints, or verification target is not explicit.
3. Apply the MUST / MUST NOT rules while editing.
4. Defer test structure and fixtures to `go-tester` for test changes.
5. Verify with the narrowest relevant Go test first, then the repository-level
   checks listed below.
6. Do not commit, push, tag, merge, rebase, cherry-pick, or rewrite Git history.

## Setup Check

Run this before non-trivial Go changes or when the package baseline is unclear:

1. **Go version** - target the version in `go.mod`.
2. **Package owner** - identify whether the behavior belongs in `cli/`,
   `configuration/`, `embedding/`, `embedding/parsing/`, `fragmentation/`, or a
   support package.
3. **Test owner** - identify the package test suite and fixtures that already
   cover the behavior.
4. **Commands** - plan `gofmt`, focused `go test`, `go vet ./...`, full
   `go test ./...`, and `go build -trimpath main.go` when integration or CLI
   behavior changes.

## Processing Flow

The normal execution path is:

1. `main.go` reads arguments, configures logging, validates input, and dispatches
   `embed` or `check` mode.
2. `cli/` reads flags or YAML and produces one or more normalized
   `configuration.Configuration` values.
3. `embedding.EmbedAll` or `embedding.CheckUpToDate` selects documentation files
   using include and exclude patterns.
4. An `embedding.Processor` processes one document at a time.
5. `embedding/parsing/` walks the document through explicit states, records each
   instruction and its code fence, and preserves unrelated document content.
6. A parsed `Instruction` resolves source content through `fragmentation/`,
   optional line patterns, indentation normalization, and comment filtering.
7. Embed mode writes changed documents. Check mode compares generated content
   with existing content and must not modify documentation.

When behavior changes, trace the complete path instead of patching only the
first function that exposes the symptom.

## MUST DO

- **Preserve package ownership.** Keep process orchestration in `main.go`,
  argument and YAML handling in `cli/`, normalized defaults in
  `configuration/`, document orchestration in `embedding/`, syntax recognition
  in `embedding/parsing/`, and source extraction in `fragmentation/`.
- **Trace the full flow.** For behavior changes, walk
  `main` -> `cli`/`configuration` -> `embedding` -> `embedding/parsing` and
  `fragmentation` -> write or compare.
- **Keep functions small and explicit.** Prefer direct code over broad helpers
  unless an abstraction removes real duplication or clarifies a shared contract.
- **Document every function and method.** Add concise doc comments to new or
  changed functions, including unexported ones. Exported comments start with
  the declaration name.
- **Return actionable errors.** Add file, pattern, instruction, or operation
  context where that context becomes known.
- **Aggregate independent failures with `errors.Join`.** Continue processing
  only when doing so produces a more complete and still safe result.
- **Preserve deterministic behavior.** Document order, fragment order,
  separators, writes, and reported stale files must stay stable.
- **Use OS-aware paths.** Preserve cross-platform behavior for separators,
  Windows drive paths, and file URLs.
- **Format changed Go files with `gofmt`.**

## MUST NOT DO

- **No new panics in library packages.** Return errors; leave process exit and
  panic recovery at the executable boundary.
- **No broad utility dumping ground.** Do not move parser, filesystem,
  fragmentation, and CLI behavior into a catch-all helper package.
- **No single-use interface by default.** Introduce an interface only at a real
  consumer, ownership, or test boundary.
- **No speculative concurrency.** This CLI is document and filesystem heavy;
  add concurrency only for a measured problem with deterministic aggregation
  and write safety.
- **No log-and-return in lower layers.** Return the error; log at the CLI
  boundary or at intentional progress/warning points.
- **No message-string error contracts.** Prefer typed errors, sentinels, or
  `errors.Is` / `errors.As`.
- **No dependency updates unless the confirmed task requires them.**
- **No commits or history-writing.**

## Project Hotspots

### Parser State Machine

- Read `embedding/parsing/constants.go`, `state.go`, `context.go`, the affected
  state, and `embedding/processor.go` together.
- Preserve self-closing, paired, and multiline instruction forms.
- Keep ordinary Markdown fences separate from embedding fences.
- Preserve the opening fence's marker length and indentation when recognizing
  the closing fence and rendering source lines.
- Report malformed instructions at their start line with a concrete reason.
- Prefer concrete parse reasons: missing tag end, missing closing tag, invalid
  XML, missing code fence, or unclosed code fence.
- Do not consume unrelated later content to recover from invalid XML.

### Embedding And Check Modes

- Embed mode may rewrite only documents whose generated content changed.
- Check mode is read-only and reports stale documents.
- Shared processing changes should preserve both modes.

### Fragmentation

- Preserve whole-file, named-fragment, exact-pattern, and range-pattern extraction.
- Keep partition ordering, separator rendering, indentation normalization, and
  comment filtering stable.
- Keep source-root names and lookup errors actionable.

## Verification

Run the narrowest relevant checks first, then broaden:

1. `gofmt -w <changed-go-files>`
2. `go test ./<affected-package>`
3. `go vet ./...`
4. `go test ./...`
5. `go build -trimpath main.go` when CLI, configuration, embedding, or package
   integration behavior changed

For user-visible CLI behavior, run a focused `go run ./main.go ...` scenario
against existing fixtures when practical.

## Output Format

When producing code:

1. A short plan after clarification.
2. The code changes.
3. A verification checklist with command results.
4. A note about any remaining risk or unrun check.

When reviewing code: call out MUST-DO / MUST-NOT violations explicitly and
suggest the minimal fix. End with a one-line verdict: `APPROVE`,
`APPROVE WITH CHANGES`, or `REQUEST CHANGES`.
