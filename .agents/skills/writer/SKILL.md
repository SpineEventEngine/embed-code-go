---
name: writer
description: >
  Writes, edits, and restructures embed-code-go documentation. Use when asked
  to create or update README.md, showcase guides, PROJECT.md, AGENTS.md,
  skills, Markdown fixtures, contributor notes, examples, command snippets,
  Go doc comments, or inline explanatory comments. Verifies claims against
  current Go code, tests, fixtures, and project flows.
---

# Write documentation

## Decide the Target and Audience

- Identify the target reader: CLI user, contributor, maintainer, or agent.
- Identify the task type: new doc, update, restructure, or documentation audit.
- Identify the acceptance criteria: what is correct when the reader is done?
- Ask clarifying questions before editing if the audience, scope, or expected
  output file is unclear.

## Choose Where The Content Should Live

Prefer updating an existing document over creating a new one.

- `README.md`: user-facing overview, short run/build instructions, and links
  to the complete guide.
- `showcase/README.md`: complete user guide entry point and runnable workflow.
- `showcase/configuration/README.md`: command-line flags, YAML configuration,
  source roots, include/exclude patterns, and multiple embedding targets.
- `showcase/embedding/README.md`: embedding instruction map; detailed syntax
  belongs in the related `showcase/embedding/positive/*.md` feature page.
- `PROJECT.md`: project overview, project map, and compact doc pointers.
- `AGENTS.md`: agent operating policy and repository-wide rules.
- `.agents/skills/<name>/SKILL.md`: language, testing, writing, review, project
  flow, and implementation constraints.
- Go doc comments: API or helper behavior.
- `test/resources/`: parser and embedding fixtures.

## Verify Against Project Flows

- For CLI flags and modes, check `main.go`, `cli/cli.go`, and the
  `cli/cli_validation.go` file.
- For YAML configuration, check `cli/`, `configuration/`, and config fixtures
  under `test/resources/config_files/`.
- For embedding syntax, check `embedding/parsing/`, `embedding/processor.go`,
  `embedding/parsing/instruction.go`, and `showcase/embedding/`.
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
- Do not duplicate long explanations between `README.md`, the showcase, and
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

Follow the git-history policy in `AGENTS.md`.
