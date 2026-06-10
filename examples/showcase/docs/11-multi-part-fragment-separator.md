# Multi-Part Fragment Separator

Fragments with the same name may appear in several source regions. The rendered
parts are joined in source order with the configured separator.

<embed-code file="$java/org/showcase/MultiPartWorkflow.java" fragment="Workflow"></embed-code>
```java
public final class MultiPartWorkflow {
    // ...
    public static void start() {
        // ...
        System.out.println("Start workflow");
    }
// ...
}
```
