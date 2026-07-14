package io.spine.embedcode.gradle;

import org.gradle.api.DefaultTask;
import org.gradle.api.GradleException;
import org.gradle.api.file.RegularFileProperty;
import org.gradle.api.provider.Property;
import org.gradle.api.tasks.Input;
import org.gradle.api.tasks.OutputFile;
import org.gradle.api.tasks.TaskAction;
import org.gradle.work.DisableCachingByDefault;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.URI;
import java.net.URLConnection;
import java.nio.file.AtomicMoveNotSupportedException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardCopyOption;
import java.util.zip.ZipEntry;
import java.util.zip.ZipInputStream;

/**
 * Downloads and prepares the Embed Code executable selected for the host.
 *
 * <p>The output file gives Gradle normal up-to-date behavior, so a successfully
 * installed version is reused by later invocations.</p>
 */
@DisableCachingByDefault(
        because = "The downloaded release asset is already reused as a task output"
)
public abstract class InstallEmbedCodeTask extends DefaultTask {

    private static final int CONNECT_TIMEOUT_MILLIS = 30_000;
    private static final int READ_TIMEOUT_MILLIS = 120_000;

    /** Returns the Embed Code release version. */
    @Input
    public abstract Property<String> getVersion();

    /** Returns the base URL containing versioned release directories. */
    @Input
    public abstract Property<String> getDownloadBaseUrl();

    /** Returns the platform-specific release asset name. */
    @Input
    public abstract Property<String> getAssetName();

    /** Returns the executable name expected inside an archive or used directly. */
    @Input
    public abstract Property<String> getExecutableName();

    /** Returns the installed executable used by Embed Code execution tasks. */
    @OutputFile
    public abstract RegularFileProperty getExecutableFile();

    /** Downloads, extracts when necessary, and marks the executable runnable. */
    @TaskAction
    public void install() {
        String requestedVersion = getVersion().get().trim();
        if (requestedVersion.isEmpty()) {
            throw new GradleException("Embed Code version must not be empty.");
        }

        String releaseTag = requestedVersion.startsWith("v")
                ? requestedVersion
                : "v" + requestedVersion;
        String asset = getAssetName().get();
        String baseUrl = trimTrailingSlashes(getDownloadBaseUrl().get());
        URI source = URI.create(baseUrl + '/' + releaseTag + '/' + asset);
        Path destination = getExecutableFile().get().getAsFile().toPath();
        Path download = getTemporaryDir().toPath().resolve(asset);
        Path preparedExecutable = getTemporaryDir().toPath()
                .resolve(getExecutableName().get());

        try {
            Files.createDirectories(destination.getParent());
            getLogger().lifecycle("Downloading Embed Code {} from {}", requestedVersion, source);
            download(source, download);

            if (asset.endsWith(".zip")) {
                extractExecutable(download, getExecutableName().get(), preparedExecutable);
            } else {
                Files.move(download, preparedExecutable, StandardCopyOption.REPLACE_EXISTING);
            }

            if (!preparedExecutable.toFile().setExecutable(true, false)) {
                throw new GradleException(
                        "Could not make `" + preparedExecutable + "` executable."
                );
            }
            moveAtomically(preparedExecutable, destination);
        } catch (IOException exception) {
            throw new GradleException(
                    "Could not install Embed Code from " + source + '.',
                    exception
            );
        }
    }

    /** Downloads {@code source} into {@code destination}, reporting HTTP failures clearly. */
    private static void download(URI source, Path destination) {
        URLConnection connection = null;
        try {
            connection = source.toURL().openConnection();
            connection.setConnectTimeout(CONNECT_TIMEOUT_MILLIS);
            connection.setReadTimeout(READ_TIMEOUT_MILLIS);

            if (connection instanceof HttpURLConnection) {
                HttpURLConnection http = (HttpURLConnection) connection;
                http.setInstanceFollowRedirects(true);
                int status = http.getResponseCode();
                if (status < 200 || status > 299) {
                    throw new GradleException(
                            "Could not download Embed Code: HTTP " + status
                                    + " from " + source + '.'
                    );
                }
            }

            try (InputStream input = connection.getInputStream();
                 OutputStream output = Files.newOutputStream(destination)) {
                copy(input, output);
            }
        } catch (IOException exception) {
            throw new GradleException(
                    "Could not download Embed Code from " + source + '.',
                    exception
            );
        } finally {
            if (connection instanceof HttpURLConnection) {
                ((HttpURLConnection) connection).disconnect();
            }
        }
    }

    /** Extracts {@code entryName} from {@code archive} into {@code destination}. */
    private static void extractExecutable(Path archive, String entryName, Path destination)
            throws IOException {
        try (ZipInputStream zip = new ZipInputStream(Files.newInputStream(archive))) {
            ZipEntry entry = zip.getNextEntry();
            while (entry != null) {
                String fileName = entry.getName();
                int slash = fileName.lastIndexOf('/');
                if (slash >= 0) {
                    fileName = fileName.substring(slash + 1);
                }
                if (!entry.isDirectory() && fileName.equals(entryName)) {
                    try (OutputStream output = Files.newOutputStream(destination)) {
                        copy(zip, output);
                    }
                    return;
                }
                zip.closeEntry();
                entry = zip.getNextEntry();
            }
        }
        throw new GradleException(
                "Archive `" + archive + "` does not contain `" + entryName + "`."
        );
    }

    /** Copies all bytes from {@code input} into {@code output}. */
    private static void copy(InputStream input, OutputStream output) throws IOException {
        byte[] buffer = new byte[8_192];
        int count = input.read(buffer);
        while (count >= 0) {
            output.write(buffer, 0, count);
            count = input.read(buffer);
        }
    }

    /** Moves {@code source} to {@code destination}, atomically when supported. */
    private static void moveAtomically(Path source, Path destination) throws IOException {
        try {
            Files.move(
                    source,
                    destination,
                    StandardCopyOption.ATOMIC_MOVE,
                    StandardCopyOption.REPLACE_EXISTING
            );
        } catch (AtomicMoveNotSupportedException ignored) {
            Files.move(source, destination, StandardCopyOption.REPLACE_EXISTING);
        }
    }

    /** Removes trailing slashes without changing a URL scheme. */
    private static String trimTrailingSlashes(String value) {
        int end = value.length();
        while (end > 0 && value.charAt(end - 1) == '/') {
            end--;
        }
        return value.substring(0, end);
    }
}
