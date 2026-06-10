# Excluded Document

The config excludes this file. The intentionally missing source file proves
that excluded documents are not processed.

## How It Works

[../../include-exclude.yml](../../include-exclude.yml) lists this file in
`doc-excludes`. The instruction points at a missing source file, but the command
still succeeds because excluded files are skipped before instruction parsing.

<embed-code file="org/showcase/DoesNotExist.java"></embed-code>
```java
```
