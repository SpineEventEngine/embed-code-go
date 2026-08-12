# Named Fragment

Use `fragment` when the source file can mark a stable region that documentation
should reuse. Named fragments are usually easier to maintain than line patterns
when the example has a clear semantic boundary, such as a method, class, or
configuration block.

## How It Works

A named fragment is declared in the source file with `#docfragment "name"`
before the first line to include and `#enddocfragment "name"` after the last
line to include. The marker text can sit inside the comment syntax of the source
language, so Java uses `//`, Kotlin uses `//`, and HTML can use `<!-- -->`.
An opening marker may also declare an [`indent-group`](indent-groups.md) when a
multi-part fragment combines independently indented source regions.

The `fragment` value in the embedding instruction must match the source marker
name exactly. During embed mode or check mode, embed-code resolves the named
region, removes the marker lines, normalizes common indentation, and compares or
updates the following code fence. If the fragment name is not present in the
source file, the run reports the missing fragment.

## Source Markers

[Greeting.java](../../code/java/org/showcase/Greeting.java)
declares the `main()` fragment like this:

```java
// #docfragment "main()"
public static void main(String[] args) {
    System.out.println(greeting("Ada"));
}
// #enddocfragment "main()"
```

## Embedding Instruction

The `file` attribute points to the source file, and `fragment` selects the
named region inside that file. A named fragment cannot be combined with
`start`, `end`, or `line`; use one source-selection method per instruction.

<embed-code file="$java/org/showcase/Greeting.java" fragment="main()"></embed-code>
```java
public static void main(String[] args) {
    System.out.println(greeting("Ada"));
}
```
