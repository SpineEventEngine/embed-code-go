# Named Java Source

This config has multiple named source roots. The `$java` prefix chooses the Java
root before resolving the relative path.

<embed-code file="$java/org/showcase/Greeting.java" fragment="main()"></embed-code>
```java
public static void main(String[] args) {
    System.out.println(greeting("Ada"));
}
```
