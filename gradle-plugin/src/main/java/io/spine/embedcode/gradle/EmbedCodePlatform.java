package io.spine.embedcode.gradle;

import org.gradle.api.GradleException;

import java.util.Locale;
import java.util.Objects;

/** A released executable selected for an operating system and architecture. */
final class EmbedCodePlatform {

    private final String assetName;
    private final String executableName;

    EmbedCodePlatform(String assetName, String executableName) {
        this.assetName = assetName;
        this.executableName = executableName;
    }

    /** Returns the platform-specific release asset name. */
    String getAssetName() {
        return assetName;
    }

    /** Returns the executable name after extraction. */
    String getExecutableName() {
        return executableName;
    }

    /** Selects the release asset for {@code osName} and {@code architecture}. */
    static EmbedCodePlatform detect(String osName, String architecture) {
        String os = osName.toLowerCase(Locale.ROOT);
        String arch = architecture.toLowerCase(Locale.ROOT);
        boolean isAmd64 = arch.equals("amd64") || arch.equals("x86_64");
        boolean isArm64 = arch.equals("aarch64") || arch.equals("arm64");

        if (os.contains("mac") && isArm64) {
            return new EmbedCodePlatform(
                    "embed-code-macos-arm64.zip",
                    "embed-code-macos-arm64"
            );
        }
        if (os.contains("mac") && isAmd64) {
            return new EmbedCodePlatform(
                    "embed-code-macos-x64.zip",
                    "embed-code-macos-x64"
            );
        }
        if (os.contains("linux") && isAmd64) {
            return new EmbedCodePlatform("embed-code-linux", "embed-code-linux");
        }
        if (os.contains("windows") && isAmd64) {
            return new EmbedCodePlatform(
                    "embed-code-windows.exe",
                    "embed-code-windows.exe"
            );
        }
        throw new GradleException(
                "Embed Code does not publish a binary for operating system `" + osName
                        + "` and architecture `" + architecture + "`."
        );
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof EmbedCodePlatform)) {
            return false;
        }
        EmbedCodePlatform that = (EmbedCodePlatform) other;
        return assetName.equals(that.assetName)
                && executableName.equals(that.executableName);
    }

    @Override
    public int hashCode() {
        return Objects.hash(assetName, executableName);
    }

    @Override
    public String toString() {
        return "EmbedCodePlatform{" +
                "assetName='" + assetName + '\'' +
                ", executableName='" + executableName + '\'' +
                '}';
    }
}
