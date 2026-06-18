# Agent Instructions

These instructions apply to the whole repository. Keep this file focused on
agent operating policy and repository-wide rules.

## Operating Policy

- Read this file before changing the workspace.
- Read [PROJECT.md](PROJECT.md) for the project overview and project map.
- Use matching discovered skills when they are available. Each skill's
  frontmatter is the routing source of truth.
- Ask clarifying questions before implementation, review, or documentation work
  when scope, acceptance criteria, or constraints are not explicit.
- Never create commits, push, tag, merge, rebase, cherry-pick, or rewrite Git
  history in this repository.
- Preserve unrelated local changes. Treat them as user work.

## Repository Hygiene

- Do not revert unrelated user changes.
- Do not edit generated binaries under `bin/` unless explicitly requested.
- Do not add temporary repo files, local binaries, IDE metadata, coverage
  output, or build artifacts to the intended change set.
- Keep changes narrowly scoped and make unrelated cleanup a separate task.
