# Start And End Patterns

Use `start` and `end` patterns to select the code snippet.

## How It Works

`start` and `end` select an inclusive source range. Embed-code first searches
for the `start` pattern, then searches for the `end` pattern after that match.
Both matched boundary lines are included in the rendered snippet.

By default, embed-code behaves as if `*` exists at the beginning and end
of the pattern, so `Hello` can match any source line that contains `Hello`.
Use `^` or `$` when the match must start or end at a line boundary.

If `start` is omitted, the range starts at the beginning of the file. 
If `end` is omitted, it continues to the end of the file.

## Embedding Instruction

The instruction below finds the first `@Scenario` in
[../../code/java/org/showcase/PatternSamples.java](../../code/java/org/showcase/PatternSamples.java).
It then stops at the next line that is exactly four spaces followed by `}`.

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
