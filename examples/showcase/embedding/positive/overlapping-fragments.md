# Overlapping Fragments

Use overlapping fragments when different documentation examples need to share
some source lines but hide different details. One marker line can name several
fragments, and each named fragment is resolved independently.

## How It Works

A marker can open or close multiple fragments by listing several quoted names:
`#docfragment "Class wrapper", "Greeting method"`.

## Embedding Instruction

[../../code/java/org/showcase/OverlappingFragments.java](../../code/java/org/showcase/OverlappingFragments.java)
uses marker lines that name both `Class wrapper` and `Greeting method`. The
instruction asks only for `Greeting method`, so the rendered snippet keeps the
shared class wrapper and the greeting method while replacing skipped details
with the configured separator.

<embed-code file="$java/org/showcase/OverlappingFragments.java" fragment="Greeting method"></embed-code>
```java
public final class OverlappingFragments {
    // ...
    public static String greeting(String name) {
        // ...
        return "Hello, " + normalized + "!";
    }
// ...
}
```
