# Named Fragment

The `fragment` attribute embeds lines between matching `#docfragment` and
`#enddocfragment` markers. Marker lines are not rendered.

<embed-code file="$java/org/showcase/Greeting.java" fragment="main()"></embed-code>
```java
public static void main(String[] args) {
    System.out.println(greeting("Ada"));
}
```
