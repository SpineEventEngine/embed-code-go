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
	"golang.org/x/net/context"
	"log/slog"
	"os"
	"runtime/debug"
)

// Handler is a custom slog.Handler that formats log records for simple console output.
//
// It displays each log message in the format:
//
//	HH:MM:SS LEVEL - message
//
// Only messages with level greater than or equal to Handler.Level are printed.
type Handler struct {
	Level slog.Level
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
	return nil
}

// WithAttrs returns a copy of the handler with extra attributes.
// This handler ignores additional attributes and returns itself.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }

// WithGroup returns a copy of the handler for a new group.
// This handler ignores groups and returns itself.
func (h *Handler) WithGroup(name string) slog.Handler { return h }

// HandlePanic is a handler for the panic.
//
// To use, defer this function in any method that calls panic
// or invokes other methods that may call panic.
//
//	defer HandlePanic(withStacktrace)
func HandlePanic(withStacktrace bool) {
	if r := recover(); r != nil {
		fmt.Printf("Panic: %v\n", r)
		if withStacktrace {
			debug.PrintStack()
		}
		os.Exit(1)
	}
}
