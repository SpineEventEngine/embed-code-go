# Whole File Source

Use a whole-file embedding when the documentation should mirror a complete
small source file. This is the smallest useful instruction: it only needs
`file`, followed by the code fence.

## How It Works

The `file` attribute is resolved from the configured source roots. 

In this example, the `$java` prefix selects the Java source root from
[../embed-code.yml](../embed-code.yml) before resolving `org/showcase/Greeting.java`.

## Embedding Instruction

The instruction below embeds the contents of
[../../code/java/org/showcase/Greeting.java](../../code/java/org/showcase/Greeting.java)
into the following code fence, without the embed-code-related instructions.

<embed-code file="$java/org/showcase/Greeting.java"></embed-code>
```java
package org.showcase;

public final class Greeting {
    private Greeting() {}

    public static void main(String[] args) {
        System.out.println(greeting("Ada"));
    }

    public static String greeting(String name) {
        return "Hello, " + name + "!";
    }
}
```
