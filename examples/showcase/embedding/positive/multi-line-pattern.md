# Multi-Line Patterns

Use `\n` inside a pattern when one source line is not specific enough.

## How It Works

A multi-line pattern is a sequence of ordinary line patterns separated by `\n`.
The match succeeds only when those patterns match neighboring source lines in
the same order. This works for `start`, `end`, and `line` patterns.

Each part keeps the same glob behavior as a one-line pattern. When a part does
not start with `^`, it may begin anywhere in the source line. When it does not
end with `$`, it may stop before the source line ends. Add `^`, `$`, or both to
the individual part that needs a stricter boundary.

## Start And End Pattern

<embed-code
  file="$java/org/showcase/PatternSamples.java"
  start="Scenario \n adds two numbers"
  end="assertEquals(2, total); \n }"></embed-code>
```java
@Scenario
@Name("adds two numbers")
void addsTwoNumbers() {
    int total = 1 + 1;

    assertEquals(2, total);
}
```

## Line Pattern

The `line` attribute also accepts `\n`. In this case the matched consecutive
source lines become the rendered snippet.

<embed-code
  file="$java/org/showcase/PatternSamples.java"
  line="Scenario \n adds two numbers"></embed-code>
```java
@Scenario
@Name("adds two numbers")
```
