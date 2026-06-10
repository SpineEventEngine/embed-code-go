# Multi-Line Patterns

Use `\n` inside a pattern when the match should span consecutive source lines.
Each pattern line still uses the same glob syntax.

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
