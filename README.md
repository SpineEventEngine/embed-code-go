# Embed Code

Embed Code provides a way to embed code snippets into Markdown files.
This allows developers to easily include code examples within their documentation.

Previously, we used the `embed-code` utility written in [Ruby for Jekyll][embed-code-jekyll].
Since we standardize our sites on Hugo, we rewrote the utility in Go.
This project is the implementation of `embed-code` utility written in Go.

## Key features
- Extracts code fragments from source files and embeds them into documentation.
- Verifies that embedded code samples are up-to-date with the source.
- Supports configuration via command-line arguments or a YAML file.
- Allows embedding specific named fragments or matching code using line patterns.
- Maps multiple code sources to various documentation folders.

For the details of the usage in the documentation and the code, please refer to the [EMBEDDING.md](EMBEDDING.md).

## Running

Embed Code operates in three modes:

1. **Embedding**: Scans documentation files for `<embed-code>` tags and performs the requested embeddings,
   overwriting the content of the target documentation files.

2. **Up-to-Date Check**: Compares the content under `<embed-code>` tags with the corresponding source code fragments.
   If they differ, the tool reports which files are out-of-date.

3. **Analysis**: Verifies that all embeddings have matching source code fragments.
   Any issues are logged to `build/analytics/problem-files.txt`.
 

The mode is selected using the mandatory `-mode` argument:
- `embed`: Performs the embedding process.
- `check`: Checks if embeddings are up-to-date.
- `analyze`: Runs the analysis process.

The tool can be run as a pre-compiled binary or via the Go compiler (requires Go [installed](#installation)).
Binaries are located in the `./bin` directory.

The code and documentation files must be prepared for embedding.
The instructions are provided in the [Setting up documentation and code files](EMBEDDING.md) document.

### Running the binary

To run the binary, use:
```bash
./bin/<binary_name> [arguments]
```

### Running the Go file

#### Running with Go

If you have Go installed (version `1.22.1` recommended), you can run the tool directly:
```bash
go run ./main.go [arguments]
```

### Arguments

The available arguments are:
  * `-mode`: (Mandatory) The execution mode: `embed`, `check`, or `analyze`.
  * `-code-path`: (Optional) Path to the source code root directory.
  * `-docs-path`: (Optional) Path to the documentation root directory.
  * `-config-path`: (Optional) Path to a YAML configuration file containing `code-path` and `docs-path`.
  * `-code-includes`: (Optional) Comma-separated glob patterns for source files to include (e.g., `"**/*.java,**/*.gradle"`). Defaults to `"**/*.*"`.
  * `-code-excludes`: (Optional) Comma-separated glob patterns for source files to exclude.
  * `-doc-includes`: (Optional) Comma-separated glob patterns for documentation files to include. Defaults to `"**/*.md,**/*.html"`.
  * `-fragments-path`: (Optional) Directory for storing code fragments. Defaults to `./build/fragments`.
  * `-separator`: (Optional) String used to separate joined code fragments. Defaults to `...`.
 
Even though the `code-path`, `docs-path`, and `config-path` arguments are optional,
Embed Code still requires the root directories for code and documentation to be set.
This can be done in one of two ways:

1. Provide the `code-path` and `docs-path` arguments, in this case the roots are read directly from the provided paths.
2. Provide the `config-path` argument, in this case the roots are read from the given configuration file.

If neither of these options is provided, the embedding process will fail.
If both options are set, the embedding will also fail.

### Configuration file

Optional settings can be defined in a YAML configuration file:

```yaml
code-path: path/to/code/root
docs-path: path/to/docs/root
code-includes: "**/*.java,**/*.gradle"
doc-excludes: "**/*-old.*,**/deprecated/*.*"
embed-mappings:
  - code-path: path/to/code/root/kotlin
    docs-path: path/to/other/docs
```

The available fields for the configuration file are:
  * `code-path`: (Mandatory) Path to the source code root.
  * `docs-path`: (Mandatory) Path to the documentation root.
  * `code-includes`: (Optional) Glob patterns for source files to include.
  * `doc-excludes`: (Optional) Glob patterns for documentation files to exclude.
  * `doc-includes`: (Optional) Glob patterns for documentation files to include.
  * `fragments-path`: (Optional) Directory for code fragments.
  * `separator`: (Optional) Separator for fragments.
  * `embed-mappings`: (Optional) A list of custom mappings, each containing `code-path` and `docs-path`.

These settings have the same role as the command-line arguments.

## Installation

* Go to https://go.dev/doc/install. Our Go version is `1.22.1`, which can be checked in the [go.mod](./go.mod) file
* Make sure your Go installed successfully with the command
    ```bash
    go version
    ```

## Compilation

Pre-compiled binaries are available in the `./bin` directory.
However, you can also compile the utility manually if Go is [installed](#installation).

Navigate to the project root and run:
```bash
go build -trimpath main.go
```

There may be issues when running `go build` outside of the directory containing `main.go`,
even if the path is specified correctly.

This command creates an executable named `embed-code` (or `embed-code.exe` on Windows).
For further information, please refer to the [docs](https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies).

Without the `-trimpath` flag, Go includes absolute file paths in stack traces 
based on the system where the binary was built. 

[embed-code-jekyll]: https://github.com/SpineEventEngine/embed-code
