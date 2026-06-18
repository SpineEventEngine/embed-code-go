# Embed Code Showcase

This is the complete usage guide for `embed-code-go`. It is also executable:
the showcase test runs the examples through the real CLI so the guide stays in
sync with the application.

## Workflow

1. Configure source roots and documentation roots.
2. Add an `<embed-code>` instruction followed by a managed code fence.
3. Run check mode in CI to detect stale snippets.
4. Run embed mode when documentation should be rewritten from source.

## Guide Map

- [Configuration](configuration/README.md): how the CLI finds source files and documentation files.
- [Embedding](embedding/README.md): how instructions select and render source content.
- [Positive examples](embedding/positive): runnable examples that should pass.
- [Negative examples](embedding/negative/docs): intentionally broken examples
  that document diagnostics.

## Run The Showcase

Run commands from the repository root.

The end-to-end suite checks positive examples, expected failures, and all
configuration shapes:

```bash
go test -tags showcase ./showcase
```

Check the positive embedding examples directly:

```bash
go run ./main.go -mode=check -config-path=showcase/embedding/embed-code.yml
```

Check the configuration examples directly:

```bash
go run ./main.go -mode=check -config-path=showcase/configuration/single-source.yml
go run ./main.go -mode=check -config-path=showcase/configuration/named-sources.yml
go run ./main.go -mode=check -config-path=showcase/configuration/include-exclude.yml
go run ./main.go -mode=check -config-path=showcase/configuration/multiple-embeddings.yml
```

The negative examples are intentionally broken, so these commands should fail:

```bash
go run ./main.go -mode=check -config-path=showcase/embedding/negative/processing-errors.yml
go run ./main.go -mode=check -config-path=showcase/embedding/negative/stale.yml
```
