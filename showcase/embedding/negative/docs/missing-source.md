# Missing Source

This scenario shows what happens when the `file` attribute cannot be resolved.

## How It Fails

The `$java` root exists, but `org/showcase/DoesNotExist.java` is not present
under [showcase/code/java/](../../../code/java). Check mode reports the missing source
file and leaves the document unchanged. In a real guide, fix the path or add the
missing source file.

<embed-code file="$java/org/showcase/DoesNotExist.java"></embed-code>
```java
```
