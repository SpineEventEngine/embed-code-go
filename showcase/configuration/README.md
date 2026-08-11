# Configuration

This guide owns command-line and YAML configuration. Start with the smallest
working shape, then add only the options your documentation needs.

Run commands from the repository root.

## Configuration Sources

Embed Code needs source roots and a documentation root. Provide them in exactly
one of these ways:

- Direct command-line roots: `-code-path` and `-docs-path`.
- A YAML file selected with `-config-path`.

Do not combine direct roots with `-config-path`.

Source roots can contain any text files that your documentation embeds.
The examples use Java, Kotlin, and plain text so the configuration stays independent
of the programming language used by the project.

## Command-Line Arguments

- `-mode`: required execution mode, either `embed` or `check`.
- `-code-path`: source root directory. Use with `-docs-path`.
- `-docs-path`: documentation root directory. Use with `-code-path`.
- `-config-path`: YAML configuration file.
- `-doc-includes`: comma-separated documentation glob patterns to include.
  Defaults to `"**/*.md,**/*.html"`.
- `-doc-excludes`: comma-separated documentation glob patterns to exclude.
- `-separator`: text inserted between joined fragment parts. Defaults to `...`.
- `-info`: enables info-level logging when set to `true`.
- `-stacktrace`: prints stack traces for panics when set to `true`.

Direct roots are useful for small projects:

```bash
go run ./main.go \
  -mode=check \
  -code-path=showcase/code/java \
  -docs-path=showcase/configuration/docs/single-source
```

YAML is easier to maintain when you need named roots, include/exclude patterns,
or multiple documentation targets.

## YAML Fields

- `code-path`: source root. Use a string for one unnamed root or a list of
  `{name, path}` entries for named roots. Do not mix named and unnamed roots.
- `docs-path`: documentation root.
- `doc-includes`: string or list of glob patterns to include.
- `doc-excludes`: string or list of glob patterns to exclude.
- `separator`: text inserted between joined fragment parts.
- `info`: enables info-level logging.
- `stacktrace`: prints stack traces for panics.
- `embeddings`: list of complete configurations for independent documentation
  targets. When `embeddings` is set, define `code-path`, `docs-path`, and
  optional settings inside each entry instead of at the root.

Each `embeddings` entry must have a unique `name`.

## Minimal Config

A YAML configuration needs one source root and one documentation root:

```yaml
code-path: showcase/code/java
docs-path: showcase/configuration/docs/single-source
```

This config is shown by [single-source.yml](single-source.yml).

The application scans files under `docs-path`, finds `<embed-code>`
instructions, and resolves each instruction's `file` path from `code-path`. For
example, see the instruction in [greeting.md](docs/single-source/greeting.md).

Relative paths in `code-path` and `docs-path` are resolved from the command's
current working directory.

Run this example:

```bash
go run ./main.go -mode=check -config-path=showcase/configuration/single-source.yml
```

## Add Document Selection

Add `doc-includes` when only some files under `docs-path` should be scanned.
Add `doc-excludes` when selected files should be skipped:

```yaml
code-path: showcase/code/java
docs-path: showcase/configuration/docs/include-exclude
doc-includes:
  - "**/*.md"
doc-excludes:
  - excluded.md
```

This shape is shown by [include-exclude.yml](include-exclude.yml). It processes
[included.md](docs/include-exclude/included.md) and skips
[excluded.md](docs/include-exclude/excluded.md).

Use include and exclude patterns to skip drafts, generated docs, deprecated
pages, or any file that should not be scanned for active instructions.

```bash
go run ./main.go -mode=check -config-path=showcase/configuration/include-exclude.yml
```

## Add Named Source Roots

Use named source roots when one documentation tree embeds snippets from several
source trees:

```yaml
code-path:
  - name: java
    path: showcase/code/java
  - name: kotlin
    path: showcase/code/kotlin
  - name: text
    path: showcase/code/text
docs-path: showcase/configuration/docs/named-sources
```

This shape is shown by [named-sources.yml](named-sources.yml). Its docs live in
[docs/named-sources](docs/named-sources).

Instructions choose a source root with the `$name` prefix:

```markdown
<embed-code file="$kotlin/org/showcase/KotlinGreeting.kt" fragment="main()"></embed-code>
```

Run the named-source example:

```bash
go run ./main.go -mode=check -config-path=showcase/configuration/named-sources.yml
```

## Add Multiple Documentation Targets

Use `embeddings` when one command should process several independent
documentation targets. Each entry has its own `name`, `code-path`, `docs-path`,
and optional settings:

```yaml
embeddings:
  - name: java-guide
    code-path: showcase/code/java
    docs-path: showcase/configuration/docs/multiple/java
  - name: kotlin-guide
    code-path:
      - name: kotlin
        path: showcase/code/kotlin
    docs-path: showcase/configuration/docs/multiple/kotlin
```

This shape is shown by [multiple-embeddings.yml](multiple-embeddings.yml). It
processes [docs/multiple/java](docs/multiple/java) and
[docs/multiple/kotlin](docs/multiple/kotlin) in one run.

```bash
go run ./main.go -mode=check -config-path=showcase/configuration/multiple-embeddings.yml
```
