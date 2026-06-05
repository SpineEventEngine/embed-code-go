---
name: project-documenter
description: >
  Inspect embed-code-go and create or update source-grounded architecture,
  onboarding, package-map, workflow, and agent documentation. Use when asked to
  explain the grand plan, map repository components, document where to look,
  refresh AGENTS.md, or align README.md and EMBEDDING.md with code. Requires a
  clarification round and writes documentation instead of returning only a
  conversational summary.
---

# Project Documenter

Document `embed-code-go` from evidence in the repository. Treat source, tests,
fixtures, and build configuration as authoritative; treat existing prose as a
hypothesis that must be checked.

## Mandatory Clarification Gate

1. Read `AGENTS.md` and inspect only enough repository structure to ask precise
   questions.
2. Ask the user to confirm:
   - the audience: contributor, maintainer, user, or agent;
   - the required output files;
   - the desired depth and whether diagrams are wanted;
   - whether stale existing prose may be rewritten or only extended.
3. Wait for answers before editing documentation or running generators.
4. If documentation would assert an architectural intention not demonstrated by
   code, ask the user instead of inventing it.

## Required Outcome

When this skill is invoked to document the repository, do not finish with only
a chat summary. Write or update the user-approved documentation files.

`AGENTS.md` is the primary maintainer and agent map. Unless the user explicitly
excludes it, update it when repository-wide architecture, workflow, package
ownership, testing policy, or agent rules are being documented.

## Evidence Collection

After clarification:

1. Read `go.mod`, `main.go`, `README.md`, and `EMBEDDING.md`.
2. Enumerate packages and non-test Go files.
3. Read package entry points, exported declarations, and central data types.
4. Trace the two user journeys:
   - embed mode from flags/YAML to document writes;
   - check mode from flags/YAML to stale-file reporting without writes.
5. Trace parser state transitions and fragment-resolution paths.
6. Read package tests and representative fixtures to discover contracts and
   malformed-input behavior.
7. Inspect CI workflows and documented build commands.
8. Check the working tree so user changes are preserved.

Use searches and symbol lists to establish coverage, then read the important
files. Do not infer the architecture from directory names alone.

## Documentation Ownership

### `AGENTS.md`

Keep repository-wide engineering guidance here:

- mandatory working protocol and Git safety;
- skill selection;
- mission and architectural direction;
- end-to-end processing flow;
- package ownership and dependency boundaries;
- parser, embedding, testing, and documentation invariants;
- canonical development commands.

Keep it operational and concise enough to read before every task.

### `README.md`

Keep user-facing material here:

- what the CLI does;
- installation and build instructions;
- modes, flags, and YAML configuration;
- ordinary examples and expected output.

Do not place internal state-machine details here unless users need them.

### `EMBEDDING.md`

Keep the embedding language here:

- supported `<embed-code>` forms and attributes;
- source fragment markers and line patterns;
- fences, indentation, separators, and comment modes;
- valid and invalid examples.

### Additional Architecture Documents

Create additional files only when the user confirms them. Prefer a small
number of durable documents over many overlapping summaries. Link to source
paths and owning documents instead of copying large sections.

## Writing Rules

- Describe what the code does now. Label future direction explicitly.
- Separate verified facts from inference.
- Name exact packages, files, types, and functions that orient a maintainer.
- Explain boundaries and invariants, not every helper.
- Include where to start for common changes such as CLI configuration, parser
  syntax, fragment extraction, comment filtering, or error formatting.
- Avoid claims like "always" or "never" unless code, tests, or confirmed policy
  supports them.
- Keep command examples executable from the repository root.
- Use ASCII unless the existing document requires another character set.

## Consistency Checks

Before finishing:

- verify every referenced path exists;
- verify command names and flags against code;
- verify package descriptions against imports and callers;
- verify parser forms and errors against tests and fixtures;
- search for stale names after renames;
- review the diff for duplicated or contradictory guidance;
- run relevant tests when documentation changes depend on behavioral claims.

## Completion Report

Report:

- documentation files changed;
- architecture or workflow facts added or corrected;
- source and tests used as evidence;
- any unresolved question or intentionally omitted area;
- verification performed.

Never commit, push, tag, or rewrite Git history.
