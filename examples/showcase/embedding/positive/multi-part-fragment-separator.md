# Multi-Part Fragment Separator

Use a multi-part fragment when one documentation example should show several
non-adjacent pieces of a source file as one snippet. This is useful for keeping
the public shape of an example while hiding setup, implementation details, or
unrelated branches between the selected parts.

## How It Works

A fragment becomes multi-part when the same `#docfragment "name"` marker is
opened and closed more than once in the same source file. Embed-code collects
the selected parts in source order, normalizes common indentation across all of
them, and inserts the configured `separator` between neighboring parts.

The default separator is `...`. This showcase uses `// ...` in
[../embed-code.yml](../embed-code.yml) so the separator is valid inside Java
snippets. Separator indentation follows the surrounding rendered code, which
keeps skipped sections readable inside classes and methods.

## Embedding Instruction

[../../code/java/org/showcase/MultiPartWorkflow.java](../../code/java/org/showcase/MultiPartWorkflow.java)
opens and closes the `Workflow` fragment several times. The instruction below
renders those selected parts as one snippet.

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
