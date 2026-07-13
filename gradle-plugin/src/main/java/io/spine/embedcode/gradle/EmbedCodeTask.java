package io.spine.embedcode.gradle;

import org.gradle.api.DefaultTask;
import org.gradle.api.GradleException;
import org.gradle.api.file.ConfigurableFileCollection;
import org.gradle.api.file.DirectoryProperty;
import org.gradle.api.file.RegularFileProperty;
import org.gradle.api.provider.ListProperty;
import org.gradle.api.provider.MapProperty;
import org.gradle.api.provider.Property;
import org.gradle.api.tasks.Input;
import org.gradle.api.tasks.InputDirectory;
import org.gradle.api.tasks.InputFile;
import org.gradle.api.tasks.InputFiles;
import org.gradle.api.tasks.Internal;
import org.gradle.api.tasks.Optional;
import org.gradle.api.tasks.PathSensitive;
import org.gradle.api.tasks.PathSensitivity;
import org.gradle.api.tasks.TaskAction;
import org.gradle.process.ExecOperations;
import org.gradle.work.DisableCachingByDefault;

import javax.inject.Inject;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import java.util.Locale;
import java.util.Map;
import java.util.TreeMap;

/** Runs Embed Code in either check or embed mode. */
@DisableCachingByDefault(because = "Embed Code checks or updates documentation files in place")
public abstract class EmbedCodeTask extends DefaultTask {

    /** Returns process execution without project access at execution time. */
    @Inject
    protected abstract ExecOperations getExecOperations();

    /** Returns the execution mode assigned by the plugin. */
    @Input
    public abstract Property<String> getMode();

    /** Returns the source root passed to {@code -code-path}. */
    @InputDirectory
    @Optional
    @PathSensitive(PathSensitivity.RELATIVE)
    public abstract DirectoryProperty getCodePath();

    /** Returns named source roots included in an internally generated configuration. */
    @Input
    public abstract MapProperty<String, String> getNamedSources();

    /** Returns named source directories with their producing task dependencies. */
    @InputFiles
    @PathSensitive(PathSensitivity.RELATIVE)
    public abstract ConfigurableFileCollection getNamedSourceDirectories();

    /** Returns the documentation root passed to {@code -docs-path}. */
    @InputDirectory
    @PathSensitive(PathSensitivity.RELATIVE)
    public abstract DirectoryProperty getDocsPath();

    /** Returns documentation include patterns passed to {@code -doc-includes}. */
    @Input
    public abstract ListProperty<String> getDocIncludes();

    /** Returns documentation exclude patterns passed to {@code -doc-excludes}. */
    @Input
    public abstract ListProperty<String> getDocExcludes();

    /** Returns the fragment separator passed to {@code -separator}. */
    @Input
    public abstract Property<String> getSeparator();

    /** Returns whether informational logging is enabled. */
    @Input
    public abstract Property<Boolean> getInfo();

    /** Returns whether panic stack traces are enabled. */
    @Input
    public abstract Property<Boolean> getStacktrace();

    /** Returns the installed platform executable. */
    @InputFile
    @PathSensitive(PathSensitivity.NONE)
    public abstract RegularFileProperty getExecutableFile();

    /** Returns the process working directory. */
    @Internal
    public abstract DirectoryProperty getWorkingDirectory();

    /** Executes Embed Code with arguments derived from the Gradle extension. */
    @TaskAction
    public void runEmbedCode() {
        Map<String, String> namedSources = new TreeMap<>(getNamedSources().get());
        boolean hasDirectSource = getCodePath().isPresent();
        boolean hasNamedSources = !namedSources.isEmpty();
        if (hasDirectSource == hasNamedSources) {
            throw new GradleException(
                    "Configure exactly one of `codePath` or `namedSource(...)` for Embed Code."
            );
        }

        List<String> arguments = new ArrayList<>();
        arguments.add("-mode=" + getMode().get());
        if (hasNamedSources) {
            arguments.add("-config-path=" + writeNamedSourceConfiguration(namedSources));
        } else {
            arguments.add("-code-path=" + getCodePath().get().getAsFile().getAbsolutePath());
            arguments.add("-docs-path=" + getDocsPath().get().getAsFile().getAbsolutePath());
            if (!getDocIncludes().get().isEmpty()) {
                arguments.add("-doc-includes=" + String.join(",", getDocIncludes().get()));
            }
            if (!getDocExcludes().get().isEmpty()) {
                arguments.add("-doc-excludes=" + String.join(",", getDocExcludes().get()));
            }
            arguments.add("-separator=" + getSeparator().get());
            arguments.add("-info=" + getInfo().get());
            arguments.add("-stacktrace=" + getStacktrace().get());
        }

        getExecOperations().exec(spec -> {
            spec.executable(getExecutableFile().get().getAsFile());
            spec.args(arguments);
            spec.setWorkingDir(getWorkingDirectory().get().getAsFile());
        });
    }

    /** Writes the generated configuration used when named source roots are configured. */
    private Path writeNamedSourceConfiguration(Map<String, String> namedSources) {
        Map<String, String> normalizedSources = new TreeMap<>();
        for (Map.Entry<String, String> source : namedSources.entrySet()) {
            Path path = Paths.get(source.getValue());
            if (!path.isAbsolute()) {
                path = getWorkingDirectory().get().getAsFile().toPath().resolve(path);
            }
            path = path.normalize().toAbsolutePath();
            if (!Files.isDirectory(path)) {
                throw new GradleException(
                        "Embed Code source `" + source.getKey() + "` is not a directory: " + path
                );
            }
            normalizedSources.put(source.getKey(), path.toString());
        }

        String json = createConfigurationJson(
                normalizedSources,
                getDocsPath().get().getAsFile().getAbsolutePath(),
                getDocIncludes().get(),
                getDocExcludes().get(),
                getSeparator().get(),
                getInfo().get(),
                getStacktrace().get()
        );
        Path configuration = getTemporaryDir().toPath().resolve("embed-code.json");
        try {
            Files.write(configuration, json.getBytes(StandardCharsets.UTF_8));
        } catch (IOException exception) {
            throw new GradleException(
                    "Could not write the generated Embed Code configuration to " + configuration + '.',
                    exception
            );
        }
        return configuration;
    }

    /** Creates a JSON document accepted by Embed Code's YAML configuration parser. */
    static String createConfigurationJson(
            Map<String, String> namedSources,
            String docsPath,
            List<String> docIncludes,
            List<String> docExcludes,
            String separator,
            boolean info,
            boolean stacktrace
    ) {
        StringBuilder json = new StringBuilder();
        json.append("{\n  \"code-path\": [\n");
        int index = 0;
        for (Map.Entry<String, String> source : namedSources.entrySet()) {
            if (index > 0) {
                json.append(",\n");
            }
            json.append("    {\"name\": ");
            appendJsonString(json, source.getKey());
            json.append(", \"path\": ");
            appendJsonString(json, source.getValue());
            json.append('}');
            index++;
        }
        json.append("\n  ],\n  \"docs-path\": ");
        appendJsonString(json, docsPath);
        json.append(",\n  \"doc-includes\": ");
        appendJsonArray(json, docIncludes);
        json.append(",\n  \"doc-excludes\": ");
        appendJsonArray(json, docExcludes);
        json.append(",\n  \"separator\": ");
        appendJsonString(json, separator);
        json.append(",\n  \"info\": ").append(info);
        json.append(",\n  \"stacktrace\": ").append(stacktrace);
        json.append("\n}\n");
        return json.toString();
    }

    /** Appends a JSON array containing {@code values}. */
    private static void appendJsonArray(StringBuilder json, List<String> values) {
        json.append('[');
        for (int i = 0; i < values.size(); i++) {
            if (i > 0) {
                json.append(", ");
            }
            appendJsonString(json, values.get(i));
        }
        json.append(']');
    }

    /** Appends {@code value} as an escaped JSON string. */
    private static void appendJsonString(StringBuilder json, String value) {
        json.append('"');
        for (int i = 0; i < value.length(); i++) {
            char character = value.charAt(i);
            switch (character) {
                case '"':
                    json.append("\\\"");
                    break;
                case '\\':
                    json.append("\\\\");
                    break;
                case '\b':
                    json.append("\\b");
                    break;
                case '\f':
                    json.append("\\f");
                    break;
                case '\n':
                    json.append("\\n");
                    break;
                case '\r':
                    json.append("\\r");
                    break;
                case '\t':
                    json.append("\\t");
                    break;
                default:
                    if (character < 0x20) {
                        json.append(String.format(Locale.ROOT, "\\u%04x", (int) character));
                    } else {
                        json.append(character);
                    }
            }
        }
        json.append('"');
    }
}
