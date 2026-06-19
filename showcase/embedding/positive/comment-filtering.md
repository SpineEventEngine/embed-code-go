# Comment Filtering

Use `comments` when examples should keep useful API documentation but omit
implementation notes.

## How It Works

The instruction first resolves the requested source content using `file` plus
any `fragment`, `start`, `end`, or `line` selection. It then applies comment
filtering before comparing or rendering the code fence. If `comments` is
omitted, the default is `all`, so every recognized comment remains in the snippet.

Supported modes are:

- `all` keeps every comment.
- `none` removes every recognized comment.
- `documentation` keeps documentation comments, such as Javadoc.
- `regular` keeps non-documentation line and block comments.
- `inline` keeps non-documentation line comments, such as `//`.
- `block` keeps non-documentation block comments, such as `/* */`.

Comment support depends on the source file extension. Unknown extensions are
embedded unchanged. Not every supported language distinguishes documentation,
regular, inline, and block comments, so unsupported categories simply have no
comments to keep.

## Language Support

| Language               | Extensions                                               | Supported `comments` modes                                   |
|------------------------|----------------------------------------------------------|--------------------------------------------------------------|
| Java, Kotlin, Groovy   | `.java`, `.kt`, `.kts`, `.groovy`                        | `all`, `none`, `documentation`, `regular`, `inline`, `block` |
| C#                     | `.cs`                                                    | `all`, `none`, `documentation`, `regular`, `inline`, `block` |
| C, C++                 | `.c`, `.h`, `.cc`, `.cpp`, `.cxx`, `.hh`, `.hpp`, `.hxx` | `all`, `none`, `inline`, `block`                             |
| JavaScript, TypeScript | `.js`, `.jsx`, `.ts`, `.tsx`                             | `all`, `none`, `documentation`, `regular`, `inline`, `block` |
| Go                     | `.go`                                                    | `all`, `none`, `inline`, `block`                             |
| Protobuf               | `.proto`                                                 | `all`, `none`, `inline`, `block`                             |
| Python                 | `.py`, `.pyi`, `.pyw`                                    | `all`, `none`                                                |
| YAML                   | `.yml`, `.yaml`                                          | `all`, `none`                                                |
| XML, HTML              | `.xml`, `.html`, `.htm`                                  | `all`, `none`                                                |
| Visual Basic           | `.vb`, `.bas`, `.vbs`, `.vbscript`                       | `all`, `none`, `documentation`, `regular`                    |

The examples below embed the same Java file with different `comments` modes so
the rendered output can be compared directly.

## All Comments

`comments="all"` keeps every recognized comment. Omitting `comments` would
produce the same result because `all` is the default.

<embed-code file="$java/org/showcase/CommentModes.java" comments="all"></embed-code>
```java
/*
 * Copyright 2026, TeamDev. All rights reserved.
 *
 * Redistribution and use in source and/or binary forms, with or without
 * modification, must retain the above copyright notice and the following
 * disclaimer.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

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
/*
 * Copyright 2026, TeamDev. All rights reserved.
 *
 * Redistribution and use in source and/or binary forms, with or without
 * modification, must retain the above copyright notice and the following
 * disclaimer.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

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
/*
 * Copyright 2026, TeamDev. All rights reserved.
 *
 * Redistribution and use in source and/or binary forms, with or without
 * modification, must retain the above copyright notice and the following
 * disclaimer.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package org.showcase;

public interface CommentModes {
    /*
     * Internal implementation note.
     */
    String URL = "http://example.org/*not-comment*/";

    String greet(String name); 
}
```
