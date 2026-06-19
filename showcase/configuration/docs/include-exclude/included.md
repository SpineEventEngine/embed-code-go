# Included Document

The `doc-includes` pattern selects this Markdown file, so the instruction is
processed normally.

## How It Works

[include-exclude.yml](../../include-exclude.yml) includes Markdown files
under this docs root. Because this file is not listed in `doc-excludes`, check
mode resolves the instruction and compares the rendered fence with the Java
source fragment.

<embed-code file="org/showcase/Greeting.java" fragment="main()"></embed-code>
```java
public static void main(String[] args) {
    System.out.println(greeting("Ada"));
}
```
