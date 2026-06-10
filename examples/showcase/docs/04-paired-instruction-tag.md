# Paired Instruction Tag

Instructions may be self-closing or paired. The self-closing form is supported,
but this showcase uses paired tags because some Markdown renderers display the
XML-style self-closing tag awkwardly.

## How It Works

The active instruction below has an opening `<embed-code>` tag and a matching
closing tag. No content is required between them; the rendered snippet still
comes from the source file and the following code fence. The whole showcase uses
this paired form so Markdown previews display the instructions consistently.

<embed-code file="$java/org/showcase/Greeting.java" fragment="main()"></embed-code>
```java
public static void main(String[] args) {
    System.out.println(greeting("Ada"));
}
```
