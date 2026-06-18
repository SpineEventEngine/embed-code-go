# Whole File Source

Use a whole-file embedding when the documentation should mirror a complete
small source file. This is the smallest useful instruction: it only needs
`file`, followed by the code fence.

## How It Works

The `file` attribute is resolved from the configured source roots.

In this example, the `$java` prefix selects the Java source root from
[embed-code.yml](../embed-code.yml) before resolving `org/showcase/Greeting.java`.

## Embedding Instruction

The instruction below embeds the contents of
[Greeting.java](../../code/java/org/showcase/Greeting.java)
into the following code fence, without the embed-code-related instructions.

<embed-code file="$java/org/showcase/Greeting.java"></embed-code>
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

public final class Greeting {
    private Greeting() {}

    public static void main(String[] args) {
        System.out.println(greeting("Ada"));
    }

    public static String greeting(String name) {
        return "Hello, " + name + "!";
    }
}
```
