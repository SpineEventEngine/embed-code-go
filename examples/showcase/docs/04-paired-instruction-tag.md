# Paired Instruction Tag

Instructions may be self-closing or paired. The self-closing form is supported,
but this showcase uses paired tags because some Markdown renderers display the
XML-style self-closing tag awkwardly.

## Short Paired Tag

This compact paired form is useful when all attributes fit naturally in one
line.

<embed-code file="$java/org/showcase/Greeting.java" fragment="main()"></embed-code>
```java
public static void main(String[] args) {
    System.out.println(greeting("Ada"));
}
```

## Paired Tag With A Larger Fragment

The same paired form works for larger snippets too.

<embed-code file="$java/org/showcase/Greeting.java" fragment="Greeter class"></embed-code>
```java
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
