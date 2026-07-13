# Project

This document gives agents and contributors the project overview, package map,
documentation ownership, and CI notes. For agent operating policy, read
[AGENTS.md](AGENTS.md).

## Overview

`embed-code-go` is a Go command-line application. It scans Markdown and HTML
documents for `embed-code` instructions, resolves source fragments, renders
them inside code fences, and checks whether existing snippets are up-to-date.

## Project Map

- `main.go`: executable entry point, mode dispatch, aggregate errors, and final
  user-facing output.
- `cli/`: flag parsing, YAML loading, validation, and runtime config conversion.
- `configuration/`: defaults and normalized settings used by processing code.
- `embedding/`: document discovery, per-document processing, embedding writes,
  and up-to-date checks.
- `embedding/parsing/`: state-machine parser for instructions, ordinary
  Markdown fences, embedding fences, and rendered-content comparison.
- `embedding/commentfilter/`: language-aware comment-retention filtering.
- `fragmentation/`: whole-file, named-fragment, and line-pattern extraction;
  source lookup; partition assembly; encoding checks; and caches.
- `gradle-plugin/`: Gradle plugin build, automatic release-binary installation,
  check/embed tasks, Kotlin DSL configuration, and TestKit tests.
- `files/`: filesystem validation and file helpers.
- `indent/`: shared indentation measurement and removal.
- `logging/`: `slog` handler, clickable file references, panic reporting, and
  terminal formatting support.
- `type/`: YAML-compatible string and named-path list types. The import path
  segment is `type`, but the Go package identifier is `_type` because `type` is
  a Go keyword.
- `scripts/release/`: helper scripts used by release workflows for signing and
  notarizing macOS binaries.
- `test/resources/`: parser, embedding, configuration, and source-code fixtures.
- `showcase/`: executable user guide and end-to-end example suite.

## Documentation Ownership

- `README.md`: project entry point, short run/build instructions, and links to
  the complete guide.
- `gradle-plugin/README.md`: Gradle plugin application, configuration, task,
  compatibility, and development guide.
- `showcase/README.md`: complete user guide entry point and runnable workflow.
- `showcase/configuration/README.md`: command-line flags, YAML configuration,
  source roots, include/exclude patterns, and multiple embedding targets.
- `showcase/embedding/README.md`: `<embed-code>` instruction syntax, source
  selection, fragments, patterns, comment filtering, and rendered examples.
- `PROJECT.md`: project map, package ownership, documentation ownership, and CI
  notes for contributors and agents.
- `AGENTS.md`: repository operating policy for agents.

Keep usage details in the showcase. Keep architecture and ownership details in
this file. Keep the root README short.

## CI

This repository is configured with these GitHub workflows:

- `check`: runs linting, the normal Go test suite, and the showcase end-to-end
  tests across supported platforms. It also verifies the Gradle plugin on each
  supported runner platform and against the minimum Gradle version.
- `coverage`: uploads Go coverage to Codecov, where `codecov.yml` requires
  90% patch coverage for pull requests. The workflow also runs on pushes to
  `master` to populate base coverage for later pull-request diffs and the
  README badge; `master` pushes do not have a coverage status gate by design.
- `release-binaries`: reads `VERSION`, builds Linux, macOS, and Windows
  binaries, signs and notarizes the macOS ARM64 and x64 ZIPs, and creates the
  matching GitHub Release on pushes to `master`. It runs on a self-hosted macOS
  ARM64 runner because Apple signing and notarization require macOS tooling.

The release tag is `v<version>` from `VERSION`. When the release already exists,
the workflow emits a warning and finishes successfully without rebuilding or
uploading binaries.
