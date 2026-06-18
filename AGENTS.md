# Agent Instructions

These instructions apply to the whole repository. Keep this file focused on
agent operating policy and repository-wide rules.

## Operating Policy

- Read this file before changing the workspace.
- Read [PROJECT.md](PROJECT.md) for architecture, processing flow, package
  ownership, parser and embedding rules, and documentation ownership.
- Always use matching automatically discovered skills when they are available.
- Ask clarifying questions before implementation, review, or documentation work
  when scope, acceptance criteria, or constraints are not explicit.
- Never create commits, push, tag, merge, rebase, cherry-pick, or rewrite Git
  history in this repository.
- Preserve unrelated local changes. Treat them as user work.

## Project Reference

See [PROJECT.md](PROJECT.md) for the project architecture, package ownership,
processing flow, parser and embedding rules, and documentation map.

## Repository Hygiene

- Do not revert unrelated user changes.
- Do not edit generated binaries under `bin/` unless explicitly requested.
- Do not add temporary repo files, local binaries, IDE metadata, coverage
  output, or build artifacts to the intended change set.
- Keep changes narrowly scoped and make unrelated cleanup a separate task.
