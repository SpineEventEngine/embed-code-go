# Instruction Tag

An embedding instruction is an XML-like tag placed immediately before the code
fence that embed-code should manage. The tag contains source-selection
attributes such as `file`, `fragment`, `line`, etc.

## Managed Code Fence

An instruction must be followed immediately by a Markdown code fence. Embed mode
replaces the fence content. Check mode compares the fence content with current
source and reports stale files without rewriting them.

Use a language label on the fence for syntax highlighting:

````markdown
<embed-code file="$java/org/showcase/Greeting.java" fragment="main()"></embed-code>
```java
```
````

The same rule applies inside HTML documents when include patterns allow `.html` files.

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
