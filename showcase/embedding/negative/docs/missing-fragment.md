# Missing Fragment

This scenario shows what happens when the source file exists but the named
fragment does not.

## How It Fails

[../../../code/java/org/showcase/Greeting.java](../../../code/java/org/showcase/Greeting.java)
contains fragments such as `main()`, but it does not contain `does not exist`.
Check mode reports the missing fragment. In a real guide, use an existing
fragment name or add matching source markers.

<embed-code file="$java/org/showcase/Greeting.java" fragment="does not exist"></embed-code>
```java
```
