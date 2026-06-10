# Named Java Source

This config has multiple named source roots. The `$java` prefix chooses the Java
root before resolving the relative path.

## How It Works

[../../named-sources.yml](../../named-sources.yml) defines a source root named
`java`. The instruction must include `$java` so the resolver knows which source
tree owns `org/showcase/Greeting.java`.

<embed-code file="$java/org/showcase/Greeting.java" fragment="main()"></embed-code>
```java
public static void main(String[] args) {
    System.out.println(greeting("Ada"));
}
```
