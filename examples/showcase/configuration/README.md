# Configuration Examples

These examples show the supported YAML configuration shapes.

Each YAML file has a matching docs root under [docs](docs/). Read the YAML file
first, then open the linked docs folder to see how instructions use that source
configuration.

## Repository Root Source

[root-source.yml](root-source.yml) uses the repository root as a named source
root. Instructions in [docs/root-source](docs/root-source/) embed files from
the project root with the `$repo` prefix.

Use this shape only when documentation really needs files from the repository
root. The main embedding showcase avoids root sources so ordinary examples stay
independent from project metadata.

```bash
go run ./main.go -mode check -config-path examples/showcase/configuration/root-source.yml
```

## Single Showcase Source Root

[single-source.yml](single-source.yml) uses one unnamed `code-path`.
Instructions in [docs/single-source](docs/single-source/) use paths relative to
that root without a `$name` prefix.

This is the simplest configuration for one source tree and one documentation
tree.

```bash
go run ./main.go -mode check -config-path examples/showcase/configuration/single-source.yml
```

## Named Source Roots

[named-sources.yml](named-sources.yml) defines Java, Kotlin, and text source
roots. Instructions in [docs/named-sources](docs/named-sources/) choose a
source root with `$java`, `$kotlin`, or `$text`.

Use this shape when one docs tree needs snippets from several source trees.

```bash
go run ./main.go -mode check -config-path examples/showcase/configuration/named-sources.yml
```

## Include And Exclude Patterns

[include-exclude.yml](include-exclude.yml) processes Markdown files in
[docs/include-exclude](docs/include-exclude/) but excludes
[excluded.md](docs/include-exclude/excluded.md). That file intentionally
references a missing source file, so the check succeeds only when
`doc-excludes` is applied.

Use this shape to skip drafts, generated docs, deprecated pages, or any file
that should not be scanned for active instructions.

```bash
go run ./main.go -mode check -config-path examples/showcase/configuration/include-exclude.yml
```

## Multiple Embeddings

[multiple-embeddings.yml](multiple-embeddings.yml) uses the `embeddings` list
to process two independent documentation roots in
[docs/multiple](docs/multiple/) in one run.

Use this shape when one command should process several independent
documentation targets with different source roots or settings.

```bash
go run ./main.go -mode check -config-path examples/showcase/configuration/multiple-embeddings.yml
```
