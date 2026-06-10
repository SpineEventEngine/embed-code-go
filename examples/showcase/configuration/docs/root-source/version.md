# Repository Root Line Pattern

The same root-source configuration can select a single line from a root file.

## How It Works

The `$repo` prefix points at the repository root configured in
[../../root-source.yml](../../root-source.yml). The anchored pattern selects the
`Version` constant from `main.go`, which keeps the example small while proving
that root files can be used as sources.

<embed-code file="$repo/main.go" line="^const Version = *"></embed-code>
```go
const Version = "1.2.2"
```
