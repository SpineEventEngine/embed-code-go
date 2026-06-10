# Overlapping Fragments

Several fragments can open or close on the same marker line. This example uses
an overlapping fragment that shares the class wrapper with another fragment.

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
