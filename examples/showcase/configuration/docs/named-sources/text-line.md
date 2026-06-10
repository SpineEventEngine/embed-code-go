# Named Text Source

Named roots do not need to be language-specific. This example embeds one line
from the text source root.

## How It Works

[../../named-sources.yml](../../named-sources.yml) defines a source root named
`text`. The `line` pattern matches one plain-text line from
`glob-patterns.txt`, showing that source roots can point at any supported text
fixture, not only programming language files.

<embed-code file="$text/glob-patterns.txt" line="The total is $5"></embed-code>
```text
The total is $5
```
