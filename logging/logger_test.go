package logging

import (
	"errors"
	"testing"
)

// TestFormatPanicMessage verifies formatting for ordinary and joined panic errors.
func TestFormatPanicMessage(t *testing.T) {
	t.Run("formats single panic value", func(t *testing.T) {
		actual := formatPanicMessage("failed")
		expected := "panic: failed"
		if actual != expected {
			t.Fatalf("expected %q, got %q", expected, actual)
		}
	})

	t.Run("formats joined panic errors as a list", func(t *testing.T) {
		actual := formatPanicMessage(errors.Join(
			errors.New("error1 text"),
			errors.New("error2 text"),
		))
		expected := "panic:\n- error1 text\n- error2 text"
		if actual != expected {
			t.Fatalf("expected %q, got %q", expected, actual)
		}
	})
}
