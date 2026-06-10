# Configuration Examples

This folder is a runnable guide to YAML configuration. Start with the smallest
working config, then add only the options your documentation needs.

## Minimal Config

A configuration needs one source root and one documentation root:

```yaml
code-path: examples/showcase/code/java
docs-path: examples/showcase/configuration/docs/single-source
```

This config is shown by [single-source.yml](single-source.yml).

The application scans files under `docs-path`, finds `<embed-code>` instructions,
and resolves each instruction's `file` path from `code-path`. For example, see
instruction in [docs/single-source/greeting.md](docs/single-source/greeting.md).

Relative paths in `code-path` and `docs-path` are resolved 
from the command's current working directory.

Run this example (from the project root):

```bash
go run ./main.go -mode=check -config-path=examples/showcase/configuration/single-source.yml
```

## Add Document Selection

Add `doc-includes` when only some files under `docs-path` should be scanned.
Add `doc-excludes` when selected files should be skipped:

```yaml
code-path: examples/showcase/code/java
docs-path: examples/showcase/configuration/docs/include-exclude
doc-includes:
  - "**/*.md"
doc-excludes:
  - excluded.md
```

This shape is shown by [include-exclude.yml](include-exclude.yml). It processes
[docs/include-exclude/included.md](docs/include-exclude/included.md) and skips
[docs/include-exclude/excluded.md](docs/include-exclude/excluded.md).

Use include and exclude patterns to skip drafts, generated docs, deprecated
pages, or any file that should not be scanned for active instructions.

```bash
go run ./main.go -mode=check -config-path=examples/showcase/configuration/include-exclude.yml
```

## Add Named Source Roots

Use named source roots when one documentation tree embeds snippets from several
source trees:

```yaml
code-path:
  - name: java
    path: examples/showcase/code/java
  - name: kotlin
    path: examples/showcase/code/kotlin
  - name: text
    path: examples/showcase/code/text
docs-path: examples/showcase/configuration/docs/named-sources
```

This shape is shown by [named-sources.yml](named-sources.yml). Its docs live in
[docs/named-sources](docs/named-sources/).

Instructions choose a source root with the `$name` prefix:

```markdown
<embed-code file="$kotlin/org/showcase/KotlinGreeting.kt" fragment="main()"></embed-code>
```

Run the named-source example:

```bash
go run ./main.go -mode=check -config-path=examples/showcase/configuration/named-sources.yml
```

## Add Multiple Documentation Targets

Use `embeddings` when one command should process several independent
documentation targets. Each entry has its own `name`, `code-path`, `docs-path`,
and optional settings:

```yaml
embeddings:
  - name: java-guide
    code-path: examples/showcase/code/java
    docs-path: examples/showcase/configuration/docs/multiple/java
  - name: kotlin-guide
    code-path:
      - name: kotlin
        path: examples/showcase/code/kotlin
    docs-path: examples/showcase/configuration/docs/multiple/kotlin
```

This shape is shown by [multiple-embeddings.yml](multiple-embeddings.yml). It
processes [docs/multiple/java](docs/multiple/java/) and
[docs/multiple/kotlin](docs/multiple/kotlin/) in one run.

```bash
go run ./main.go -mode=check -config-path=examples/showcase/configuration/multiple-embeddings.yml
```

## All Configuration Checks

Run commands from the project root.

```bash
go run ./main.go -mode=check -config-path=examples/showcase/configuration/single-source.yml
go run ./main.go -mode=check -config-path=examples/showcase/configuration/named-sources.yml
go run ./main.go -mode=check -config-path=examples/showcase/configuration/include-exclude.yml
go run ./main.go -mode=check -config-path=examples/showcase/configuration/multiple-embeddings.yml
```
