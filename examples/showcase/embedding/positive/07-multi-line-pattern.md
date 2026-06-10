# Multi-Line Patterns

Use `\n` inside a pattern when one source line is not specific enough.

## How It Works

The `start` value is split into two consecutive line patterns: one that matches
the `@Scenario` line and one that matches the display-name line. The `end`
value works the same way for the assertion and closing brace. Each pattern line
still uses the normal glob rules, so anchors are optional unless you need exact
line boundaries.

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
