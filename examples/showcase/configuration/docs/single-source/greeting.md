# Single Source Root

This config uses one unnamed `code-path`.

## How It Works

[../../single-source.yml](../../single-source.yml) points `code-path` directly
at `examples/showcase/code/java`. The instruction therefore uses
`org/showcase/Greeting.java` instead of `$java/org/showcase/Greeting.java`.

<embed-code file="org/showcase/Greeting.java" fragment="main()"></embed-code>
```java
public static void main(String[] args) {
    System.out.println(greeting("Ada"));
}
```
