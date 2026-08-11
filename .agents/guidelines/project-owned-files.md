# Project-owned files

## Rule for this repository

This repository has no `.gitmodules`, so it has no submodule contents or
`config`-distributed files to exclude. Treat every tracked path as
project-owned, including root documentation, `.agents/`, `.github/`,
`scripts/`, `showcase/`, and `test/resources/`.

Do not apply the reusable skip rules below in this repository.

## Why ownership matters

Some skills edit or stamp files across a whole repository. For example,
`proofread` rewrites prose. These skills must touch only the files the
repository **owns**. Editing an upstream file is worse than a no-op: the
next `./config/pull` (or submodule update) overwrites the change or turns it
into a merge conflict, and the same edit is needed again.

This page defines *project-owned* and the two kinds of upstream-owned paths to
skip. Each consuming skill applies the rule to its file types; this page defines
the shared mechanism, while the file filter remains in the consuming skill.

## Reusable rules for repositories with submodules

Apply these rules only when `.gitmodules` declares the relevant submodule.
Determine submodule ownership from `.gitmodules`, never from a directory name.
A repository does **not** own the two kinds of tracked files below. Skip them
in every mode: full sweeps, scoped sweeps, and incremental branch-diff runs alike.

### 1. Submodule contents

Skip every path declared as a submodule in `.gitmodules`, and everything
beneath it — the `config` submodule (usually `config/`), the shared-agents
submodule at `.agents/shared`, and any others a repo adds (for example
`BuildSpeed/`, or example projects mounted under `docs/`).

When a skill enumerates files with `git ls-files`, submodules are already
excluded — a submodule appears to the parent repo as a single gitlink, not
as its files. Apply the skip explicitly anyway, so a diff-based run (which
lists changed *paths*) drops a submodule entry the same way.

### 2. Files distributed by the `config` repository

When a repository consumes `config` (its `.gitmodules` declares a submodule
with `path = config`), `config`'s `migrate` step copies shared files *out
of* the submodule into the project tree, where they become ordinary tracked
files. `git ls-files` and `git diff` therefore surface them even though the
project does not own them, and `./config/pull` overwrites them on every pull.

Skip the config-distributed set below. A consuming skill only encounters
the members that match its own file filter — a prose skill meets the
Markdown and source members; a header-stamping skill meets the
copyright-bearing members — so each skill skips the relevant subset without
needing its own copy of this list.

- **`buildSrc/`** — the entire Gradle build-logic tree (hundreds of Kotlin
  files with KDoc); by far the largest source of false edits.
  **Exception:** `buildSrc/src/main/kotlin/module.gradle.kts` is
  project-owned — `migrate` saves and restores it, so `./config/pull` never
  overwrites it. Only that exact file is project-owned; its
  `*-module.gradle.kts` siblings (for example `jvm-module.gradle.kts`) *are*
  distributed and stay skipped.
- **Root documentation** — `AGENTS.md`, `CLAUDE.md`, `CODE_OF_CONDUCT.md`.
- **Root configuration** — `gradle.properties`, `.codecov.yml`, `lychee.toml`.
- **Agent and IDE files** — `.junie/guidelines.md`,
  `.github/copilot-instructions.md`, and `.idea/`.
- **`.github/workflows/<name>`** when the `config` submodule carries a
  workflow of the same basename in its `.github/workflows/` (or
  `.github-workflows/`) directory. Repo-specific workflows — including
  variants that replace a distributed workflow via a `config:replaces`
  directive and so ship under a different basename — stay project-owned.
  When the submodule is not checked out, the comparison is impossible; treat
  the workflow as project-owned.
- **`CONTRIBUTING.md`** only when it is the unmodified org-wide copy —
  byte-for-byte identical to the submodule's `config/CONTRIBUTING.md`.
  `config` writes it **only into a repo that lacks one** (an
  *initialize if absent* step), so a repo shipping its own contributor
  guide owns it. If it differs, or the `config` submodule is not checked
  out so you cannot confirm, treat it as project-owned.

This set mirrors what `config`'s `migrate` script copies. If that script
changes what it distributes, update this page (and any skill that encodes
the set in a script) to match.

## The `config` and `agents` source repositories

The `config` and `agents` repositories declare no `config` submodule, so rule 2
is inert there: their own `AGENTS.md`, `buildSrc/`, `gradle.properties`, and the
rest are project-owned and stay in scope. This is correct: a fix or stamp must
originate at the source that floats to every consumer.

Do **not** instead skip a path merely because a same-named file exists
under the `config/` submodule. `config` carries files it does *not*
distribute (for example its own `README.md`), and skipping by name would
wrongly exclude a project's own `README.md`.

