# Embed Code Gradle Plugin

The `io.spine.embed-code` plugin runs Embed Code without requiring developers
or CI jobs to download an executable manually. It selects the released binary
for the current platform, installs it under the project's `build/` directory,
and exposes separate check and embed tasks.

## Apply and Configure

After the plugin is published, apply its released version:

```kotlin
plugins {
    id("io.spine.embed-code") version "0.1.0"
}
```

Until then, test the plugin directly from this checkout by adding its build to
the consuming project's `settings.gradle.kts`:

```kotlin
pluginManagement {
    includeBuild("../embed-code-go/gradle-plugin")
}
```

The consuming `build.gradle.kts` can then apply `id("io.spine.embed-code")`
without a version while using that included build.

Configure Embed Code directly in `build.gradle.kts`; no `embed-code.yml` file
is required:

```kotlin
embedCode {
    version.set("1.2.4")
    codePath.set(layout.projectDirectory.dir("src/main/java"))
    docsPath.set(layout.projectDirectory.dir("docs"))
    docIncludes.set(listOf("**/*.md", "**/*.html"))
    docExcludes.set(listOf("drafts/**", "generated/**"))
    separator.set("...")
    info.set(false)
    stacktrace.set(false)
}
```

`version` and `docsPath` are required. Configure either one unnamed `codePath`
or one or more named sources. The other properties use the same defaults as
the Embed Code command-line application.

| Property | Default | Purpose |
|---|---|---|
| `version` | Required | Selects the GitHub release tag and executable version. |
| `codePath` | Required without named sources | Sets one unnamed source root. |
| `namedSource(name, directory)` | Required without `codePath` | Adds a source root selected with `$name/`. |
| `docsPath` | Required | Sets the documentation root to scan. |
| `docIncludes` | `**/*.md`, `**/*.html` | Selects documentation files. |
| `docExcludes` | Empty | Skips matching documentation files. |
| `separator` | `...` | Separates joined fragment parts. |
| `info` | `false` | Enables informational logging. |
| `stacktrace` | `false` | Prints stack traces after panics. |
| `downloadBaseUrl` | GitHub Releases | Selects a release mirror or test repository. |

### Named Source Roots

Use `namedSource` when documentation embeds code from multiple modules:

```kotlin
embedCode {
    version.set("1.2.4")
    namedSource(
        "company-site",
        layout.projectDirectory.dir("company-site"),
    )
    namedSource(
        "jxbrowser",
        layout.projectDirectory.dir("browser"),
    )
    docsPath.set(layout.projectDirectory)
}
```

Embedding instructions select these roots with `$company-site/` and
`$jxbrowser/`. The plugin writes the corresponding Embed Code configuration
into the Gradle task's temporary directory and passes it to the executable;
the project does not need an `embed-code.yml` file.

`codePath` and `namedSource(...)` are mutually exclusive. Multiple independent
documentation targets are not exposed by this Gradle DSL.

## Run

Check that documentation already contains current source snippets:

```bash
./gradlew checkEmbedCode
```

Update documentation in place:

```bash
./gradlew embedCode
```

`installEmbedCode` is an internal preparation task. Gradle runs it
automatically before either execution task and reuses its output until the
requested version, platform, download URL, or build directory changes.

The plugin supports the platforms for which Embed Code currently publishes
release assets:

- Linux AMD64.
- Windows AMD64.
- macOS AMD64 and ARM64.

## Compatibility

The published plugin implementation targets Java 8 bytecode. The tested range
is Gradle 7.6.3 through 9.5.0. The JVM used to run Gradle must also satisfy the
selected Gradle version's own Java compatibility requirements.

The plugin build uses Kotlin DSL and Kotlin tests, while its published classes
are Java. Keeping Kotlin 2.x off the consumer plugin classpath allows older
Gradle Kotlin DSL compilers to load the plugin.

## Develop

Run compilation, plugin validation, unit tests, and TestKit functional tests:

```bash
cd gradle-plugin
./gradlew check
```

The functional tests create local fake release assets. They do not download or
execute a real GitHub release.

Publish the current plugin version to the local Maven repository when testing
it from another checkout:

```bash
./gradlew publishToMavenLocal
```
