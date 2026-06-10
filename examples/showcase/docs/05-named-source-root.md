# Named Source Roots

Named roots let one documentation set embed source from several directories.

## How It Works

[../embed-code.yml](../embed-code.yml) defines `java`, `kotlin`, and `text`
source roots. The `$kotlin` prefix chooses the Kotlin root before resolving
`org/showcase/KotlinGreeting.kt`. The same docs root can therefore mix Java,
Kotlin, and text snippets without changing the command line.

<embed-code file="$kotlin/org/showcase/KotlinGreeting.kt" fragment="main()"></embed-code>
```kotlin
@JvmStatic
fun main(args: Array<String>) {
    println("Hello from Kotlin")
}
```
