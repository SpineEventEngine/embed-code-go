# Project

This document gives agents and contributors the project overview and package
map. For agent operating policy, read [AGENTS.md](AGENTS.md). For implementation
details, use the matching discovered skill.

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
- `files/`: filesystem validation and file helpers.
- `indent/`: shared indentation measurement and removal.
- `logging/`: `slog` handler, clickable file references, panic reporting, and
  terminal formatting support.
- `type/`: YAML-compatible string and named-path list types. The import path
  segment is `type`, but the Go package identifier is `_type` because `type` is
  a Go keyword.
- `test/resources/`: parser, embedding, configuration, and source-code fixtures.

## Documentation

- `README.md`: user-facing usage.
- `EMBEDDING.md`: embedding syntax and examples.
- `.agents/skills/`: task-specific rules for implementation, testing, writing,
  and review.
