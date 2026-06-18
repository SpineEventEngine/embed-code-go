# Instruction Tag

An embedding instruction is an XML-like tag placed immediately before the code
fence that embed-code should manage. The tag contains source-selection
attributes such as `file`, `fragment`, `line`, etc.

## Attributes And Selection Rules

Every instruction needs a `file` attribute. The value is resolved from the
configured `code-path`. With named source roots, start the value with `$name/`,
such as `$java/org/showcase/Greeting.java`. With one unnamed source root, use a
path relative to that root.

Use one source-selection shape per instruction:

- `file` alone embeds the whole source file.
- `file` with `fragment` embeds one named fragment.
- `file` with `line` embeds the first matching source line, or consecutive
  source lines when the pattern contains `\n`.
- `file` with `start`, `end`, or both embeds an inclusive source range. Without
  `start`, the range begins at the first source line. Without `end`, it
  continues through the last source line.

The `fragment` attribute cannot be combined with `start`, `end`, or `line`.
The `line` attribute cannot be combined with `start` or `end`.

Add `comments` when embedded snippets should keep only selected source comments.
Supported values are `all`, `none`, `documentation`, `regular`, `inline`, and
`block`; omitting `comments` is the same as `comments="all"`. See
[comment-filtering.md](comment-filtering.md) for language support and examples.

Pattern details live in [source-line-pattern.md](source-line-pattern.md),
[start-end-pattern.md](start-end-pattern.md),
[multi-line-pattern.md](multi-line-pattern.md), and
[pattern-escaping.md](pattern-escaping.md).

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
