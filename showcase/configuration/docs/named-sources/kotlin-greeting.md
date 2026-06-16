# Named Kotlin Source

The same config can embed from a different root by changing the source-root
prefix in the instruction.

## How It Works

[../../named-sources.yml](../../named-sources.yml) also defines a source root
named `kotlin`. Changing the prefix to `$kotlin` resolves the same relative
package path under the Kotlin source tree.

<embed-code file="$kotlin/org/showcase/KotlinGreeting.kt" fragment="main()"></embed-code>
```kotlin
@JvmStatic
fun main(args: Array<String>) {
    println("Hello from Kotlin")
}
```
