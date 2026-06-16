# Missing Pattern

This scenario shows what happens when a line pattern matches nothing.

## How It Fails

The source file is found, but no line in
[../../../code/java/org/showcase/Greeting.java](../../../code/java/org/showcase/Greeting.java)
matches `doesNotExistPattern`. Check mode reports the unmatched pattern. In a
real guide, loosen the glob pattern, add anchors only where needed, or point the
instruction at the intended source file.

<embed-code file="$java/org/showcase/Greeting.java" line="doesNotExistPattern"></embed-code>
```java
```
