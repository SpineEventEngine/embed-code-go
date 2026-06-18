# Kotlin Embedding Entry

This document is processed by the `kotlin-guide` entry in the same
`embeddings` config.

## How It Works

[multiple-embeddings.yml](../../../multiple-embeddings.yml) contains a
separate `kotlin-guide` entry. That entry defines a named `kotlin` source root,
so this instruction uses `$kotlin` even though it is processed by the same
top-level command as the Java entry.

<embed-code file="$kotlin/org/showcase/KotlinGreeting.kt" fragment="main()"></embed-code>
```kotlin
@JvmStatic
fun main(args: Array<String>) {
    println("Hello from Kotlin")
}
```
