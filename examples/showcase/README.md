# Embed Code Showcase

This folder is an opt-in, executable guide to `embed-code-go`. It is not part of
the normal `go test ./...` flow. The Go test is guarded by the `showcase` build
tag, so run it only when you want to verify the examples end to end.

Run commands from the repository root.

The showcase-owned source examples live under [code](code/). Repository-root
source examples are kept in the configuration showcase only.

## How To Use This Guide

Read the files in [docs](docs/) from `01` onward when learning the embedding
syntax for the first time. Each file owns one positive case, explains how the
instruction is resolved, and shows the rendered code fence that embed mode owns.

Then inspect [negative/docs](negative/docs/) to see the failures that check mode
reports without rewriting files. Finally, read [configuration](configuration/)
to compare the YAML shapes that point documentation roots at source roots.

## Positive Flow

Refresh the generated snippets:

```bash
go run ./main.go -mode embed -config-path examples/showcase/embed-code.yml
```

Verify that the snippets are up-to-date:

```bash
go run ./main.go -mode check -config-path examples/showcase/embed-code.yml
```

Run the opt-in test:

```bash
go test -tags showcase ./examples/showcase
```

The test copies the showcase docs to a temporary directory, runs check mode,
intentionally makes one copied snippet stale, repairs it with embed mode, and
then verifies the negative and configuration cases. The build tag keeps this
larger documentation test out of the default test flow.

The positive showcase covers:

| Case                     | File                                                                         | What it verifies                                                      |
|--------------------------|------------------------------------------------------------------------------|-----------------------------------------------------------------------|
| Whole file               | [docs/01-whole-file-source.md](docs/01-whole-file-source.md)                 | Omitting selection attributes embeds the whole showcase source file.  |
| Source line pattern      | [docs/02-source-line-pattern.md](docs/02-source-line-pattern.md)             | A showcase source file can be matched with a `line` pattern.          |
| Named fragment           | [docs/03-named-fragment.md](docs/03-named-fragment.md)                       | `fragment` uses `#docfragment` markers and omits marker lines.        |
| Paired instruction tags  | [docs/04-paired-instruction-tag.md](docs/04-paired-instruction-tag.md)       | Paired tags are preferred; self-closing tags are still supported.     |
| Named source roots       | [docs/05-named-source-root.md](docs/05-named-source-root.md)                 | `$java/...` and `$kotlin/...` select different configured roots.      |
| Start and end patterns   | [docs/06-start-end-pattern.md](docs/06-start-end-pattern.md)                 | `start` and `end` select an inclusive source range.                   |
| Multi-line patterns      | [docs/07-multi-line-pattern.md](docs/07-multi-line-pattern.md)               | `\n` matches consecutive source lines.                                |
| Escaped glob characters  | [docs/08-escaped-glob-character.md](docs/08-escaped-glob-character.md)       | `\*` matches a literal asterisk.                                      |
| Escaped newline text     | [docs/09-escaped-newline-text.md](docs/09-escaped-newline-text.md)           | `\\n` matches a literal backslash-n sequence.                         |
| Comment filtering        | [docs/10-comment-filtering.md](docs/10-comment-filtering.md)                 | `comments="documentation"` keeps Java documentation comments.         |
| Multi-part fragments     | [docs/11-multi-part-fragment-separator.md](docs/11-multi-part-fragment-separator.md) | Repeated fragment markers are joined with the configured separator.   |
| Overlapping fragments    | [docs/12-overlapping-fragments.md](docs/12-overlapping-fragments.md)         | Multiple fragment names may share marker lines.                       |
| Markdown fence shielding | [docs/13-markdown-fence-shielding.md](docs/13-markdown-fence-shielding.md)   | Instruction-looking text inside a regular fence is ignored.           |
| HTML documents           | [docs/html-showcase.html](docs/html-showcase.html)                           | HTML files can be scanned when included by configuration.             |
| Excludes                 | [docs/ignored-by-exclude.md](docs/ignored-by-exclude.md)                     | Excluded files are not processed even when they contain instructions. |

## Negative Flow

The negative examples are intentionally broken. They are still documentation:
each file describes the mistake, the expected failure reason, and what a user
should fix in a real document. These commands should fail.

```bash
go run ./main.go -mode check -config-path examples/showcase/negative/processing-errors.yml
go run ./main.go -mode check -config-path examples/showcase/negative/stale.yml
```

The processing-error config verifies:

| Case                | File                                                                         | Expected behavior                                    |
|---------------------|------------------------------------------------------------------------------|------------------------------------------------------|
| Missing source      | [negative/docs/missing-source.md](negative/docs/missing-source.md)           | Reports that the code file cannot be found.          |
| Missing fragment    | [negative/docs/missing-fragment.md](negative/docs/missing-fragment.md)       | Reports that the requested fragment cannot be found. |
| Missing pattern     | [negative/docs/missing-pattern.md](negative/docs/missing-pattern.md)         | Reports that no source line matches the pattern.     |
| Invalid attributes  | [negative/docs/invalid-attributes.md](negative/docs/invalid-attributes.md)   | Rejects mutually exclusive selection attributes.     |
| Missing code fence  | [negative/docs/missing-code-fence.md](negative/docs/missing-code-fence.md)   | Requires a fence immediately after an instruction.   |
| Unclosed code fence | [negative/docs/unclosed-code-fence.md](negative/docs/unclosed-code-fence.md) | Reports an instruction fence that reaches EOF.       |

The stale config verifies:

| Case          | File                                                               | Expected behavior                                                      |
|---------------|--------------------------------------------------------------------|------------------------------------------------------------------------|
| Stale snippet | [negative/docs/stale-snippet.md](negative/docs/stale-snippet.md)   | Check mode reports the file as needing an update without rewriting it. |

## Configuration Flow

Configuration examples live under [configuration](configuration/). They cover a
repository-root source, a single source root, named source roots,
include/exclude patterns, and the `embeddings` list for multiple documentation
roots.

Each configuration has its own docs root so the examples can be run separately.
Open the YAML file first, then follow the linked docs root to see how the
instructions use that configuration.

| Case                 | Config                                                                       | Docs root                                                                  |
|----------------------|------------------------------------------------------------------------------|----------------------------------------------------------------------------|
| Repository root      | [configuration/root-source.yml](configuration/root-source.yml)               | [configuration/docs/root-source](configuration/docs/root-source/)           |
| Single source root   | [configuration/single-source.yml](configuration/single-source.yml)           | [configuration/docs/single-source](configuration/docs/single-source/)       |
| Named source roots   | [configuration/named-sources.yml](configuration/named-sources.yml)           | [configuration/docs/named-sources](configuration/docs/named-sources/)       |
| Include and exclude  | [configuration/include-exclude.yml](configuration/include-exclude.yml)       | [configuration/docs/include-exclude](configuration/docs/include-exclude/)   |
| Multiple embeddings  | [configuration/multiple-embeddings.yml](configuration/multiple-embeddings.yml) | [configuration/docs/multiple](configuration/docs/multiple/)                 |

```bash
go run ./main.go -mode check -config-path examples/showcase/configuration/root-source.yml
go run ./main.go -mode check -config-path examples/showcase/configuration/single-source.yml
go run ./main.go -mode check -config-path examples/showcase/configuration/named-sources.yml
go run ./main.go -mode check -config-path examples/showcase/configuration/include-exclude.yml
go run ./main.go -mode check -config-path examples/showcase/configuration/multiple-embeddings.yml
```
