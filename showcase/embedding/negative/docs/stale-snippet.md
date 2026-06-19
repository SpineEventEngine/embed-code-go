# Stale Snippet

This scenario is syntactically valid, but the rendered code is out of date.

## How It Fails

Check mode resolves the `main()` fragment from
[Greeting.java](../../../code/java/org/showcase/Greeting.java)
and compares it with the existing fence. The fence contains different text, so
check mode reports the document as stale without rewriting it. Embed mode would
replace the fence with the current source fragment.

<embed-code file="$java/org/showcase/Greeting.java" fragment="main()"></embed-code>
```java
public static void main(String[] args) {
    System.out.println("Out of date");
}
```
