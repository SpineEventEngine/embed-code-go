# Named Source Roots

When a configuration has named code paths, the instruction chooses one with the
`$name/relative/path` prefix. This example reads Kotlin code from the `kotlin`
source root while the Java examples read from `java`.

<embed-code file="$kotlin/org/showcase/KotlinGreeting.kt" fragment="main()"></embed-code>
```kotlin
@JvmStatic
fun main(args: Array<String>) {
    println("Hello from Kotlin")
}
```
