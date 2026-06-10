# Repository Root Source

This configuration example uses the repository root as a named source root. The
instruction embeds the module declaration from `go.mod` through the `$repo`
prefix.

## How It Works

[../../root-source.yml](../../root-source.yml) maps the repository root to the
name `repo`. The instruction uses `$repo/go.mod` and a `line` pattern so this
configuration test proves root-source lookup without embedding the whole
project metadata file.

<embed-code file="$repo/go.mod" line="^module *"></embed-code>
```go
module embed-code/embed-code-go
```
