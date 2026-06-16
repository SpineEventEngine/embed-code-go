# Java Embedding Entry

This document is processed by the `java-guide` entry in an `embeddings` config.

## How It Works

[../../../multiple-embeddings.yml](../../../multiple-embeddings.yml) contains a
`java-guide` entry with its own `code-path` and `docs-path`. This document lives
under that entry's docs root, so the unprefixed source path resolves against the
Java source tree.

<embed-code file="org/showcase/Greeting.java" fragment="main()"></embed-code>
```java
public static void main(String[] args) {
    System.out.println(greeting("Ada"));
}
```
