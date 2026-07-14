# Embed Code

[![Coverage](https://codecov.io/gh/SpineEventEngine/embed-code-go/branch/master/graph/badge.svg)](https://codecov.io/gh/SpineEventEngine/embed-code-go)

Embed Code helps keep code examples in Markdown and HTML documentation up to
date without manual copying. It takes snippets directly from source files and
reports when embedded content no longer matches the source.

## Typical Usage

For example, consider a simple project with a Java source file and Markdown documentation:

```text
.
|-- src/
|   `-- com/example/Greeting.java
`-- docs/
    `-- greeting.md
```

The project contains `src/com/example/Greeting.java`:

```java
package com.example;

public final class Greeting {
    public static String message() {
        return "Hello from Embed Code";
    }
}
```

In `docs/greeting.md`, add an instruction `<embed-code>` followed by an empty code fence:

````markdown
# Greeting

How to use our Greeting system:
<embed-code file="com/example/Greeting.java"></embed-code>
```java
```

Additional documentation.
````

The configured embed-code application run embeds code into the documentation file:

````markdown
# Greeting

How to use our Greeting system:
<embed-code file="com/example/Greeting.java"></embed-code>
```java
package com.example;

public final class Greeting {
    public static String message() {
        return "Hello from Embed Code";
    }
}
```

Additional documentation.
````

## Language Support

Embed Code works with any programming language, provided its source files use
valid UTF-8 text. Basic embedding treats source files as text and does not
require a language compiler or parser.

## Next Steps

- Run the [quick start](showcase/quick-start/README.md) for a complete example
  with source, configuration, and documentation files.
- Use the [showcase](showcase/README.md) as the entry point for all user guides
  and runnable examples.
- Read the [configuration guide](showcase/configuration/README.md) for direct
  command-line paths, YAML options, named source roots, document selection, and
  multiple documentation targets.
- Read the [embedding guide](showcase/embedding/README.md) for fragments, line
  and range patterns, comment filtering, and instruction attributes.

## Download

Download the asset for your platform from [GitHub Releases][releases]. 
You do not need to install Go to use a release binary.

| Platform            | Release asset                | Executable               |
|---------------------|------------------------------|--------------------------|
| Linux x64           | `embed-code-linux.zip`       | `embed-code-linux`       |
| macOS Apple silicon | `embed-code-macos-arm64.zip` | `embed-code-macos-arm64` |
| macOS Intel         | `embed-code-macos-x64.zip`   | `embed-code-macos-x64`   |
| Windows x64         | `embed-code-windows.exe`     | `embed-code-windows.exe` |

## Build From Source

Using Embed Code does not require Go. 
To build the application from this repository, install Go `1.26.4` and run:

```bash
go build -trimpath -o embed-code main.go
```

This creates an executable named `embed-code`. 
On Windows, use`-o embed-code.exe` to give it the standard `.exe` suffix. 
The `-trimpath` flag prevents local absolute paths from appearing in stack traces.

You can also run the application directly from the source checkout:

```bash
go run ./main.go -mode=check -config-path=showcase/embedding/embed-code.yml
```

## Testing

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

This project replaces the earlier [`embed-code` utility for Ruby/Jekyll][embed-code-jekyll].

[embed-code-jekyll]: https://github.com/SpineEventEngine/embed-code
[releases]: https://github.com/SpineEventEngine/embed-code-go/releases
