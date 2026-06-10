# Configuration Examples

These examples show the supported YAML configuration shapes. Run commands from
the repository root.

## Repository Root Source

[root-source.yml](root-source.yml) uses the repository root as a named source
root. Instructions in [docs/root-source](docs/root-source/) embed files from
the project root with the `$repo` prefix.

```bash
go run ./main.go -mode check -config-path examples/showcase/configuration/root-source.yml
```

## Single Showcase Source Root

[single-source.yml](single-source.yml) uses one unnamed `code-path`.
Instructions in [docs/single-source](docs/single-source/) use paths relative to
that root without a `$name` prefix.

```bash
go run ./main.go -mode check -config-path examples/showcase/configuration/single-source.yml
```

## Named Source Roots

[named-sources.yml](named-sources.yml) defines Java, Kotlin, and text source
roots. Instructions in [docs/named-sources](docs/named-sources/) choose a
source root with `$java`, `$kotlin`, or `$text`.

```bash
go run ./main.go -mode check -config-path examples/showcase/configuration/named-sources.yml
```

## Include And Exclude Patterns

[include-exclude.yml](include-exclude.yml) processes Markdown files in
[docs/include-exclude](docs/include-exclude/) but excludes
[excluded.md](docs/include-exclude/excluded.md). That file intentionally
references a missing source file, so the check succeeds only when
`doc-excludes` is applied.

```bash
go run ./main.go -mode check -config-path examples/showcase/configuration/include-exclude.yml
```

## Multiple Embeddings

[multiple-embeddings.yml](multiple-embeddings.yml) uses the `embeddings` list
to process two independent documentation roots in
[docs/multiple](docs/multiple/) in one run.

```bash
go run ./main.go -mode check -config-path examples/showcase/configuration/multiple-embeddings.yml
```
