package io.spine.embedcode.gradle;

import org.gradle.api.GradleException;
import org.gradle.api.Plugin;
import org.gradle.api.Project;
import org.gradle.api.tasks.TaskProvider;

import java.io.IOException;
import java.io.InputStream;
import java.util.Arrays;
import java.util.Collections;
import java.util.Properties;

/** Registers automatic installation and execution tasks for Embed Code. */
public final class EmbedCodePlugin implements Plugin<Project> {

    private static final String DEFAULT_DOWNLOAD_BASE_URL =
            "https://github.com/SpineEventEngine/embed-code-go/releases/download";
    private static final String VERSION_RESOURCE =
            "/io/spine/embedcode/gradle/version.properties";
    private static final String TASK_GROUP = "embed code";

    /** Applies the plugin to {@code project}. */
    @Override
    public void apply(Project project) {
        String checkTaskName = availableTaskName(project, "checkEmbedding");
        String embedTaskName = availableTaskName(project, "embedCode");
        EmbedCodeExtension extension = project.getExtensions().create(
                "embedCode",
                EmbedCodeExtension.class
        );
        extension.getVersion().convention(pluginVersion());
        extension.getDocIncludes().convention(Arrays.asList("**/*.md", "**/*.html"));
        extension.getDocExcludes().convention(Collections.emptyList());
        extension.getNamedSources().convention(Collections.emptyMap());
        extension.getSeparator().convention("...");
        extension.getInfo().convention(false);
        extension.getStacktrace().convention(false);
        extension.getDownloadBaseUrl().convention(DEFAULT_DOWNLOAD_BASE_URL);

        EmbedCodePlatform platform = EmbedCodePlatform.detect(
                System.getProperty("os.name"),
                System.getProperty("os.arch")
        );
        TaskProvider<InstallEmbedCodeTask> installTask = project.getTasks().register(
                "installEmbedCode",
                InstallEmbedCodeTask.class,
                task -> {
                    task.setDescription("Installs the requested Embed Code executable");
                    task.getVersion().set(extension.getVersion());
                    task.getDownloadBaseUrl().set(extension.getDownloadBaseUrl());
                    task.getAssetName().set(platform.getAssetName());
                    task.getExecutableName().set(platform.getExecutableName());
                    task.getExecutableFile().set(
                            project.getLayout().getBuildDirectory().file(
                                    extension.getVersion().map(
                                            version -> "embed-code/" + version
                                                    + '/' + platform.getExecutableName()
                                    )
                            )
                    );
                }
        );

        registerExecutionTask(
                project,
                extension,
                installTask,
                checkTaskName,
                "Checks embedded code snippets are up to date",
                "check"
        );
        registerExecutionTask(
                project,
                extension,
                installTask,
                embedTaskName,
                "Updates embedded code snippets from source files",
                "embed"
        );
    }

    /** Registers one mode-specific execution task backed by {@code installTask}. */
    private static void registerExecutionTask(
            Project project,
            EmbedCodeExtension extension,
            TaskProvider<InstallEmbedCodeTask> installTask,
            String name,
            String description,
            String mode
    ) {
        project.getTasks().register(name, EmbedCodeTask.class, task -> {
            task.setGroup(TASK_GROUP);
            task.setDescription(description);
            task.getMode().set(mode);
            task.getCodePath().set(extension.getCodePath());
            task.getNamedSources().set(extension.getNamedSources());
            task.getNamedSourceDirectories().from(extension.getNamedSourceDirectories());
            task.getDocsPath().set(extension.getDocsPath());
            task.getDocIncludes().set(extension.getDocIncludes());
            task.getDocExcludes().set(extension.getDocExcludes());
            task.getSeparator().set(extension.getSeparator());
            task.getInfo().set(extension.getInfo());
            task.getStacktrace().set(extension.getStacktrace());
            task.getExecutableFile().set(installTask.flatMap(InstallEmbedCodeTask::getExecutableFile));
            task.getWorkingDirectory().set(project.getLayout().getProjectDirectory());
        });
    }

    /** Returns {@code preferredName}, prepending underscores until the task name is unused. */
    private static String availableTaskName(Project project, String preferredName) {
        String candidate = preferredName;
        while (project.getTasks().getNames().contains(candidate)) {
            candidate = '_' + candidate;
        }
        return candidate;
    }

    /** Returns the Embed Code version packaged into the plugin at build time. */
    private static String pluginVersion() {
        Properties properties = new Properties();
        try (InputStream resource = EmbedCodePlugin.class.getResourceAsStream(VERSION_RESOURCE)) {
            if (resource == null) {
                throw new GradleException(
                        "Embed Code plugin version resource is missing: " + VERSION_RESOURCE
                );
            }
            properties.load(resource);
        } catch (IOException error) {
            throw new GradleException("Could not read the Embed Code plugin version.", error);
        }

        String version = properties.getProperty("version", "").trim();
        if (version.isEmpty()) {
            throw new GradleException("The Embed Code plugin version must not be empty.");
        }
        return version;
    }
}
