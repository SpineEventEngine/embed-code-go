import org.gradle.api.publish.maven.MavenPublication
import org.gradle.external.javadoc.StandardJavadocDocletOptions
import org.gradle.plugin.compatibility.compatibility
import org.jetbrains.kotlin.gradle.dsl.JvmTarget
import org.jetbrains.kotlin.gradle.tasks.KotlinJvmCompile

plugins {
    alias(libs.plugins.kotlin.jvm)
    alias(libs.plugins.gradle.plugin.publish)
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

tasks.named<Jar>("jar") {
    from(layout.projectDirectory.file("../LICENSE")) {
        into("META-INF")
    }
}

// Getter docs use concise "Returns..." prose instead of duplicate `@return` tags.
tasks.withType<Javadoc>().configureEach {
    (options as StandardJavadocDocletOptions).addBooleanOption("Xdoclint:-missing", true)
}

dependencies {
    testImplementation(libs.junit.jupiter)
    testImplementation(libs.kotest.assertions.core)
    testRuntimeOnly(libs.junit.launcher)
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
    website.set(
        "https://github.com/SpineEventEngine/embed-code-go/blob/master/gradle-plugin/README.md",
    )
    vcsUrl.set("https://github.com/SpineEventEngine/embed-code-go")
    plugins {
        create("embedCode") {
            id = "io.spine.embed-code"
            implementationClass = "io.spine.embedcode.gradle.EmbedCodePlugin"
            displayName = "Embed Code Gradle Plugin"
            description =
                "Runs Embed Code from Gradle without a separately installed executable."
            tags.set(listOf("documentation", "code-samples"))
            compatibility {
                features {
                    configurationCache = true
                }
            }
        }
    }
}

publishing {
    publications.withType<MavenPublication>().configureEach {
        pom {
            name.set("Embed Code Gradle Plugin")
            description.set(
                "Runs Embed Code from Gradle without a separately installed executable.",
            )
            url.set("https://github.com/SpineEventEngine/embed-code-go")
            licenses {
                license {
                    name.set("The Apache License, Version 2.0")
                    url.set("https://www.apache.org/licenses/LICENSE-2.0.txt")
                    distribution.set("repo")
                }
            }
            developers {
                developer {
                    id.set("SpineEventEngine")
                    name.set("Spine Event Engine")
                    url.set("https://github.com/SpineEventEngine")
                }
            }
            scm {
                url.set("https://github.com/SpineEventEngine/embed-code-go")
                connection.set(
                    "scm:git:https://github.com/SpineEventEngine/embed-code-go.git",
                )
                developerConnection.set(
                    "scm:git:ssh://git@github.com/SpineEventEngine/embed-code-go.git",
                )
            }
        }
    }
}
