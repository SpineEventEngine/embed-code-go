# Escaped Newline Text

Pattern escaping distinguishes a real multi-line pattern from source text that
contains the characters backslash and `n`.

## How It Works

The pattern uses `\\n` because the source line contains a string literal with
backslash-n text. The quote characters are written as `\"` so the instruction
remains valid XML. The result is one source line, not a two-line match.

<embed-code
  file="$java/org/showcase/PatternSamples.java"
  line="ESCAPED_NEWLINE = \"\\n\""></embed-code>
```java
private static final String ESCAPED_NEWLINE = "\n";
```
