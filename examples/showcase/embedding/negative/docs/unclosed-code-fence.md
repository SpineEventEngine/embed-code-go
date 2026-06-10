# Unclosed Code Fence

This scenario shows what happens when the instruction has an opening fence but
no closing fence.

## How It Fails

The parser finds the opening fence after the instruction and then reaches the
end of the file before seeing a matching closing fence. In a real guide, close
the fence with the same marker style and at least the same marker length.

<embed-code file="$java/org/showcase/Greeting.java" fragment="main()"></embed-code>
```java
public static void main(String[] args) {
}
