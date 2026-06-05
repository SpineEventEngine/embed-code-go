# Agent Instructions

These instructions apply to the whole repository. They are the primary operating
contract for every agent working in `embed-code-go`.

## Non-Negotiable Workflow

1. Read this file and every applicable skill under `skills/` before starting.
2. Begin every new task with a clarification round. Ask questions and wait for
   the answers before editing files, formatting code, running generators,
   installing dependencies, or otherwise changing the workspace.
3. Read-only inspection is allowed before the clarification round only when it
   is needed to ask precise questions.
4. Establish crystal-clear agreement on:
   - the desired outcome and acceptance criteria;
   - the files, packages, and behavior that are in scope;
   - compatibility constraints and behavior that must not change;
   - expected tests and documentation updates;
   - any ambiguity discovered during repository inspection.
5. If requirements conflict, remain incomplete, or depend on an unconfirmed
   assumption, stop and ask. Do not guess.
6. After the user confirms the scope, state the plan and the invariants that
   will be preserved. Then implement, verify, and report the result.
7. Follow these rules even when a faster shortcut appears harmless.

## Git Safety

- Never create commits in this repository.
- Never push, tag, merge, rebase, cherry-pick, or rewrite history.
- Do not run destructive Git commands that discard local work.
- Leave completed changes in the working tree and report the changed files and
  verification results to the user.
- Treat existing uncommitted changes as user work. Preserve them and work with
  them when they overlap the task.

## Skill Selection

- Use `skills/go-engineer/SKILL.md` for Go implementation, debugging,
  refactoring, API changes, and test changes.
- Use `skills/go-code-reviewer/SKILL.md` for reviews of Go code or repository
  changes. A review is read-only unless the user starts a separate fix task.
- Use `skills/project-documenter/SKILL.md` for architecture mapping,
  onboarding documentation, and updates to this file, `README.md`, or `EMBEDDING.md`.
- Apply more than one skill when a task crosses those boundaries.

## Project Mission

`embed-code-go` is a Go 1.22.1 command-line application. It scans Markdown and
HTML documents for `<embed-code>` instructions, resolves source fragments,
renders them inside code fences, and checks whether existing snippets are up-to-date.

The architecture follows a one-way flow from process orchestration to document
processing and source resolution. Keep syntax recognition, source extraction,
filesystem handling, and CLI concerns within their owning packages. Do not
collapse them into a broad utility layer or introduce circular dependencies.

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

## Package Map

- `main.go`: executable entry point, mode dispatch, aggregate errors, and final
  user-facing output.
- `cli/`: flag parsing, YAML loading, validation, and conversion to runtime
  configuration.
- `configuration/`: defaults and normalized settings used by processing code.
- `embedding/`: document discovery, per-document processing, embedding writes,
  and up-to-date checks.
- `embedding/parsing/`: state-machine parser for instructions, ordinary
  Markdown fences, embedding fences, and rendered-content comparison.
- `embedding/commentfilter/`: language-aware filtering for comment-retention
  modes.
- `fragmentation/`: whole-file, named-fragment, and line-pattern extraction;
  source lookup; partition assembly; encoding checks; and caches.
- `files/`: filesystem validation and file helpers.
- `indent/`: shared indentation measurement and removal.
- `logging/`: `slog` handler, clickable file references, panic reporting, and
  terminal formatting support.
- `type/`: YAML-compatible string and named-path list types. Preserve this
  existing package name even though it is unusual.
- `test/resources/`: parser, embedding, configuration, and source-code fixtures.

## Go Engineering Rules

- Prefer small, explicit functions and the repository's existing package
  boundaries over new layers or generic helpers.
- Keep the happy path clear with early returns where they reduce nesting.
- Accept interfaces where behavior is consumed and return concrete types by
  default. Do not introduce an interface for a single implementation unless it
  creates a real test or ownership boundary.
- Avoid global mutable state. Existing caches must remain bounded and have
  explicit lifecycle behavior.
- Do not introduce concurrency unless it solves a measured problem. Document
  ordering, writes, and error aggregation must remain deterministic.
- Use `filepath` for operating-system paths. Preserve cross-platform behavior
  for Windows drive paths, separators, and file URLs.
- Close files and other resources promptly, and preserve meaningful close or
  flush errors when they can affect correctness.
- Keep library packages free of new panics. Return errors; reserve process exit
  and panic recovery for the executable boundary.

## Errors And Logging

- Return errors from library packages and add context at the layer that knows
  the failing operation, file, pattern, or instruction line.
- Wrap underlying errors with `%w` when callers may inspect them with
  `errors.Is` or `errors.As`.
- Use `%v` only when rendering terminal-facing messages that are not returned
  for further inspection.
- Do not compare errors by message text when a type, sentinel, `errors.Is`, or
  `errors.As` can express the contract.
- Do not both log and return the same error in lower layers. Logging belongs at
  the CLI boundary or at intentional progress/warning points.
- Aggregate independent failures with `errors.Join` when continuing provides a
  more complete and still actionable result.

## Function Documentation

- Every function and method must have a doc comment, including unexported ones.
- Exported declarations must start with the exact exported name, following Go
  documentation style.
- Unexported comments should state intent in one concise sentence and start
  with the function name when natural.
- Document non-obvious side effects, filesystem writes, state transitions,
  errors, and panic behavior.
- Do not narrate obvious assignments or restate the signature.

## Parser And Embedding Invariants

- The parser is a state machine. Inspect `embedding/parsing/constants.go`,
  `state.go`, `context.go`, and the affected state together before changing it.
- Preserve these instruction forms:
  - `<embed-code file="..."/>`
  - `<embed-code file="..."></embed-code>`
  - multiline instructions already represented by fixtures.
- An instruction inside an unrelated Markdown code fence is ordinary content,
  not an active embedding instruction.
- Preserve the opening fence's marker length and indentation when recognizing
  the corresponding closing fence and rendering source lines.
- Malformed-instruction errors must point to the instruction start line, not
  EOF or a later fence line.
- Prefer concrete parse reasons: missing `>`, missing closing tag, invalid XML,
  missing code fence, or unclosed code fence.
- Do not scan arbitrary later document content to guess where an invalid
  instruction ends.
- Embed mode may rewrite only documents whose generated content changed.
- Check mode is read-only and must identify every stale document it can safely
  inspect.

## Testing Rules

- Use Ginkgo/Gomega `Describe` and `It` style in packages that already use it.
- Add a focused regression test near the package that owns the bug.
- Add or update fixtures under `test/resources/` for parser and end-to-end
  embedding behavior.
- Parser changes should cover successful input and relevant malformed input.
- Verify both `embed` and `check` semantics when shared processing behavior
  changes.
- Keep tests deterministic, cross-platform, and independent of machine-specific
  absolute paths unless the output intentionally contains one.
- Prefer behavioral assertions over assertions about incidental internal state.

## Commands

Run commands from the repository root unless a fixture requires another
working directory.

- Format changed Go files: `gofmt -w <files>`.
- Run focused tests while iterating: `go test ./embedding/...` or the affected
  package.
- Run static checks for Go changes: `go vet ./...`.
- Run the full suite before finishing: `go test ./...`.
- Build the executable when CLI or integration behavior changes:
  `go build -trimpath main.go`.
- Run from source: `go run ./main.go -mode <embed|check> ...`.

If a required command cannot be run, report the exact reason and the remaining
risk. Do not describe unrun checks as passing.

## Documentation Stewardship

- Keep this file aligned with actual code. Verify package and flow descriptions
  against source before editing them.
- Keep `README.md` user-oriented: installation, configuration, modes, and
  common operation.
- Keep `EMBEDDING.md` focused on instruction syntax, source markers, patterns,
  comment modes, and examples.
- Do not duplicate long user documentation in `AGENTS.md`; link to the owning
  document and retain only the engineering invariant here.
- When architecture changes, update this package map and processing flow in the
  same task after the user confirms documentation is in scope.

## Repository Hygiene

- Do not revert unrelated user changes.
- Do not edit generated binaries under `bin/` unless explicitly requested.
- Do not add temporary repro files, local binaries, IDE metadata, coverage
  output, or build artifacts to the intended change set.
- Do not update dependencies unless the confirmed task requires it.
- Keep changes narrowly scoped and make unrelated cleanup a separate task.
