# Pattern Escaping

Pattern escaping distinguishes glob syntax from source text that happens to use
the same characters.

## How It Works

Patterns use glob control characters, so `*`, `?`, and character classes have
special meaning unless they are escaped with a backslash. The anchors `^` and
`$` are special only at the beginning and end of a pattern part. Use `^^` at the
beginning to match a literal caret and `$$` at the end to match a literal dollar
sign.

The sequence `\n` separates consecutive pattern lines. Use `\\n` when the source
line contains the literal characters `\n`.

## Literal Asterisk

The pattern `Use \* to multiply` treats `*` as source text instead of a
wildcard. It matches a line in
[glob-patterns.txt](../../code/text/glob-patterns.txt).

<embed-code file="$text/glob-patterns.txt" line="Use \* to multiply"></embed-code>
```text
Use * to multiply
```

## Literal Dollar At The End

`$` is an end anchor only at the end of a pattern. Use `$$` there when the
source line itself ends with a dollar sign.

<embed-code file="$text/glob-patterns.txt" line="The value ends with $$"></embed-code>
```text
The value ends with $
```

## Literal Caret At The Start

`^` is a start anchor only at the start of a pattern. Use `^^` there when the
source line itself starts with a caret.

<embed-code file="$text/glob-patterns.txt" line="^^ starts with caret"></embed-code>
```text
^ starts with caret
```

## Literal Backslash-N Text

Use `\\n` when the source line contains the characters backslash and `n`. The
quote characters are written as `\"` so the instruction remains valid XML.

<embed-code
  file="$java/org/showcase/PatternSamples.java"
  line="ESCAPED_NEWLINE = \"\\n\""></embed-code>
```java
private static final String ESCAPED_NEWLINE = "\n";
```
