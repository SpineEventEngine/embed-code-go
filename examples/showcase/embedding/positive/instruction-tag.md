# Instruction Tag

Instructions may be self-closing or paired. The self-closing form is supported,
but it is preferred to use paired tags because some Markdown renderers display the
XML-style self-closing tag awkwardly.

### Paired tag version

<embed-code file="$java/org/showcase/Greeting.java" fragment="main()"></embed-code>
```java
public static void main(String[] args) {
    System.out.println(greeting("Ada"));
}
```

### Self-closing tag version

<embed-code file="$java/org/showcase/Greeting.java" fragment="main()" />
```java
public static void main(String[] args) {
    System.out.println(greeting("Ada"));
}
```
