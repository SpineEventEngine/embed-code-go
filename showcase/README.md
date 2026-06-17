# Embed Code Showcase

This is an executable showcase guide to `embed-code-go` and the end-to-end tests.

## How To Use This Guide

This guide is divided into two categories:

1. [Configuration](configuration/README.md) - describes how to configure the whole embed-code application.
2. [Embedding](embedding/README.md) - describes how to work with the embedding instructions.

## How To Run Tests

Run commands from the repository root.

Run the opt-in end-to-end test with the `showcase` build tag:

```bash
go test -tags showcase ./showcase
```

Verify the positive embedding examples:

```bash
go run ./main.go -mode=check -config-path=showcase/embedding/embed-code.yml
```

Verify the configuration examples:

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
