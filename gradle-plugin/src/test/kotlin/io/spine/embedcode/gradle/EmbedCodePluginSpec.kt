package io.spine.embedcode.gradle

import io.kotest.matchers.collections.shouldContain
import io.kotest.matchers.shouldBe
import io.kotest.matchers.string.shouldContain
import org.gradle.testkit.runner.GradleRunner
import org.gradle.testkit.runner.TaskOutcome
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.DisplayName
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.Assumptions.assumeTrue
import org.junit.jupiter.api.condition.EnabledOnOs
import org.junit.jupiter.api.condition.OS
import org.junit.jupiter.api.io.TempDir
import java.nio.file.Files
import java.nio.file.Path
import java.util.zip.ZipEntry
import java.util.zip.ZipOutputStream

@DisplayName("`EmbedCodePlugin` should")
@EnabledOnOs(OS.LINUX, OS.MAC)
internal class EmbedCodePluginSpec {

    @TempDir
    private lateinit var projectDirectory: Path

    private lateinit var releaseDirectory: Path

    @BeforeEach
    fun setUp() {
        Files.createDirectories(projectDirectory.resolve("code"))
        Files.createDirectories(projectDirectory.resolve("docs"))
        Files.writeString(projectDirectory.resolve("settings.gradle.kts"), "rootProject.name = \"test-project\"\n")

        releaseDirectory = projectDirectory.resolve("releases")
        createFakeRelease(releaseDirectory)
        writeBuildFile()
    }

    @Test
    fun `run check mode with Gradle configuration`() {
        val result = runner("checkEmbedCode").build()

        result.task(":installEmbedCode")?.outcome shouldBe TaskOutcome.SUCCESS
        result.task(":checkEmbedCode")?.outcome shouldBe TaskOutcome.SUCCESS
        Files.readString(projectDirectory.resolve("mode.txt")).trim() shouldBe "check"

        val arguments = Files.readAllLines(projectDirectory.resolve("arguments.txt"))
        arguments shouldContain "-mode=check"
        arguments shouldContain "-code-path=${projectDirectory.resolve("code").toRealPath()}"
        arguments shouldContain "-docs-path=${projectDirectory.resolve("docs").toRealPath()}"
        arguments shouldContain "-doc-includes=**/*.md,**/*.html"
        arguments shouldContain "-doc-excludes=drafts/**,generated/**"
        arguments shouldContain "-separator=---"
        arguments shouldContain "-info=true"
        arguments shouldContain "-stacktrace=true"
    }

    @Test
    fun `run check mode with Gradle 6_9_4`() {
        val javaHome = System.getenv("EMBED_CODE_GRADLE_6_JAVA_HOME")
        assumeTrue(
            !javaHome.isNullOrBlank(),
            "Set EMBED_CODE_GRADLE_6_JAVA_HOME to a JDK supported by Gradle 6.9.4.",
        )

        val result = runner("checkEmbedCode", useConfigurationCache = false)
            .withGradleVersion("6.9.4")
            .withEnvironment(System.getenv() + ("JAVA_HOME" to javaHome))
            .build()

        result.task(":installEmbedCode")?.outcome shouldBe TaskOutcome.SUCCESS
        result.task(":checkEmbedCode")?.outcome shouldBe TaskOutcome.SUCCESS
        Files.readString(projectDirectory.resolve("mode.txt")).trim() shouldBe "check"
    }

    @Test
    fun `reuse installation when running embed mode`() {
        runner("checkEmbedCode").build()
        releaseDirectory.toFile().deleteRecursively()

        val result = runner("embedCode").build()

        result.task(":installEmbedCode")?.outcome shouldBe TaskOutcome.UP_TO_DATE
        result.task(":embedCode")?.outcome shouldBe TaskOutcome.SUCCESS
        Files.readString(projectDirectory.resolve("mode.txt")).trim() shouldBe "embed"
    }

    @Test
    fun `run with named source roots and generated configuration`() {
        Files.createDirectories(projectDirectory.resolve("company-site"))
        Files.createDirectories(projectDirectory.resolve("browser"))
        writeNamedSourcesBuildFile()

        val result = runner("checkEmbedCode").build()

        result.task(":checkEmbedCode")?.outcome shouldBe TaskOutcome.SUCCESS
        val arguments = Files.readAllLines(projectDirectory.resolve("arguments.txt"))
        arguments shouldContain "-mode=check"
        arguments.single { it.startsWith("-config-path=") }

        val configuration = Files.readString(projectDirectory.resolve("generated-config.json"))
        configuration shouldContain "\"name\": \"company-site\""
        configuration shouldContain "\"path\": \"${projectDirectory.resolve("company-site").toRealPath()}\""
        configuration shouldContain "\"name\": \"jxbrowser\""
        configuration shouldContain "\"path\": \"${projectDirectory.resolve("browser").toRealPath()}\""
        configuration shouldContain "\"docs-path\": \"${projectDirectory.toRealPath()}\""
    }

    @Test
    fun `reject direct and named source roots together`() {
        Files.createDirectories(projectDirectory.resolve("browser"))
        writeNamedSourcesBuildFile(includeDirectSource = true)

        val result = runner("checkEmbedCode").buildAndFail()

        result.output shouldContain
            "Configure exactly one of `codePath` or `namedSource(...)` for Embed Code."
    }

    @Test
    fun `report missing release asset`() {
        releaseDirectory.toFile().deleteRecursively()

        val result = runner("checkEmbedCode").buildAndFail()

        result.output shouldContain "Could not download Embed Code"
    }

    /** Creates a runner using the plugin-under-test classpath. */
    private fun runner(
        vararg arguments: String,
        useConfigurationCache: Boolean = true,
    ): GradleRunner {
        val gradleArguments = arguments.toMutableList()
        if (useConfigurationCache) {
            gradleArguments.add("--configuration-cache")
        }
        gradleArguments.add("--stacktrace")
        return GradleRunner.create()
            .withProjectDir(projectDirectory.toFile())
            .withArguments(gradleArguments)
            .withPluginClasspath()
    }

    /** Writes a consuming build configured entirely through the plugin extension. */
    private fun writeBuildFile() {
        val baseUrl = releaseDirectory.toUri().toString().trimEnd('/')
        Files.writeString(
            projectDirectory.resolve("build.gradle.kts"),
            """
            plugins {
                id("io.spine.embed-code")
            }

            embedCode {
                version.set("1.2.4")
                downloadBaseUrl.set("$baseUrl")
                codePath.set(layout.projectDirectory.dir("code"))
                docsPath.set(layout.projectDirectory.dir("docs"))
                docIncludes.set(listOf("**/*.md", "**/*.html"))
                docExcludes.set(listOf("drafts/**", "generated/**"))
                separator.set("---")
                info.set(true)
                stacktrace.set(true)
            }
            """.trimIndent(),
        )
    }

    /** Writes a consuming build with two named source roots and no YAML file. */
    private fun writeNamedSourcesBuildFile(includeDirectSource: Boolean = false) {
        val baseUrl = releaseDirectory.toUri().toString().trimEnd('/')
        val directSource = if (includeDirectSource) {
            "codePath.set(layout.projectDirectory.dir(\"code\"))"
        } else {
            ""
        }
        Files.writeString(
            projectDirectory.resolve("build.gradle.kts"),
            """
            plugins {
                id("io.spine.embed-code")
            }

            embedCode {
                version.set("1.2.4")
                downloadBaseUrl.set("$baseUrl")
                $directSource
                namedSource("company-site", layout.projectDirectory.dir("company-site"))
                namedSource("jxbrowser", layout.projectDirectory.dir("browser"))
                docsPath.set(layout.projectDirectory)
            }
            """.trimIndent(),
        )
    }

    /** Creates a host-specific fake release asset that records received arguments. */
    private fun createFakeRelease(root: Path) {
        val platform = EmbedCodePlatform.detect(
            System.getProperty("os.name"),
            System.getProperty("os.arch"),
        )
        val versionDirectory = root.resolve("v1.2.4")
        Files.createDirectories(versionDirectory)
        val executable = projectDirectory.resolve(platform.executableName)
        Files.writeString(
            executable,
            """
            #!/bin/sh
            : > arguments.txt
            for argument in "${'$'}@"; do
              printf '%s\n' "${'$'}argument" >> arguments.txt
              case "${'$'}argument" in
                -mode=check) printf 'check\n' > mode.txt ;;
                -mode=embed) printf 'embed\n' > mode.txt ;;
                -config-path=*) cp "${'$'}{argument#-config-path=}" generated-config.json ;;
              esac
            done
            """.trimIndent() + "\n",
        )

        val asset = versionDirectory.resolve(platform.assetName)
        if (platform.assetName.endsWith(".zip")) {
            ZipOutputStream(Files.newOutputStream(asset)).use { zip ->
                zip.putNextEntry(ZipEntry(platform.executableName))
                Files.newInputStream(executable).use { it.copyTo(zip) }
                zip.closeEntry()
            }
        } else {
            Files.copy(executable, asset)
        }
    }
}
