# Embed Code

Embed Code is a Go command-line tool that keeps documentation snippets in sync
with source files. It scans Markdown and HTML documents for `<embed-code>`
instructions, resolves the requested source content, and manages the following
code fence.

This project replaces the earlier [`embed-code` utility for Ruby/Jekyll][embed-code-jekyll].

## Start Here

The complete usage guide lives in the [showcase](showcase/README.md). It covers
configuration, embedding instructions, check mode, embed mode, expected
failures, and runnable examples.

## What It Does

- Embeds whole files, named fragments, source ranges, or matching source lines.
- Supports multiple named source roots for one documentation tree.
- Filters comments when examples should omit implementation notes.
- Processes Markdown and HTML documents.
- Runs in `check` mode for CI and `embed` mode to update documentation.

## Run

Download the asset for your platform from [GitHub Releases][releases].

On Linux, for example:

```bash
./embed-code-linux -mode=check -config-path=showcase/embedding/embed-code.yml
```

> It may be necessary to give the executable permission with `chmod +x` on Unix-like systems:
> ```bash
> chmod +x embed-code-linux
> ```

On macOS, download `embed-code-macos.zip`, unzip it, and run the binary:

```bash
unzip embed-code-macos.zip
./embed-code-macos -mode=check -config-path=showcase/embedding/embed-code.yml
```

Or run it with Go:

```bash
go run ./main.go -mode=check -config-path=showcase/embedding/embed-code.yml
```

Use `-mode=embed` when documentation should be rewritten with current source
content. See the [configuration guide](showcase/configuration/README.md) for
all command-line flags and YAML options.

## Build

Use Go `1.26.4`.

```bash
go build -trimpath -o embed-code main.go
```

This creates `embed-code` on Unix-like systems or `embed-code.exe` on Windows.
The `-trimpath` flag prevents local absolute paths from appearing in stack traces.

## Development

Run the normal test suite:

```bash
go test ./...
```

Run the executable showcase:

```bash
go test -tags showcase ./showcase
```

## Documentation

The main user guide is the [showcase](showcase/README.md).
Go package docs are useful for maintainers who need to browse package comments
and exported APIs.

Generate static API docs:

```bash
./scripts/godoc
```

The script writes `build/godoc/` and prints the generated main index file link.

Launch the GoDoc server:

```bash
./scripts/godoc-serve
```

Open the package link printed by the script.

[embed-code-jekyll]: https://github.com/SpineEventEngine/embed-code
[releases]: https://github.com/SpineEventEngine/embed-code-go/releases
