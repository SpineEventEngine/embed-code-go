# Instruction Tag

An embedding instruction is an XML-like tag placed immediately before the code
fence that embed-code should manage. The tag contains source-selection
attributes such as `file`, `fragment`, `line`, etc.

The following Markdown fences keeps the language label used by renderers for
syntax highlighting.

## Paired Tag

The paired form is preferred in Markdown because it is displayed consistently
by most renderers.

<embed-code file="$java/org/showcase/Greeting.java" fragment="main()"></embed-code>
```java
public static void main(String[] args) {
    System.out.println(greeting("Ada"));
}
```

## Self-Closing Tag

The self-closing form is supported and resolves the same source content,
but it is preferred to use paired tags elsewhere because they tend to look better
in Markdown previews.

<embed-code file="$java/org/showcase/Greeting.java" fragment="main()" />
```java
public static void main(String[] args) {
    System.out.println(greeting("Ada"));
}
```
