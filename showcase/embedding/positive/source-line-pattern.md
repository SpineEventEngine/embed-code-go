# Source Line Pattern

Use `line` when the documentation needs one source line instead of a whole
fragment or range.

## How It Works

The `line` attribute uses the same glob-style pattern syntax as `start` and `end`. 
By default, embed-code behaves as if `*` exists at the beginning and end
of the pattern, so `Hello` can match any source line that contains `Hello`.
Use `^` or `$` when the match must start or end at a line boundary.

Only the first matching source line is rendered into the fence. A `line` pattern
cannot be combined with `fragment`, `start`, or `end`.

## Embedding Instruction

The instruction below searches
[Greeting.java](../../code/java/org/showcase/Greeting.java)
and renders the first line that contains `Hello`.

<embed-code file="$java/org/showcase/Greeting.java" line="Hello"></embed-code>
```java
return "Hello, " + name + "!";
```
