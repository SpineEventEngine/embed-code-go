# Comment Filtering

Use `comments` when examples should keep useful API documentation but omit
implementation notes.

## How It Works

The instruction embeds the whole Java file and applies
`comments="documentation"`. Javadoc is retained, regular block comments and
inline comments are removed, and comment-like text inside string literals stays
unchanged because it is not a real comment.

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
