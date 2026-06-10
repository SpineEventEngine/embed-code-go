# Invalid Attributes

This scenario shows a structurally invalid instruction.

## How It Fails

`fragment` selects a named source region, while `line`, `start`, and `end`
select by pattern. The instruction combines `fragment` and `line`, so the tool
rejects it before reading the source file. In a real guide, choose one selection
style for each instruction.

<embed-code file="$java/org/showcase/Greeting.java" fragment="main()" line="main"></embed-code>
```java
```
