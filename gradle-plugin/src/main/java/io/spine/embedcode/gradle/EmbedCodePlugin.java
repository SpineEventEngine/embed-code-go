package io.spine.embedcode.gradle;

import org.gradle.api.Plugin;
import org.gradle.api.Project;
import org.gradle.api.tasks.TaskProvider;

import java.util.Arrays;
import java.util.Collections;

/** Registers automatic installation and execution tasks for Embed Code. */
public final class EmbedCodePlugin implements Plugin<Project> {

    private static final String DEFAULT_DOWNLOAD_BASE_URL =
            "https://github.com/SpineEventEngine/embed-code-go/releases/download";
    private static final String SPINE_TASK_GROUP = "spine";

    /** Applies the plugin to {@code project}. */
    @Override
    public void apply(Project project) {
        EmbedCodeExtension extension = project.getExtensions().create(
                "embedCode",
                EmbedCodeExtension.class
        );
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
                    task.setGroup(SPINE_TASK_GROUP);
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
                "checkEmbedCode",
                "Checks embedded code snippets are up to date",
                "check"
        );
        registerExecutionTask(
                project,
                extension,
                installTask,
                "embedCode",
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
            task.setGroup(SPINE_TASK_GROUP);
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
}
