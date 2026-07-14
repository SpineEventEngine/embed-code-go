package io.spine.embedcode.gradle

import io.kotest.assertions.throwables.shouldThrow
import io.kotest.matchers.shouldBe
import org.gradle.api.GradleException
import org.junit.jupiter.api.DisplayName
import org.junit.jupiter.api.Test

@DisplayName("`EmbedCodePlatform` should")
internal class EmbedCodePlatformSpec {

    @Test
    fun `select Apple silicon asset`() {
        EmbedCodePlatform.detect("Mac OS X", "aarch64") shouldBe
            EmbedCodePlatform("embed-code-macos-arm64.zip", "embed-code-macos-arm64")
    }

    @Test
    fun `select Intel macOS asset`() {
        EmbedCodePlatform.detect("Mac OS X", "x86_64") shouldBe
            EmbedCodePlatform("embed-code-macos-x64.zip", "embed-code-macos-x64")
    }

    @Test
    fun `select Linux asset`() {
        EmbedCodePlatform.detect("Linux", "amd64") shouldBe
            EmbedCodePlatform("embed-code-linux", "embed-code-linux")
    }

    @Test
    fun `select Windows asset`() {
        EmbedCodePlatform.detect("Windows 11", "amd64") shouldBe
            EmbedCodePlatform("embed-code-windows.exe", "embed-code-windows.exe")
    }

    @Test
    fun `reject platform without release binary`() {
        val error = shouldThrow<GradleException> {
            EmbedCodePlatform.detect("Linux", "aarch64")
        }

        error.message shouldBe
            "Embed Code does not publish a binary for operating system `Linux`" +
            " and architecture `aarch64`."
    }
}
