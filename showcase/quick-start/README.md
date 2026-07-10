# Quick Start

This folder is the smallest runnable Embed Code application setup. It contains:

- [config.yml](config.yml): connects the source and target folders.
- [source](source): source files that can be embedded into documentation.
- [target](target): Markdown or HTML files with `<embed-code>` instructions.

## Run

Download the latest release for your platform from
[GitHub Releases](https://github.com/SpineEventEngine/embed-code-go/releases).
Extract binary file from archive, if necessary, and place it to the `quick-start` folder.

Open this folder before running the example:

```bash
cd showcase/quick-start
```

For this example the `embed-code-linux` binary is used,
but it is the same usage for all other binaries.

Run check mode:

```bash
./embed-code-linux -mode=check -config-path=config.yml
```

When working from this repository's source checkout, the same config can be run
with `go run`:

```bash
go run ../../main.go -mode=check -config-path=config.yml
```

Check mode is read-only. It fails when a rendered snippet in `target` no longer
matches the source file.

Run embed mode to fill or refresh snippets:

```bash
./embed-code-linux -mode=embed -config-path=config.yml
```

> You can try modifying the source file and re-running the embedding to see the changes.

## How It Works

[config.yml](config.yml) points `code-path` at [source](source) and `docs-path`
at [target](target):

```yaml
code-path: source
docs-path: target
```

[target/greeting.md](target/greeting.md) contains an embedding instruction
followed by a managed code fence:

````markdown
<embed-code file="com/example/Greeting.java"></embed-code>
```java
```
````

The `file` value is a relative path resolved from `code-path`, so this instruction reads
[source/com/example/Greeting.java](source/com/example/Greeting.java). 
Embed mode writes the current source content into the managed fence. 
Check mode verifies that the fence is already up-to-date.

For more source/target configuration see [Configuration guide](../configuration/README.md).
For more embed-code instructions see [Embedding guide](../embedding/README.md).
