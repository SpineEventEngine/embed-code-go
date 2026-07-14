package io.spine.embedcode.gradle;

import org.gradle.api.InvalidUserDataException;
import org.gradle.api.file.ConfigurableFileCollection;
import org.gradle.api.file.Directory;
import org.gradle.api.file.DirectoryProperty;
import org.gradle.api.provider.ListProperty;
import org.gradle.api.provider.MapProperty;
import org.gradle.api.provider.Property;
import org.gradle.api.provider.Provider;

/**
 * Configures Embed Code for a Gradle project.
 *
 * <p>The extension maps directly to Embed Code command-line options and does
 * not create or require a YAML configuration file.</p>
 */
public abstract class EmbedCodeExtension {

    /** Returns the release version to download and run, defaulting to the plugin version. */
    public abstract Property<String> getVersion();

    /** Returns the root directory containing source files used by embedding instructions. */
    public abstract DirectoryProperty getCodePath();

    /** Returns named source roots keyed by the name used in embedding instructions. */
    public abstract MapProperty<String, String> getNamedSources();

    /** Returns named source directories with their task dependencies. */
    public abstract ConfigurableFileCollection getNamedSourceDirectories();

    /**
     * Adds a named source root.
     *
     * @param name the name referenced as {@code $name} in an embedding instruction
     * @param directory the source root directory
     */
    public void namedSource(String name, Directory directory) {
        String normalizedName = name.trim();
        if (normalizedName.isEmpty()) {
            throw new InvalidUserDataException("An Embed Code source name must not be empty.");
        }
        getNamedSources().put(normalizedName, directory.getAsFile().getAbsolutePath());
        getNamedSourceDirectories().from(directory);
    }

    /**
     * Adds a named source root supplied by another Gradle provider.
     *
     * @param name the name referenced as {@code $name} in an embedding instruction
     * @param directory the source root provider, including its task dependency
     */
    public void namedSource(String name, Provider<Directory> directory) {
        String normalizedName = name.trim();
        if (normalizedName.isEmpty()) {
            throw new InvalidUserDataException("An Embed Code source name must not be empty.");
        }
        getNamedSources().put(
                normalizedName,
                directory.map(value -> value.getAsFile().getAbsolutePath())
        );
        getNamedSourceDirectories().from(directory);
    }

    /** Returns the root directory containing Markdown or HTML documentation. */
    public abstract DirectoryProperty getDocsPath();

    /** Returns glob patterns selecting documentation files to process. */
    public abstract ListProperty<String> getDocIncludes();

    /** Returns glob patterns selecting documentation files to skip. */
    public abstract ListProperty<String> getDocExcludes();

    /** Returns text inserted between joined fragment parts. */
    public abstract Property<String> getSeparator();

    /** Returns whether Embed Code should print informational log messages. */
    public abstract Property<Boolean> getInfo();

    /** Returns whether Embed Code should print stack traces after panics. */
    public abstract Property<Boolean> getStacktrace();

    /**
     * <p>The plugin appends {@code /v<version>/<platform-asset>} to this URL.
     * This property primarily supports release mirrors and functional testing.</p>
     *
     * @return the base URL containing versioned Embed Code release directories
     */
    public abstract Property<String> getDownloadBaseUrl();
}
