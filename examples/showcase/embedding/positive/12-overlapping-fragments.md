# Overlapping Fragments

Several fragments can open or close on the same marker line.

## How It Works

[../code/java/org/showcase/OverlappingFragments.java](../code/java/org/showcase/OverlappingFragments.java)
uses marker lines that name both `Class wrapper` and `Greeting method`. The
instruction asks only for `Greeting method`, so embed mode keeps the shared
class wrapper, skips unrelated method details, and renders the selected method
inside the wrapper.

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
