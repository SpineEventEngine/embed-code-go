# Start And End Patterns

Use `start` and `end` when the source does not contain named fragment markers.

## How It Works

The `start` pattern finds the first line in
[../code/java/org/showcase/PatternSamples.java](../code/java/org/showcase/PatternSamples.java)
that contains `@Scenario`. The `end` pattern then searches after that start
match and stops at the first line that is exactly four spaces followed by `}`.
Both boundary lines are included in the rendered snippet.

<embed-code
  file="$java/org/showcase/PatternSamples.java"
  start="@Scenario"
  end="^    }$"></embed-code>
```java
@Scenario
@Name("adds two numbers")
void addsTwoNumbers() {
    int total = 1 + 1;

    assertEquals(2, total);
}
```
