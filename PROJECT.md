# Project Reference

This document owns durable project knowledge for agents and contributors. For
agent operating policy, read [AGENTS.md](AGENTS.md). For user-facing usage, read
[README.md](README.md) and [EMBEDDING.md](EMBEDDING.md).

## Overview

`embed-code-go` is a Go command-line application. It scans Markdown and HTML
documents for `embed-code` instructions, resolves source fragments, renders
them inside code fences, and checks whether existing snippets are up-to-date.

## Architecture Boundaries

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

## Parser And Embedding Rules

- The parser is state-machine based. When changing parsing behavior, inspect
  `embedding/parsing/constants.go`, `state.go`, `context.go`, and the relevant
  state implementation together.
- Preserve the supported instruction forms:
  - `<embed-code file="..."/>`;
  - `<embed-code file="..."></embed-code>`;
  - multiline instructions represented by existing fixtures.
- Instruction-looking text inside an unrelated Markdown code fence is ordinary
  content, not an active embedding instruction.
- Preserve the opening fence's marker length and indentation when recognizing
  the closing fence and rendering source lines.
- Malformed-instruction errors should point to the instruction start line, not
  EOF or a later fence line.
- Prefer concrete parse reasons: missing tag end, missing closing tag, invalid
  XML, missing code fence, or unclosed code fence.
- Do not scan arbitrary later document content to guess where an invalid
  instruction ends.
- Embed mode may rewrite only documents whose generated content changed.
- Check mode is read-only and should identify every stale document it can safely
  inspect before reporting.

## Documentation Map

- `AGENTS.md`: agent operating policy and repository-wide rules.
- `PROJECT.md`: project architecture, processing flow, package ownership,
  parser and embedding rules, and documentation ownership.
- `README.md`: user-facing overview, installation, configuration, modes, flags,
  and normal operation.
- `EMBEDDING.md`: embedding syntax, source markers, patterns, fences,
  separators, comment modes, and examples.
- `.agents/skills/`: detailed operational rules for agents. Keep
  language-specific, testing, and documentation authoring/review policy there
  rather than duplicating it here.
