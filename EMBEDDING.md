# Setting Up Code Embedding

The `embed-code` utility uses a custom `<embed-code>` tag to insert code snippets from source files into Markdown documentation.

## Embedding options

There are two ways to specify which code fragment to embed:

### Option 1: Named fragments

Use a named fragment defined within the source file.
```markdown
<embed-code file="path/to/file" fragment="Fragment Name"></embed-code>
```

### Option 2: Line patterns

Use glob-style patterns to match the start and end lines of the fragment.
```markdown
<embed-code file="path/to/file" start="first-line-pattern" end="last-line-pattern"></embed-code>
```

## Embedding instruction format

An `<embed-code>` instruction must always be followed by a Markdown code fence 
(opening and closing triple backticks). 

````markdown
<embed-code file="java/lang/String.java" fragment="Constructor"></embed-code>
```java
// The utility will automatically overwrite this content.
```
````

The content inside the code fence is irrelevant as it is automatically updated by the tool.
However, you should specify the language for syntax highlighting (e.g., ` ```java `).

This is true even when embedding into HTML.

## Named fragments

### Marking up source code

To define a named fragment in your source code, wrap the desired lines with
`#docfragment` and `#enddocfragment` comments:

```java
public final class String
    implements java.io.Serializable, Comparable<String>, CharSequence {
    
    // #docfragment "Constructor"
    public String() {
        this.value = new char[0];
    }
    // #enddocfragment "Constructor"
}
```

The `#docfragment` and `#enddocfragment` tags are excluded from the embedded snippet.

### Usage in documentation

To embed a named fragment, add the following to your Markdown file:

````markdown
<embed-code file="java/lang/String.java" fragment="Constructor"></embed-code>
```java
```
````

- **`file`**: The path to the source file relative to the `code-path` defined in your configuration.
- **`fragment`**: The name of the fragment to embed. If omitted, the entire file will be embedded.

Fragment names can be any string, but avoid using double quotes (`"`) or characters reserved by XML.

## Pattern-based fragments

Alternatively, you can specify a fragment using `start` and `end` patterns:

````markdown
<embed-code file="java/lang/String.java" start="*class Hello*" end="}*"></embed-code>
```java
```
````

Patterns match the first and last lines of the desired fragment.
If a pattern is omitted, the fragment will start at the beginning or end at the end of the file, respectively.

### Pattern syntax

The tool supports an extended glob syntax for matching lines:

- `?` — Matches any single character.
- `*` — Matches zero or more characters.
- `[set]` — Matches any single character from the specified set (similar to regex character classes).
- `^` — When used at the start of a pattern, matches the beginning of the line.
- `$` — When used at the end of a pattern, matches the end of the line.

**Note on anchors:**
By default, patterns imply a wildcard (`*`) at both the start and end.
Use `^` and `$` to disable this behavior and match the exact line start or end.

If you need to match a literal `^` at the start of a line, use `^^`.
Similarly, use `$$` to match a literal `$` at the end of a line.

## Advanced use cases

### Joining several parts of code into one fragment

A named fragment may consist of one or several pieces declared in a single file. 
When rendered, the pieces that belong to a single fragment are joined together. 
It is possible to specify a separator between the joined pieces, 
see [Configuration](./README.md#arguments) for the corresponding parameter.

Here is an example of how a multi-piece fragment is rendered.

**Code:**

```java
public final class String
    implements java.io.Serializable, Comparable<String>, CharSequence {

    // #docfragment "Standard Object methods"
    public int hashCode() {
        // Method logic.
        return hash;
    }
    // #enddocfragment "Standard Object methods"
    
    /* Here goes irrelevant code */

    // #docfragment "Standard Object methods"
    public boolean equals(Object anObject) {
        // Method logic.
        return false;
    }
    // #enddocfragment "Standard Object methods"

    /* Here goes more irrelevant code */

    // #docfragment "Standard Object methods"
    public String toString() {
        return this;
    }
    // #enddocfragment "Standard Object methods"
}
```

**Result:**

```java
public int hashCode() {
    // Method logic.
    return hash;
}
...
public boolean equals(Object anObject) {
    // Method logic.
    return false;
}
...
public String toString() {
    return this;
}
```

### Declaring multiple fragments in one line

Multiple fragments can start or end on a single line and may overlap:

```java
public final class String
    implements java.io.Serializable, Comparable<String>, CharSequence {

    // #docfragment "Standard Object methods", "All methods"
    public int hashCode() {
        // Method logic.
        return hash;
    }

    public boolean equals(Object anObject) {
        // Method logic.
        return false;
    }

    public String toString() {
        return this;
    }
    // #enddocfragment "Standard Object methods"

    public boolean startsWith(String prefix, int toffset) {
        // Method logic.
        return true;
    }
    // #enddocfragment "All methods"
}
``` 

### Usage with other languages

The fragments can also be used in other languages:

```html
<html lang="en">
<body>
<!-- #docfragment "html-only" -->
<span class="counter" id="counter"></span>
<!-- #enddocfragment "html-only" -->
</body>
</html>
```
