package logging

import (
	"fmt"
	"golang.org/x/net/context"
	"log/slog"
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

// Enabled returns true if the log level `l` is greater than or equal to the Handler's Level.
func (h *Handler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.Level
}

// Handle formats the log record and writes it to standard output in a simple readable format:
//
//	HH:MM:SS LEVEL - message
func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	time := r.Time.Format("15:04:05")
	fmt.Printf("%s %s - %s\n",
		time,
		r.Level.String(),
		r.Message,
	)
	return nil
}

// WithAttrs returns a copy of the handler with extra attributes.
// This handler ignores additional attributes and returns itself.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }

// WithGroup returns a copy of the handler for a new group.
// This handler ignores groups and returns itself.
func (h *Handler) WithGroup(name string) slog.Handler { return h }
