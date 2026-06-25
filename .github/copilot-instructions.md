# GitHub Copilot Instructions

## Repository Context

This repository is **embed-code-go**, a TeamDev Go command-line application.
It scans Markdown and HTML documents for `<embed-code>` instructions, resolves
source fragments, renders them inside code fences, and checks whether existing
snippets are up-to-date.

Read [`AGENTS.md`](../AGENTS.md) first. It contains repository-wide agent
rules, safety constraints, and Git history policy.

For project-specific context, also consult:

- [`PROJECT.md`](../PROJECT.md) for the project overview, package map,
  documentation ownership, and CI notes.
- [`README.md`](../README.md) for the short user entry point, run commands,
  build command, and links to the complete guide.
- [`showcase/README.md`](../showcase/README.md) before changing user-facing
  examples, embedding behavior, configuration examples, or showcase tests.
- [`showcase/configuration/README.md`](../showcase/configuration/README.md)
  before changing CLI flags, YAML configuration, source roots, include/exclude
  patterns, or multiple embedding targets.
- [`showcase/embedding/README.md`](../showcase/embedding/README.md) before
  changing `<embed-code>` instruction syntax, source selection, fragments,
  patterns, comment filtering, or rendered examples.

## Project Shape

- Language: Go module `embed-code/embed-code-go`, using Go `1.26.4` from
  [`go.mod`](../go.mod).
- Entry point: [`main.go`](../main.go), including mode dispatch, aggregate
  errors, version embedding from [`VERSION`](../VERSION), and final
  user-facing output.
- CLI and config: `cli/` parses flags, reads YAML, validates settings, and
  converts user input into runtime configuration from `configuration/`.
- Embedding flow: `embedding/` discovers documents, processes each document,
  writes embedded snippets, and checks rendered snippets for staleness.
- Parser: `embedding/parsing/` is a state-machine parser for instructions,
  ordinary Markdown fences, embedding fences, and rendered-content comparison.
- Comment filtering: `embedding/commentfilter/` applies language-aware
  comment-retention filtering for embedded snippets.
- Source resolution: `fragmentation/` extracts whole files, named fragments,
  line patterns, partitioned fragments, source-root lookups, encoding checks,
  and caches.
- Filesystem, indentation, logging, and YAML helper types live in `files/`,
  `indent/`, `logging/`, and `type/`. The import path segment is `type`, but
  the Go package identifier is `_type`.
- Parser, embedding, configuration, and source-code fixtures live in
  `test/resources/`.
- User guide and executable end-to-end examples live in `showcase/`.
- Repository-specific Codex skills live under `.agents/skills/`; they are part
  of the repository's agent workflow and should stay in sync with code and docs.

Do not assume Gradle, npm, frontend, infrastructure, Spine/config helper
scripts, shared `.agents` submodules, or config-managed skip rules exist in
this repository.

## Implementation Notes

- Keep parser changes local to the existing state-machine flow in
  `embedding/parsing/` unless the task explicitly requires a broader rewrite.
- Preserve user-facing error wording, exit behavior, and aggregate error shape
  unless tests and documentation are updated with the behavior change.
- Use `filepath` for filesystem paths and keep path behavior portable across
  Windows, macOS, and Linux.
- Keep embedding writes constrained to configured documentation roots and files
  selected by `doc-includes` and `doc-excludes`.
- Prefer existing helpers from `files/`, `indent/`, `logging/`,
  `fragmentation/`, and `configuration/` over new ad hoc utilities.

## Testing Style

- Tests use Ginkgo and Gomega. Follow the existing `Describe`, `Context`, `It`,
  `DescribeTable`, and matcher style in nearby tests.
- Parser, pattern, and instruction changes should usually add focused coverage
  in `embedding/parsing/` and fixtures under `test/resources/` when needed.
- Fragment extraction, source-root lookup, partition, cache, and encoding
  changes should usually add focused coverage in `fragmentation/`.
- CLI flag, YAML, validation, and runtime config changes should usually add
  focused coverage in `cli/` and update configuration fixtures or showcase
  examples when behavior changes.
- User-facing embedding workflows, rendered snippets, and documentation examples
  should be proven by the showcase tests or by a focused `go run ./main.go`
  command against a real showcase config.

## Documentation Rules

- Keep [`README.md`](../README.md) short: it is the user-facing entry point with
  links to the complete guide.
- Keep complete usage details in `showcase/`, especially configuration and
  embedding syntax details.
- Changes to public CLI flags, YAML keys, defaults, instruction attributes,
  comment modes, or mode behavior must update the owning showcase documentation.
- Keep [`PROJECT.md`](../PROJECT.md) focused on project map, ownership, and CI
  notes. Keep [`AGENTS.md`](../AGENTS.md) focused on repository-wide agent
  policy.
- Verify embedded examples through `embed` mode, `check` mode, or tests instead
  of editing rendered snippets by hand.

## Review Focus

Prioritize:

- Correctness bugs and behavioral regressions in CLI validation, parser state
  transitions, embedding writes, check mode, and source-fragment resolution.
- Missing tests for changed parsing, configuration, embedding, fragmentation,
  comment filtering, filesystem, logging, or showcase behavior.
- Cross-platform path, newline, file permission, encoding, and temporary-file
  issues. CI runs on Windows, macOS, and Ubuntu.
- Error reporting quality, especially user-facing aggregate errors and clickable
  file references from `logging/`.
- Documentation drift between code behavior, `README.md`, `PROJECT.md`, and
  the showcase guides.
- Security and safety issues such as unsafe filesystem traversal, surprising
  writes outside configured docs roots, secret leakage in logs, or broad global
  state.

For reviews, lead with findings ordered by severity and include precise
file/line references.

## Do Not Review

Skip routine review of generated, vendored, or tool-owned files when present:

- Local binaries such as `embed-code` or `embed-code.exe`
- Coverage output and temporary files
- IDE metadata such as `.idea/**`
- Go build cache, module cache, and other local tool artifacts

Do review repository-owned workflow, documentation, agent, and configuration
files, including `AGENTS.md`, `PROJECT.md`, `.github/copilot-instructions.md`,
`.github/workflows/**`, `.golangci.yml`, and `.agents/skills/**`.

## Output Hygiene

- Do not add local binaries, coverage files, temporary outputs, local caches,
  IDE metadata, or machine-specific files to the intended change set.
- Do not commit generated showcase output or rewritten snippets unless the task
  explicitly changes rendered documentation examples.
- Do not hand-edit rendered snippets that should be produced by `embed` mode.
- Keep unrelated cleanup separate from the requested change.

## Universal Rules

Do not suggest:

- Git history operations such as `git commit`, `git push`, `git tag`,
  `git rebase`, `git merge`, `git cherry-pick`, or `gh pr merge` unless the
  current user request explicitly asks for them.
- Auto-updating dependency versions outside a dedicated dependency update task.
- Changing public CLI flags, YAML keys, defaults, output wording, or exit
  behavior without updating tests and documentation.
- Rewriting showcase output by hand when `embed` mode or tests can prove the
  expected rendered snippets.
- Broad parser rewrites when a local state transition or fixture can cover the
  behavior.
- Reflection, unsafe code, hidden background work, or broad global state
  without explicit approval.
- Committing secrets, credentials, tokens, customer data, local absolute paths,
  or machine-specific artifacts.

## Verification Hints

Prefer the smallest verification command that proves the change.

Format changed Go files:

```bash
gofmt -w <files>
```

Format imports when `goimports` is available:

```bash
goimports -w <files>
```

Run focused package tests when they prove the change:

```bash
go test -v ./cli -p 1
go test -v ./embedding/parsing ./fragmentation -p 1
```

Run the normal Go test suite sequentially, matching CI behavior:

```bash
go test -v ./... -p 1
```

Run the executable showcase tests:

```bash
go test -v -tags showcase ./showcase -p 1
```

Run a focused CLI check when changing embedding or configuration examples:

```bash
go run ./main.go -mode=check -config-path=showcase/embedding/embed-code.yml
```

Linting is configured by [`.golangci.yml`](../.golangci.yml) and CI uses
`golangci-lint-action` version `v8` with `golangci-lint` version `v2.12.2`.
Run the local equivalent only when available:

```bash
golangci-lint run ./...
```
