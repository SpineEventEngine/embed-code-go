# Escaped Glob Characters

Glob characters are useful in patterns, but sometimes the source text contains
one literally.

## How It Works

The pattern `Use \* to multiply` treats `*` as source text instead of a wildcard.
It matches the line in [../../code/text/glob-patterns.txt](../../code/text/glob-patterns.txt)
that contains a literal asterisk and embeds only that line.

<embed-code file="$text/glob-patterns.txt" line="Use \* to multiply"></embed-code>
```text
Use * to multiply
```
