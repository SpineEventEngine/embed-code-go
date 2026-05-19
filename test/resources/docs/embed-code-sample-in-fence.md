# Example with a literal `<embed-code>` sample

````markdown
<embed-code
  file="$root/version.gradle.kts"
  start="val validationVersion"
  end="val validationVersion">
</embed-code>
```kotlin
val validationVersion by extra("2.0.0-SNAPSHOT.419")
```
````
