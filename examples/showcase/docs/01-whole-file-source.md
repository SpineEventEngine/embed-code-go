# Whole File From A Showcase Source

Omitting `fragment`, `start`, `end`, and `line` embeds the whole source file.

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
