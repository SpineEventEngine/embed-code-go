# Start And End Patterns

The `start` and `end` attributes select an inclusive source range. Patterns are
glob-like and the end search starts after the matched start.

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
