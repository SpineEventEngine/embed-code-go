# Whole File From A Showcase Source

This is the smallest useful instruction: it names a source file and leaves the
selection attributes empty.

## How It Works

The `$java` prefix selects the Java source root from
[../embed-code.yml](../embed-code.yml). Because `fragment`, `start`, `end`, and
`line` are omitted, embed mode copies every line from
[../../code/java/org/showcase/Greeting.java](../../code/java/org/showcase/Greeting.java)
into the following code fence.

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
