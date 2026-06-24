// Copyright 2026, TeamDev. All rights reserved.
//
// Redistribution and use in source and/or binary forms, with or without
// modification, must retain the above copyright notice and the following
// disclaimer.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package logging

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"

	"golang.org/x/net/context"
)

const fileScheme = "file"

// Handler is a custom slog.Handler that formats log records for simple console output.
//
// It displays each log message in the format:
//
//	HH:MM:SS LEVEL - message
//
// Only messages with level greater than or equal to Handler.Level are printed.
type Handler struct {
	// Level is the minimum enabled logging level.
	Level slog.Level

	// attributes contains attributes added through WithAttrs.
	attributes []slog.Attr

	// groups contains group names added through WithGroup.
	groups []string
}

// Enabled returns true if the log level is greater than or equal to the Handler's Level.
func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.Level
}

// Handle formats the log record and writes it to standard output in a simple readable format:
//
//	HH:MM:SS LEVEL - message
func (h *Handler) Handle(_ context.Context, record slog.Record) error {
	time := record.Time.Format("15:04:05")
	fmt.Printf("%s %s - %s\n",
		time,
		record.Level.String(),
		record.Message,
	)

	prefix := strings.Join(h.groups, ".")
	if prefix != "" {
		prefix += "."
	}

	for _, attr := range h.attributes {
		fmt.Printf(" %s%s=%v\n", prefix, attr.Key, attr.Value)
	}

	record.Attrs(func(attr slog.Attr) bool {
		fmt.Printf(" %s%s=%v\n", prefix, attr.Key, attr.Value)

		return true
	})

	return nil
}

// WithAttrs returns a copy of the handler with extra attributes.
func (h *Handler) WithAttrs(attributes []slog.Attr) slog.Handler {
	newHandler := *h
	newHandler.attributes = append(append([]slog.Attr{}, h.attributes...), attributes...)

	return &newHandler
}

// WithGroup returns a copy of the handler for a new group.
func (h *Handler) WithGroup(name string) slog.Handler {
	newHandler := *h
	newHandler.groups = append(append([]string{}, h.groups...), name)

	return &newHandler
}

// FileReference returns a clickable file URL when the path can be made absolute.
func FileReference(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return path
	}

	return fileURLFromAbsolutePath(absPath)
}

// FileReferenceWithLine returns a clickable file URL with an optional line suffix.
func FileReferenceWithLine(path string, line int) string {
	reference := FileReference(path)
	if line <= 0 {
		return reference
	}

	return reference + ":" + strconv.Itoa(line)
}

// fileURLFromAbsolutePath formats an absolute local path as an OS-neutral file URL.
func fileURLFromAbsolutePath(path string) string {
	normalizedPath := filepath.ToSlash(strings.ReplaceAll(path, "\\", "/"))
	if isWindowsDrivePath(normalizedPath) {
		return (&url.URL{
			Scheme: fileScheme,
			Path:   "/" + normalizedPath,
		}).String()
	}
	if strings.HasPrefix(normalizedPath, "//") {
		withoutSlashes := strings.TrimPrefix(normalizedPath, "//")
		host, pathAfterHost, _ := strings.Cut(withoutSlashes, "/")

		return (&url.URL{
			Scheme: fileScheme,
			Host:   host,
			Path:   "/" + pathAfterHost,
		}).String()
	}

	return (&url.URL{
		Scheme: fileScheme,
		Path:   normalizedPath,
	}).String()
}

// isWindowsDrivePath reports whether a slash-normalized path starts with a drive letter.
func isWindowsDrivePath(path string) bool {
	if len(path) < 2 || path[1] != ':' {
		return false
	}

	driveLetter := path[0]

	return (driveLetter >= 'A' && driveLetter <= 'Z') ||
		(driveLetter >= 'a' && driveLetter <= 'z')
}

// HandlePanic is a handler for the panic.
//
// To use, defer this function in any method that calls panic
// or invokes other methods that may call panic.
//
//	defer HandlePanic(withStacktrace)
func HandlePanic(withStacktrace bool) {
	if r := recover(); r != nil {
		fmt.Println(formatPanicMessage(r))
		if withStacktrace {
			debug.PrintStack()
		}
		os.Exit(1)
	}
}

// formatPanicMessage formats panic values for console output.
func formatPanicMessage(recovered any) string {
	err, isError := recovered.(error)
	if !isError {
		return fmt.Sprintf("panic: %v", recovered)
	}

	return FormatError("panic", err)
}
