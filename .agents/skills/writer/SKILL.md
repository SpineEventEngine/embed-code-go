---
name: writer
description: >
  Writes, edits, and restructures embed-code-go documentation. Use when asked
  to create or update README.md, EMBEDDING.md, AGENTS.md, skills, Markdown
  fixtures, contributor notes, examples, command snippets, Go doc comments, or
  inline explanatory comments. Verifies documentation claims against current Go
  code, tests, fixtures, and project flows.
---

# Write documentation (repo-specific)

## Decide the Target and Audience

- Identify the target reader: CLI user, contributor, maintainer, or agent.
- Identify the task type: new doc, update, restructure, or documentation audit.
- Identify the acceptance criteria: what is correct when the reader is done?
- Ask clarifying questions before editing if the audience, scope, or expected
  output file is unclear.

## Choose Where The Content Should Live

- Prefer updating an existing document over creating a new one.
- Put project entry-point and usage content in `README.md`.
- Put embedding syntax, source markers, patterns, fences, separators, comment
  modes, and examples in `EMBEDDING.md`.
- Put project flow, package ownership, skill routing, and repository-level
  invariants in `AGENTS.md`.
- Put language, testing, writing, and review procedures in the relevant
  `.agents/skills/<name>/SKILL.md`.
- Put API or helper behavior that belongs with code in Go doc comments.
- Put parser and embedding fixtures under `test/resources/` only when they are
  part of tests or examples already in scope.

## Verify Against Project Flows

- For CLI flags and modes, check `main.go`, `cli/cli.go`, and
  `cli/cli_validation.go`.
- For YAML configuration, check `cli/`, `configuration/`, and config fixtures
  under `test/resources/config_files/`.
- For embedding syntax, check `embedding/parsing/`, `embedding/processor.go`,
  `embedding/parsing/instruction.go`, and `EMBEDDING.md`.
- For source fragments and patterns, check `fragmentation/`,
  `embedding/parsing/pattern.go`, and relevant source fixtures.
- For comments modes, check `embedding/commentfilter/`.
- For examples that claim command output or write behavior, verify with tests or
  a focused `go run ./main.go ...` command when practical.

## Follow Local Documentation Conventions

- Use fenced code blocks for commands, YAML, Markdown, and source examples.
- Render file paths, package paths, flags, config keys, instruction attributes,
  function names, and command names as code.
- Keep headings hierarchical: one top-level `#`, then ordered levels without skips.
- Preserve the existing link style of the file being edited. Use local relative
  links for repository files. Prefer reference-style links for long external
  URLs when they improve readability.
- Keep examples small enough to verify and copy.
- Use consistent terminology: embed mode, check mode, embedding instruction,
  code fence, fragment, source root, docs root, include pattern, exclude pattern.
- Do not leave orphans in prose: no paragraph, list item, or table cell should
  end with a final line containing only one word. Reflow or rewrite the text.
- Do not duplicate long explanations between `README.md`, `EMBEDDING.md`, and
  `AGENTS.md`; link to the owning document instead.

## Go Doc Comment Guidance

- Every new or changed function and method should have a useful doc comment in
  this project, including unexported functions.
- Exported comments start with the exact declaration name.
- Unexported comments state intent and start with the function name when it
  reads naturally.
- Document non-obvious state transitions, filesystem writes, returned errors,
  panics, and parser constraints.
- Do not restate the signature or narrate obvious assignments.
- Inline comments in production Go should explain why a constraint exists, not
  what the next line does.

## Make Docs Actionable

- Prefer executable steps, expected outcomes, and concrete examples over broad descriptions.
- Include prerequisites such as Go version, working directory, fixture location,
  or mode when they are easy to miss.
- When documenting failure behavior, include the concrete reason and where the
  user should look.
- When documenting architecture, describe ownership boundaries and the normal
  flow rather than every helper function.

## Validate Changes

- Verify every referenced path exists.
- Verify flags, config keys, defaults, and instruction attributes against code.
- Verify Markdown examples and local links.
- Run `go test ./...` when documentation changes depend on behavior described by
  tests or fixtures.
- Run a focused `go run ./main.go ...` scenario when adding or changing a CLI example.

## Output Format

When writing documentation:

1. State the target audience and file location.
2. Summarize the documentation changed.
3. List source files, tests, or fixtures used to verify claims.
4. Report validation commands run and any remaining unverified claims.

Never commit, push, tag, or rewrite Git history.
