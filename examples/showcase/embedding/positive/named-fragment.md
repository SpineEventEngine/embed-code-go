# Named Fragment

Use `fragment` when the source file already marks a reusable documentation
region.

## How It Works

[../../code/java/org/showcase/Greeting.java](../../code/java/org/showcase/Greeting.java)
wraps the `main()` method with matching `#docfragment` and `#enddocfragment`
comments. The instruction resolves the named region, removes the marker lines,
normalizes indentation, and replaces the following fence with the method body.

<embed-code file="$java/org/showcase/Greeting.java" fragment="main()"></embed-code>
```java
public static void main(String[] args) {
    System.out.println(greeting("Ada"));
}
```
