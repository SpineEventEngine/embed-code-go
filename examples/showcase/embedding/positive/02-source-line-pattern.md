# One Line From A Showcase Source

Use `line` when the documentation needs one source line instead of a whole
fragment.

## How It Works

The pattern is matched against
[../code/java/org/showcase/Greeting.java](../code/java/org/showcase/Greeting.java).
The opening quote is escaped because instruction attributes are parsed as XML,
and the trailing `*` lets the pattern match the rest of the return expression.
Only the first matching source line is rendered into the fence.

<embed-code file="$java/org/showcase/Greeting.java" line="Hello"></embed-code>
```java
return "Hello, " + name + "!";
```
