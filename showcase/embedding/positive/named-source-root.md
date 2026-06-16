# Named Source Roots

Use named source roots when one documentation set needs snippets from several
source trees. Names keep instructions explicit and avoid relying on whichever
configured root happens to contain a matching relative path.

## How It Works

In config file, each entry in `code-path` can have a `name` and a `path`. 
When an instruction starts its `file` value with `$name/`, embed-code selects 
only that named root and then resolves the remaining relative path inside it.

[embed-code.yml](../embed-code.yml) defines `java`, `kotlin`, and `text`
source roots. The `$kotlin` prefix below chooses the Kotlin root before
resolving `org/showcase/KotlinGreeting.kt`.

## Embedding Instruction

<embed-code file="$kotlin/org/showcase/KotlinGreeting.kt" fragment="main()"></embed-code>
```kotlin
@JvmStatic
fun main(args: Array<String>) {
    println("Hello from Kotlin")
}
```
