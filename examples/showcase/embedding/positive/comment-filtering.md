# Comment Filtering

Use `comments` when examples should keep useful API documentation but omit
implementation notes.

## How It Works

The instruction first resolves the source content, then applies comment
filtering before rendering the code fence. If `comments` is omitted, the default
is `all`, so every recognized comment remains in the snippet.

Supported modes are:

- `all` keeps every comment.
- `none` removes every recognized comment.
- `documentation` keeps documentation comments, such as Javadoc.
- `regular` keeps non-documentation line and block comments.
- `inline` keeps non-documentation line comments, such as `//`.
- `block` keeps non-documentation block comments, such as `/* */`.

Comment support depends on the source file extension. Unknown extensions are
embedded unchanged, and not every supported language distinguishes
documentation, regular, inline, and block comments. See
[Comment filtering](../../../../EMBEDDING.md#comment-filtering) for the full
language matrix.

The examples below embed the same Java file with different `comments` modes so
the rendered output can be compared directly.

## All Comments

`comments="all"` keeps every recognized comment. Omitting `comments` would
produce the same result because `all` is the default.

<embed-code file="$java/org/showcase/CommentModes.java" comments="all"></embed-code>
```java
package org.showcase;

/**
 * Creates public greetings.
 */
public interface CommentModes {
    /*
     * Internal implementation note.
     */
    String URL = "http://example.org/*not-comment*/";

    // Regular inline comment.
    String greet(String name); // trailing inline comment.
}
```

## No Comments

`comments="none"` removes every recognized comment.

<embed-code file="$java/org/showcase/CommentModes.java" comments="none"></embed-code>
```java
package org.showcase;

public interface CommentModes {
    String URL = "http://example.org/*not-comment*/";

    String greet(String name); 
}
```

## Documentation Comments

`comments="documentation"` keeps Javadoc and removes regular block and inline
comments. Comment-like text inside string literals stays unchanged because it is
not a real comment.

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

## Regular Comments

`comments="regular"` keeps non-documentation line and block comments and removes
documentation comments.

<embed-code file="$java/org/showcase/CommentModes.java" comments="regular"></embed-code>
```java
package org.showcase;

public interface CommentModes {
    /*
     * Internal implementation note.
     */
    String URL = "http://example.org/*not-comment*/";

    // Regular inline comment.
    String greet(String name); // trailing inline comment.
}
```

## Inline Comments

`comments="inline"` keeps non-documentation line comments and removes block and
documentation comments.

<embed-code file="$java/org/showcase/CommentModes.java" comments="inline"></embed-code>
```java
package org.showcase;

public interface CommentModes {
    String URL = "http://example.org/*not-comment*/";

    // Regular inline comment.
    String greet(String name); // trailing inline comment.
}
```

## Block Comments

`comments="block"` keeps non-documentation block comments and removes inline and
documentation comments.

<embed-code file="$java/org/showcase/CommentModes.java" comments="block"></embed-code>
```java
package org.showcase;

public interface CommentModes {
    /*
     * Internal implementation note.
     */
    String URL = "http://example.org/*not-comment*/";

    String greet(String name); 
}
```
