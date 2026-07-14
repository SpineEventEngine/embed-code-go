import org.jetbrains.kotlin.gradle.dsl.JvmTarget
import org.jetbrains.kotlin.gradle.tasks.KotlinJvmCompile

plugins {
    alias(libs.plugins.kotlin.jvm)
    `java-gradle-plugin`
    `maven-publish`
}

val versionFile = layout.projectDirectory.file("../VERSION")
val embedCodeVersion = providers.fileContents(
    providers.provider { versionFile },
).asText.map { it.trim() }.get()
require(embedCodeVersion.isNotEmpty()) {
    "The Embed Code version in ../VERSION must not be empty."
}

group = "io.spine.tools"
version = embedCodeVersion

kotlin {
    jvmToolchain(17)
    compilerOptions {
        jvmTarget.set(JvmTarget.JVM_1_8)
        freeCompilerArgs.add("-Xjsr305=strict")
    }
}

java {
    toolchain {
        languageVersion.set(JavaLanguageVersion.of(17))
    }
}

tasks.named<JavaCompile>("compileJava") {
    options.release.set(8)
}

tasks.named<KotlinJvmCompile>("compileTestKotlin") {
    compilerOptions.jvmTarget.set(JvmTarget.JVM_17)
}

dependencies {
    testImplementation(libs.junit.jupiter)
    testImplementation(libs.kotest.assertions.core)
    testRuntimeOnly(libs.junit.platform.launcher)
}

tasks.test {
    useJUnitPlatform()
    inputs.property(
        "embedCodeGradle6JavaHome",
        providers.environmentVariable("EMBED_CODE_GRADLE_6_JAVA_HOME").orElse(""),
    )
}

tasks.processResources {
    val versionProperties = mapOf("embedCodeVersion" to project.version.toString())
    inputs.properties(versionProperties)
    filesMatching("**/version.properties") {
        expand(versionProperties)
    }
}

gradlePlugin {
    plugins {
        create("embedCode") {
            id = "io.spine.embed-code"
            implementationClass = "io.spine.embedcode.gradle.EmbedCodePlugin"
            displayName = "Embed Code"
            description = "Runs Embed Code from Gradle without a separately installed executable"
            tags.set(listOf("documentation", "code-samples"))
        }
    }
}
