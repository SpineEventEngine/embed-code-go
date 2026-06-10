# Comment Filtering

The `comments` attribute controls which recognized comments remain in the
rendered snippet. This example keeps documentation comments and removes regular
comments from Java source.

<embed-code file="$java/org/showcase/CommentModes.java" comments="documentation"></embed-code>
```java
package org.showcase;

/**
 * Creates public greetings.
 */
public interface CommentModes {
    String URL = "http://example.org/*not-comment*/";

    String greet(String name); 
}
```
