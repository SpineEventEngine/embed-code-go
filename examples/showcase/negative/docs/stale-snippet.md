# Stale Snippet

This scenario is syntactically valid, but check mode reports it as stale because
the rendered code fence does not match the current source fragment.

<embed-code file="$java/org/showcase/Greeting.java" fragment="main()"></embed-code>
```java
public static void main(String[] args) {
    System.out.println("Out of date");
}
```
