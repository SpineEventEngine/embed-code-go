# Excluded Showcase File

This file is intentionally present in the positive docs root but absent from
the positive processing flow.

## How It Works

[../embed-code.yml](../embed-code.yml) lists this file in `doc-excludes`, so
embed mode and check mode skip it even though it matches the include patterns.
The missing source path proves the exclude is active: processing this file would
fail immediately.

<embed-code file="$java/org/showcase/ExcludedAndMissing.java"></embed-code>
```go
```
