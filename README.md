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
    May be represented as:
    * single path
        ```yaml
        code-path: path/to/code/root
        ```
    * multiple named paths:
        ```yaml
        code-path: 
          - name: examples
            path: path/to/code/root1
          - name: production
            path: path/to/code/root2
        ```
      When a named path is specified, fragments must be referenced in the embedding instructions 
      using the corresponding path name:
      ```md
      <embed-code file="$PATH_NAME/path/to/file"></embed-code>
      ```
      **Do not forget the dollar sign (`$`) before the path name.**
    
      It is possible to specify a path without a name or with an empty name.
      In this case, fragments will be stored in the root defined by `fragments-path`.

      It is also possible to specify multiple paths with the same name,
      but this may lead to fragments being overwritten if they have the same relative path and name.

  * `docs-path`: (Mandatory) Path to the documentation root.
  * `code-includes`: (Optional) Glob patterns for source files to include.
    It may be represented as a comma-separated string list or as a YAML sequence.
  * `doc-excludes`: (Optional) Glob patterns for documentation files to exclude.
    It may be represented as a comma-separated string list or as a YAML sequence.
  * `doc-includes`: (Optional) Glob patterns for documentation files to include.
    It may be represented as a comma-separated string list or as a YAML sequence.
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

Run following command to build binaries for macOS, Windows and Ubuntu:
```bash
mkdir -p bin && \
GOOS=darwin GOARCH=amd64 go build -trimpath -o bin/embed-code-macos main.go && \
GOOS=windows GOARCH=amd64 go build -trimpath -o bin/embed-code-windows.exe main.go && \
GOOS=linux GOARCH=amd64 go build -trimpath -o bin/embed-code-linux main.go
```

## Development Notes

This repository is configured with the following GitHub workflows:
- `check` — runs tests across different platforms.
- `build_binaries` — builds binaries on push to the `master` branch.
   > Note: This workflow uses a **Deploy Key** instead of the default GitHub Actions bot
   > to bypass the `master` branch protection against direct pushes.
   >
   > If it is necessary to update the Deploy Key, follow these steps:
   > 1. Generate an SSH key pair for GitHub: `ssh -i ~/.ssh/workflow_deploy_key -T git@github.com`.
   > 2. Add the public key (`workflow_deploy_key.pub`) as a **Deploy Key** in GitHub with write access.
   > 3. Add the private key (`workflow_deploy_key`) as a repository secret named `WORKFLOW_DEPLOY_KEY`.

[embed-code-jekyll]: https://github.com/SpineEventEngine/embed-code
