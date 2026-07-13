import org.jetbrains.kotlin.gradle.dsl.JvmTarget
import org.jetbrains.kotlin.gradle.tasks.KotlinJvmCompile

plugins {
    alias(libs.plugins.kotlin.jvm)
    `java-gradle-plugin`
    `maven-publish`
}

group = "io.spine.tools"
version = "0.1.0"

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
